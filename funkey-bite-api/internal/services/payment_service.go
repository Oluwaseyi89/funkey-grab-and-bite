package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

// ErrInvalidWebhookSignature lets callers (the HTTP handler) distinguish a
// rejected signature from any other processing failure without string matching.
var ErrInvalidWebhookSignature = errors.New("invalid webhook signature")

type PaymentService interface {
	// InitiateTransferPayment starts a Pay-with-Transfer charge for an order that
	// was created with payment_method = "transfer", and persists the returned
	// account details on the order.
	InitiateTransferPayment(order *models.Order) (*models.Order, error)
	// HandleWebhook verifies and processes a Paystack webhook payload. rawBody
	// must be the exact, unparsed request body — the signature is computed over
	// those exact bytes.
	HandleWebhook(rawBody []byte, signature string) error
}

type paymentService struct {
	gateway             utils.PaymentGateway
	orderRepo           repository.IOrderRepository
	notificationService NotificationService
}

func NewPaymentService(gateway utils.PaymentGateway, orderRepo repository.IOrderRepository, notificationService NotificationService) PaymentService {
	return &paymentService{
		gateway:             gateway,
		orderRepo:           orderRepo,
		notificationService: notificationService,
	}
}

func (s *paymentService) InitiateTransferPayment(order *models.Order) (*models.Order, error) {
	if order.CustomerEmail == nil || *order.CustomerEmail == "" {
		return nil, fmt.Errorf("a transfer payment requires a customer email")
	}

	reference := utils.GeneratePaymentReference(order.OrderNumber)

	charge, err := s.gateway.InitiateTransfer(*order.CustomerEmail, order.TotalAmount, reference)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate transfer payment: %w", err)
	}

	if err := s.orderRepo.SetOrderPaymentInitiated(
		order.ID,
		charge.Reference,
		charge.AccountNumber,
		charge.AccountName,
		charge.BankName,
		charge.ExpiresAt,
	); err != nil {
		return nil, fmt.Errorf("failed to persist payment initiation: %w", err)
	}

	return s.orderRepo.GetOrderWithItems(order.ID)
}

type paystackWebhookEvent struct {
	Event string `json:"event"`
	Data  struct {
		Reference string `json:"reference"`
		Amount    int64  `json:"amount"`
		Status    string `json:"status"`
		Currency  string `json:"currency"`
	} `json:"data"`
}

func (s *paymentService) HandleWebhook(rawBody []byte, signature string) error {
	if !s.gateway.VerifyWebhookSignature(rawBody, signature) {
		return ErrInvalidWebhookSignature
	}

	var event paystackWebhookEvent
	if err := json.Unmarshal(rawBody, &event); err != nil {
		return fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	if event.Event != "charge.success" {
		log.Printf("[Paystack webhook] ignoring event type %q", event.Event)
		return nil
	}

	if event.Data.Reference == "" {
		return fmt.Errorf("webhook charge.success event carried no reference")
	}

	order, err := s.orderRepo.GetOrderByPaymentReference(event.Data.Reference)
	if err != nil {
		return fmt.Errorf("failed to look up order for payment reference %q: %w", event.Data.Reference, err)
	}
	if order == nil {
		return fmt.Errorf("no order found for payment reference %q", event.Data.Reference)
	}

	// Idempotent: Paystack may retry webhook delivery for the same event.
	if order.PaymentStatus == models.PaymentStatusPaid {
		log.Printf("[Paystack webhook] order #%s already marked paid, ignoring duplicate delivery", order.OrderNumber)
		return nil
	}

	expectedKobo := int64(math.Round(order.TotalAmount * 100))
	if event.Data.Amount != expectedKobo {
		// Don't silently mark this paid — a mismatch here means either a partial
		// transfer or a reference collision, both of which need a human to look at.
		log.Printf(
			"[Paystack webhook] CRITICAL: amount mismatch for order #%s (reference %s): expected %d kobo, got %d kobo — payment left pending for manual review",
			order.OrderNumber, event.Data.Reference, expectedKobo, event.Data.Amount,
		)
		return nil
	}

	if err := s.orderRepo.MarkOrderPaymentPaid(order.ID, time.Now()); err != nil {
		return fmt.Errorf("failed to mark order #%s as paid: %w", order.OrderNumber, err)
	}

	if err := s.orderRepo.UpdateOrderStatus(order.ID, string(models.OrderStatusConfirmed)); err != nil {
		return fmt.Errorf("failed to auto-confirm order #%s after payment: %w", order.OrderNumber, err)
	}

	if err := s.notificationService.SendOrderStatusUpdate(order.ID, string(models.OrderStatusConfirmed)); err != nil {
		log.Printf("Failed to send payment-confirmed notification for order #%s: %v", order.OrderNumber, err)
	}

	return nil
}
