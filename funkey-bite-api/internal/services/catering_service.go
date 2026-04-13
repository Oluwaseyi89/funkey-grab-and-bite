package services

import (
	"fmt"
	"log"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
)

type CateringService interface {
	CreateRequest(input models.CateringRequestInput, userID *int) (*models.CateringRequest, error)
	GetRequestByID(id int) (*models.CateringRequest, error)
	GetUserRequests(userID int) ([]models.CateringRequest, error)
	GetAllRequests() ([]models.CateringRequest, error)
	UpdateRequestStatus(id int, status string) error
	GetPackages() []models.CateringPackage
}

type cateringService struct {
	cateringRepo        repository.CateringRepository
	notificationService NotificationService
}

func NewCateringService(cateringRepo repository.CateringRepository, notificationService NotificationService) CateringService {
	return &cateringService{
		cateringRepo:        cateringRepo,
		notificationService: notificationService,
	}
}

func (s *cateringService) CreateRequest(input models.CateringRequestInput, userID *int) (*models.CateringRequest, error) {
	eventDate, err := time.Parse("2006-01-02", input.EventDate)
	if err != nil {
		return nil, fmt.Errorf("invalid event date format. Use YYYY-MM-DD")
	}

	if eventDate.Before(time.Now().AddDate(0, 0, -1)) {
		return nil, fmt.Errorf("event date must be in the future")
	}

	if input.GuestCount < 1 {
		return nil, fmt.Errorf("guest count must be at least 1")
	}
	if input.GuestCount > 1000 {
		return nil, fmt.Errorf("guest count cannot exceed 1000")
	}

	request := &models.CateringRequest{
		UserID:          userID,
		EventName:       input.EventName,
		ContactName:     input.ContactName,
		ContactPhone:    input.ContactPhone,
		ContactEmail:    input.ContactEmail,
		EventDate:       input.EventDate,
		EventTime:       input.EventTime,
		GuestCount:      input.GuestCount,
		EventType:       input.EventType,
		Package:         input.Package,
		Budget:          input.Budget,
		SpecialRequests: input.SpecialRequests,
		Status:          models.CateringStatusPending,
		CreatedAt:       time.Now(),
	}

	createdRequest, err := s.cateringRepo.Create(request)
	if err != nil {
		return nil, fmt.Errorf("failed to create catering request: %w", err)
	}

	go func() {
		if err := s.notificationService.SendCateringConfirmation(createdRequest); err != nil {
			log.Printf("Failed to send catering confirmation notification: %v", err)
		}
	}()

	return createdRequest, nil
}

func (s *cateringService) GetRequestByID(id int) (*models.CateringRequest, error) {
	request, err := s.cateringRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get catering request: %w", err)
	}
	if request == nil {
		return nil, fmt.Errorf("catering request not found")
	}
	return request, nil
}

func (s *cateringService) GetUserRequests(userID int) ([]models.CateringRequest, error) {
	requests, err := s.cateringRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user catering requests: %w", err)
	}
	return requests, nil
}

func (s *cateringService) GetAllRequests() ([]models.CateringRequest, error) {
	requests, err := s.cateringRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all catering requests: %w", err)
	}
	return requests, nil
}

func (s *cateringService) UpdateRequestStatus(id int, status string) error {
	validStatuses := map[string]bool{
		string(models.CateringStatusPending):   true,
		string(models.CateringStatusConfirmed): true,
		string(models.CateringStatusDeclined):  true,
		string(models.CateringStatusCompleted): true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	return s.cateringRepo.UpdateStatus(id, status)
}

func (s *cateringService) GetPackages() []models.CateringPackage {
	return []models.CateringPackage{
		{
			ID:             "standard",
			Name:           "Standard Package",
			Description:    "Perfect for casual gatherings and office meetings",
			PricePerPerson: 18.99,
			MinGuests:      10,
			Includes: []string{
				"Buffet-style setup",
				"Basic beverages",
				"Disposable serveware",
				"2-hour service",
			},
		},
		{
			ID:             "premium",
			Name:           "Premium Package",
			Description:    "Ideal for weddings and corporate events",
			PricePerPerson: 24.99,
			MinGuests:      20,
			Includes: []string{
				"Plated or buffet service",
				"Premium beverage selection",
				"Professional staff",
				"Table setup and decor",
				"4-hour service",
			},
		},
		{
			ID:             "executive",
			Name:           "Executive Package",
			Description:    "Luxury service for high-profile events",
			PricePerPerson: 34.99,
			MinGuests:      50,
			MaxGuests:      intPtr(200),
			Includes: []string{
				"Plated gourmet service",
				"Premium open bar",
				"Event coordinator",
				"Custom menu planning",
				"Full setup and cleanup",
				"6-hour service",
			},
		},
	}
}

func intPtr(i int) *int {
	return &i
}
