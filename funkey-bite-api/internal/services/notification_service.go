package services

import (
	"fmt"
	"log"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

type NotificationType string

const (
	NotificationTypeOrderConfirmation    NotificationType = "order_confirmation"
	NotificationTypeOrderStatusUpdate    NotificationType = "order_status_update"
	NotificationTypeCateringConfirmation NotificationType = "catering_confirmation"
	NotificationTypePasswordReset        NotificationType = "password_reset"
	NotificationTypeVerificationCode     NotificationType = "verification_code"
)

type NotificationService interface {
	SendOrderConfirmation(order *models.Order) error
	SendOrderStatusUpdate(orderID int, newStatus string) error
	SendCateringConfirmation(request *models.CateringRequest) error
	SendPasswordReset(email, resetToken string) error
	SendPhoneVerification(phoneNumber, code string) error
	GetNotificationHistory(userID int, limit int) ([]models.Notification, error)
	MarkNotificationAsRead(notificationID int) error
}

type notificationService struct {
	emailService     utils.EmailService
	smsService       utils.SMSService
	orderRepo        repository.OrderRepository
	userRepo         repository.UserRepository
	cateringRepo     repository.CateringRepository
	notificationRepo repository.NotificationRepository
}

func NewNotificationService(
	emailService utils.EmailService,
	smsService utils.SMSService,
	orderRepo repository.OrderRepository,
	userRepo repository.UserRepository,
	cateringRepo repository.CateringRepository,
	notificationRepo repository.NotificationRepository,
) NotificationService {
	return &notificationService{
		emailService:     emailService,
		smsService:       smsService,
		orderRepo:        orderRepo,
		userRepo:         userRepo,
		cateringRepo:     cateringRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *notificationService) SendOrderConfirmation(order *models.Order) error {
	// Determine recipient info
	var email, phone, name string
	name = order.CustomerName
	phone = order.CustomerPhone

	if order.CustomerEmail != nil && *order.CustomerEmail != "" {
		email = *order.CustomerEmail
	}

	// Send SMS notification
	if phone != "" {
		if err := s.smsService.SendOrderConfirmation(phone, order.OrderNumber, order.TotalAmount); err != nil {
			log.Printf("Failed to send order confirmation SMS: %v", err)
		}
	}

	// Send email notification if email is provided
	if email != "" {
		if err := s.emailService.SendOrderConfirmation(email, name, order.OrderNumber, order.TotalAmount); err != nil {
			log.Printf("Failed to send order confirmation email: %v", err)
		}
	}

	// Record notification in database
	if order.UserID != nil {
		notification := &models.Notification{
			UserID:        *order.UserID,
			Type:          string(NotificationTypeOrderConfirmation),
			Title:         "Order Confirmation",
			Message:       fmt.Sprintf("Your order #%s has been received.", order.OrderNumber),
			IsRead:        false,
			ReferenceID:   &order.ID,
			ReferenceType: "order",
		}
		s.notificationRepo.Create(notification)
	}

	return nil
}

func (s *notificationService) SendOrderStatusUpdate(orderID int, newStatus string) error {
	// Get order details
	order, err := s.orderRepo.GetOrderWithItems(orderID)
	if err != nil || order == nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Determine recipient info
	var email, phone, name string
	name = order.CustomerName
	phone = order.CustomerPhone

	if order.CustomerEmail != nil && *order.CustomerEmail != "" {
		email = *order.CustomerEmail
	}

	// Send SMS notification
	if phone != "" {
		if err := s.smsService.SendOrderStatusUpdate(phone, order.OrderNumber, newStatus); err != nil {
			log.Printf("Failed to send order status update SMS: %v", err)
		}
	}

	// Send email notification if email is provided
	if email != "" {
		if err := s.emailService.SendOrderStatusUpdate(email, name, order.OrderNumber, newStatus); err != nil {
			log.Printf("Failed to send order status update email: %v", err)
		}
	}

	// Record notification in database
	if order.UserID != nil {
		notification := &models.Notification{
			UserID:        *order.UserID,
			Type:          string(NotificationTypeOrderStatusUpdate),
			Title:         "Order Status Updated",
			Message:       fmt.Sprintf("Your order #%s status changed to: %s", order.OrderNumber, newStatus),
			IsRead:        false,
			ReferenceID:   &order.ID,
			ReferenceType: "order",
		}
		s.notificationRepo.Create(notification)
	}

	return nil
}

func (s *notificationService) SendCateringConfirmation(request *models.CateringRequest) error {
	// Determine recipient info
	var email, phone, name, eventName string
	name = request.ContactName
	phone = request.ContactPhone
	eventName = "Your Event"
	if request.EventName != nil && *request.EventName != "" {
		eventName = *request.EventName
	}

	if request.ContactEmail != nil && *request.ContactEmail != "" {
		email = *request.ContactEmail
	}

	requestID := fmt.Sprintf("CATER-%d", request.ID)

	// Send SMS notification
	if phone != "" {
		if err := s.smsService.SendCateringConfirmation(phone, requestID); err != nil {
			log.Printf("Failed to send catering confirmation SMS: %v", err)
		}
	}

	// Send email notification if email is provided
	if email != "" {
		if err := s.emailService.SendCateringConfirmation(email, name, requestID, eventName); err != nil {
			log.Printf("Failed to send catering confirmation email: %v", err)
		}
	}

	// Record notification in database
	if request.UserID != nil {
		notification := &models.Notification{
			UserID:        *request.UserID,
			Type:          string(NotificationTypeCateringConfirmation),
			Title:         "Catering Request Received",
			Message:       fmt.Sprintf("Your catering request #%s has been received.", requestID),
			IsRead:        false,
			ReferenceID:   &request.ID,
			ReferenceType: "catering",
		}
		s.notificationRepo.Create(notification)
	}

	return nil
}

func (s *notificationService) SendPasswordReset(email, resetToken string) error {
	// Send email notification
	if err := s.emailService.SendPasswordReset(email, resetToken); err != nil {
		log.Printf("Failed to send password reset email: %v", err)
		return err
	}

	// Record notification in database (need to get user by email)
	user, err := s.userRepo.FindByPhoneOrEmail("", email)
	if err == nil && user != nil {
		notification := &models.Notification{
			UserID:        user.ID,
			Type:          string(NotificationTypePasswordReset),
			Title:         "Password Reset Requested",
			Message:       "A password reset has been requested for your account.",
			IsRead:        false,
			ReferenceID:   nil,
			ReferenceType: "auth",
		}
		s.notificationRepo.Create(notification)
	}

	return nil
}

func (s *notificationService) SendPhoneVerification(phoneNumber, code string) error {
	// Send SMS notification
	if err := s.smsService.SendVerificationCode(phoneNumber, code); err != nil {
		log.Printf("Failed to send verification code SMS: %v", err)
		return err
	}

	return nil
}

func (s *notificationService) GetNotificationHistory(userID int, limit int) ([]models.Notification, error) {
	return s.notificationRepo.GetByUserID(userID, limit)
}

func (s *notificationService) MarkNotificationAsRead(notificationID int) error {
	return s.notificationRepo.MarkAsRead(notificationID)
}

// Helper method to schedule order ready notifications
func (s *notificationService) ScheduleOrderReadyNotification(orderID int, readyTime time.Time) {
	// This would be implemented with a task queue (e.g., Redis, SQS)
	// For now, just log it
	log.Printf("Scheduled order ready notification for order %d at %v", orderID, readyTime)
}
