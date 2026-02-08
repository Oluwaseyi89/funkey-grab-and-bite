package services

import (
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
)

type PromotionService interface {
	CreatePromotion(promotion *models.PromotionCreate) (*models.Promotion, error)
	GetPromotionByID(id int) (*models.Promotion, error)
	GetPromotionByCode(code string) (*models.Promotion, error)
	GetAllPromotions(page, limit int, status string) ([]models.Promotion, int, error)
	UpdatePromotion(id int, updates *models.PromotionUpdate) (*models.Promotion, error)
	DeletePromotion(id int) error
	ValidatePromotion(code string, orderAmount float64, customerID *int) (*models.PromotionValidation, error)
	ApplyPromotion(promotionID, orderID int, customerID *int, orderAmount float64) (float64, error)
	GetActivePromotions() ([]models.Promotion, error)
}

type promotionService struct {
	promotionRepo repository.PromotionRepository
}

func NewPromotionService(promotionRepo repository.PromotionRepository) PromotionService {
	return &promotionService{
		promotionRepo: promotionRepo,
	}
}

func (s *promotionService) CreatePromotion(create *models.PromotionCreate) (*models.Promotion, error) {
	// Check if code already exists
	existing, err := s.promotionRepo.GetByCode(create.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to check promotion code: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("promotion code already exists")
	}

	// Validate dates
	if create.ValidFrom.After(create.ValidUntil) {
		return nil, fmt.Errorf("valid from date must be before valid until date")
	}

	promotion := &models.Promotion{
		Code:           create.Code,
		Title:          create.Title,
		Description:    create.Description,
		PromotionType:  create.PromotionType,
		DiscountValue:  create.DiscountValue,
		MaxDiscount:    create.MaxDiscount,
		MinOrderAmount: create.MinOrderAmount,
		ValidFrom:      create.ValidFrom,
		ValidUntil:     create.ValidUntil,
		UsageLimit:     create.UsageLimit,
		UsedCount:      0,
		IsActive:       create.IsActive,
	}

	return s.promotionRepo.Create(promotion)
}

func (s *promotionService) GetPromotionByID(id int) (*models.Promotion, error) {
	return s.promotionRepo.GetByID(id)
}

func (s *promotionService) GetPromotionByCode(code string) (*models.Promotion, error) {
	return s.promotionRepo.GetByCode(code)
}

func (s *promotionService) GetAllPromotions(page, limit int, status string) ([]models.Promotion, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Convert status string to boolean for query
	var statusFilter string
	switch status {
	case "active":
		statusFilter = "true"
	case "inactive":
		statusFilter = "false"
	default:
		statusFilter = ""
	}

	return s.promotionRepo.GetAll(limit, offset, statusFilter)
}

func (s *promotionService) UpdatePromotion(id int, updates *models.PromotionUpdate) (*models.Promotion, error) {
	// Get existing promotion
	promotion, err := s.promotionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if promotion == nil {
		return nil, fmt.Errorf("promotion not found")
	}

	// Apply updates
	if updates.Title != nil {
		promotion.Title = *updates.Title
	}
	if updates.Description != nil {
		promotion.Description = *updates.Description
	}
	if updates.PromotionType != nil {
		promotion.PromotionType = *updates.PromotionType
	}
	if updates.DiscountValue != nil {
		promotion.DiscountValue = *updates.DiscountValue
	}
	if updates.MaxDiscount != nil {
		promotion.MaxDiscount = updates.MaxDiscount
	}
	if updates.MinOrderAmount != nil {
		promotion.MinOrderAmount = updates.MinOrderAmount
	}
	if updates.ValidUntil != nil {
		promotion.ValidUntil = *updates.ValidUntil
	}
	if updates.UsageLimit != nil {
		promotion.UsageLimit = updates.UsageLimit
	}
	if updates.IsActive != nil {
		promotion.IsActive = *updates.IsActive
	}

	// Validate dates
	if promotion.ValidFrom.After(promotion.ValidUntil) {
		return nil, fmt.Errorf("valid from date must be before valid until date")
	}

	err = s.promotionRepo.Update(promotion)
	if err != nil {
		return nil, err
	}

	return s.promotionRepo.GetByID(id)
}

