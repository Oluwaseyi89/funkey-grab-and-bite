package services

import (
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
)

type SettingsService interface {
	GetSettings() (*models.BusinessSettings, error)
	UpdateSettings(updates *models.BusinessSettingsUpdate) (*models.BusinessSettings, error)
	GetOpeningHours() ([]models.OpeningHours, error)
	IsBusinessOpen() (bool, error)
	ValidateOrderTime(orderType string, requestedTime *time.Time) (bool, string)
	CanAcceptOrders(orderType string) (bool, string)
	CalculateOrderFees(subtotal float64, orderType string) (deliveryFee, taxAmount, total float64)
	ValidateMinimumOrder(subtotal float64, orderType string) (bool, string)
}

type settingsService struct {
	settingsRepo repository.SettingsRepository
	location     *time.Location
	now          func() time.Time
}

func NewSettingsService(settingsRepo repository.SettingsRepository) SettingsService {
	return &settingsService{
		settingsRepo: settingsRepo,
		location:     LoadBusinessLocation(),
		now:          time.Now,
	}
}

func (s *settingsService) GetSettings() (*models.BusinessSettings, error) {
	return s.settingsRepo.GetSettings()
}

func (s *settingsService) UpdateSettings(updates *models.BusinessSettingsUpdate) (*models.BusinessSettings, error) {
	if updates.DeliveryFee != nil && *updates.DeliveryFee < 0 {
		return nil, fmt.Errorf("delivery fee cannot be negative")
	}
	if updates.MinOrderAmount != nil && *updates.MinOrderAmount < 0 {
		return nil, fmt.Errorf("minimum order amount cannot be negative")
	}
	if updates.TaxRate != nil && (*updates.TaxRate < 0 || *updates.TaxRate > 100) {
		return nil, fmt.Errorf("tax rate must be between 0 and 100")
	}
	if err := normalizeSettingsUpdate(updates); err != nil {
		return nil, err
	}

	return s.settingsRepo.UpdateSettings(updates)
}

func (s *settingsService) GetOpeningHours() ([]models.OpeningHours, error) {
	settings, err := s.settingsRepo.GetSettings()
	if err != nil {
		return nil, err
	}

	return settings.OpeningHours, nil
}

func (s *settingsService) IsBusinessOpen() (bool, error) {
	settings, err := s.settingsRepo.GetSettings()
	if err != nil {
		return false, err
	}

	return isOpenAt(settings.OpeningHours, s.now().In(s.location)), nil
}

func (s *settingsService) ValidateOrderTime(orderType string, requestedTime *time.Time) (bool, string) {
	if requestedTime == nil {
		return true, ""
	}

	settings, err := s.settingsRepo.GetSettings()
	if err != nil {
		return false, "Unable to validate order time"
	}

	now := s.now()
	if requestedTime.Before(now) {
		return false, "Order time must be in the future"
	}

	maxFuture := now.Add(7 * 24 * time.Hour)
	if requestedTime.After(maxFuture) {
		return false, "Orders can only be placed up to 7 days in advance"
	}

	local := requestedTime.In(s.location)
	if isOpenAt(settings.OpeningHours, local) {
		return true, ""
	}

	requestedDay := weekdayNames[local.Weekday()]
	if hours, found := hoursFor(settings.OpeningHours, local.Weekday()); found && hours.IsOpen {
		return false, fmt.Sprintf("Business is closed at %s. %s hours: %s - %s",
			local.Format("15:04"), requestedDay, hours.OpenTime, hours.CloseTime)
	}

	return false, "Business is closed on " + requestedDay
}

func (s *settingsService) CanAcceptOrders(orderType string) (bool, string) {
	settings, err := s.settingsRepo.GetSettings()
	if err != nil {
		return false, "Unable to check order acceptance"
	}

	isOpen, err := s.IsBusinessOpen()
	if err != nil || !isOpen {
		return false, "Business is currently closed"
	}

	if orderType == "delivery" && !settings.IsDeliveryOpen {
		return false, "Delivery service is currently unavailable"
	}
	if orderType == "pickup" && !settings.IsPickupOpen {
		return false, "Pickup service is currently unavailable"
	}

	return true, ""
}

func (s *settingsService) CalculateOrderFees(subtotal float64, orderType string) (deliveryFee, taxAmount, total float64) {
	settings, err := s.settingsRepo.GetSettings()
	if err != nil {
		return 0, 0, subtotal
	}

	if orderType == "delivery" {
		deliveryFee = settings.DeliveryFee
	}

	taxAmount = subtotal * (settings.TaxRate / 100)

	total = subtotal + deliveryFee + taxAmount

	return deliveryFee, taxAmount, total
}

func (s *settingsService) ValidateMinimumOrder(subtotal float64, orderType string) (bool, string) {
	settings, err := s.settingsRepo.GetSettings()
	if err != nil {
		return false, "Unable to validate order minimum"
	}

	minAmount := settings.MinOrderAmount
	if subtotal < minAmount {
		return false, fmt.Sprintf("Minimum order amount is ₦%.2f", minAmount)
	}

	return true, ""
}

func (s *settingsService) GetPublicInfo() (map[string]interface{}, error) {
	settings, err := s.settingsRepo.GetSettings()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"businessName":   settings.BusinessName,
		"phoneNumber":    settings.PhoneNumber,
		"email":          settings.Email,
		"address":        settings.Address,
		"city":           settings.City,
		"state":          settings.State,
		"postalCode":     settings.PostalCode,
		"country":        settings.Country,
		"openingHours":   settings.OpeningHours,
		"deliveryFee":    settings.DeliveryFee,
		"minOrderAmount": settings.MinOrderAmount,
		"taxRate":        settings.TaxRate,
		"isDeliveryOpen": settings.IsDeliveryOpen,
		"isPickupOpen":   settings.IsPickupOpen,
	}, nil
}