func (s *promotionService) DeletePromotion(id int) error {
	return s.promotionRepo.Delete(id)
}

func (s *promotionService) ValidatePromotion(code string, orderAmount float64, customerID *int) (*models.PromotionValidation, error) {
	promotion, err := s.promotionRepo.GetByCode(code)
	if err != nil {
		return &models.PromotionValidation{IsValid: false, Message: "Failed to validate promotion"}, err
	}

	if promotion == nil {
		return &models.PromotionValidation{IsValid: false, Message: "Promotion code not found"}, nil
	}

	// Check if promotion is active
	if !promotion.IsActive {
		return &models.PromotionValidation{IsValid: false, Message: "Promotion is not active"}, nil
	}

	// Check validity period
	now := time.Now()
	if now.Before(promotion.ValidFrom) {
		return &models.PromotionValidation{IsValid: false, Message: "Promotion not yet valid"}, nil
	}
	if now.After(promotion.ValidUntil) {
		return &models.PromotionValidation{IsValid: false, Message: "Promotion has expired"}, nil
	}

	// Check usage limit
	if promotion.UsageLimit != nil && promotion.UsedCount >= *promotion.UsageLimit {
		return &models.PromotionValidation{IsValid: false, Message: "Promotion usage limit reached"}, nil
	}

	// Check minimum order amount
	if promotion.MinOrderAmount != nil && orderAmount < *promotion.MinOrderAmount {
		return &models.PromotionValidation{
			IsValid: false,
			Message: fmt.Sprintf("Minimum order amount of $%.2f required", *promotion.MinOrderAmount),
		}, nil
	}

	// Calculate discount
	discount := s.calculateDiscount(promotion, orderAmount)

	return &models.PromotionValidation{
		IsValid:     true,
		Message:     "Promotion is valid",
		Discount:    discount,
		PromotionID: promotion.ID,
	}, nil
}

func (s *promotionService) ApplyPromotion(promotionID, orderID int, customerID *int, orderAmount float64) (float64, error) {
	promotion, err := s.promotionRepo.GetByID(promotionID)
	if err != nil || promotion == nil {
		return 0, fmt.Errorf("promotion not found")
	}

	// Validate promotion again
	validation, err := s.ValidatePromotion(promotion.Code, orderAmount, customerID)
	if err != nil || !validation.IsValid {
		return 0, fmt.Errorf("promotion validation failed: %s", validation.Message)
	}

	// Calculate discount
	discount := s.calculateDiscount(promotion, orderAmount)

	// Record usage
	usage := &models.PromotionUsage{
		PromotionID:     promotionID,
		OrderID:         orderID,
		CustomerID:      customerID,
		DiscountApplied: discount,
	}

	if err := s.promotionRepo.RecordUsage(usage); err != nil {
		return 0, fmt.Errorf("failed to record promotion usage: %w", err)
	}

	// Increment usage count
	if err := s.promotionRepo.IncrementUsage(promotionID); err != nil {
		return 0, fmt.Errorf("failed to increment promotion usage: %w", err)
	}

	return discount, nil
}

func (s *promotionService) GetActivePromotions() ([]models.Promotion, error) {
	promotions, _, err := s.promotionRepo.GetAll(100, 0, "true")
	if err != nil {
		return nil, err
	}

	// Filter for currently valid promotions
	now := time.Now()
	var activePromotions []models.Promotion
	for _, promotion := range promotions {
		if promotion.IsActive && now.After(promotion.ValidFrom) && now.Before(promotion.ValidUntil) {
			activePromotions = append(activePromotions, promotion)
		}
	}

	return activePromotions, nil
}

func (s *promotionService) calculateDiscount(promotion *models.Promotion, orderAmount float64) float64 {
	var discount float64

	switch promotion.PromotionType {
	case models.PromotionTypePercentage:
		discount = orderAmount * (promotion.DiscountValue / 100)
		if promotion.MaxDiscount != nil && discount > *promotion.MaxDiscount {
			discount = *promotion.MaxDiscount
		}

	case models.PromotionTypeFixed:
		discount = promotion.DiscountValue
		if discount > orderAmount {
			discount = orderAmount
		}

	case models.PromotionTypeBOGO:
		// Buy One Get One - implement based on your business logic
		// For now, return 0
		discount = 0
	}

	return discount
}
