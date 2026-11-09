package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/jobs"
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
	// SetCateringHandler routes catering payments, which share the webhook,
	// to the catering contract service.
	SetCateringHandler(handler CateringPaymentHandler)
	// RenewTransferPayment gives a customer payment details for an unpaid
	// online transfer order whose first attempt failed or whose window
	// closed, and returns the order as it now stands.
	RenewTransferPayment(order *models.Order) (*models.Order, error)
	// CheckTransferPayment is the payment check job: once an online transfer
	// order's payment window has closed it confirms the payment with
	// Paystack, or cancels the unpaid order and returns its stock.
	CheckTransferPayment(ctx context.Context, orderID int) error
}

// ErrPaymentNotRenewable is returned for orders that aren't unpaid online
// transfers, e.g. ones already paid or cancelled.
var ErrPaymentNotRenewable = errors.New("this order is not awaiting an online transfer")

// unpaidOrderCanceller cancels unpaid orders and returns their stock (see
// repository.CheckoutRepository).
type unpaidOrderCanceller interface {
	CancelUnpaidOrder(ctx context.Context, orderID int) (bool, error)
}

const (
	// initiateAttempts bounds how often a checkout asks Paystack for an
	// account before telling the customer to try again.
	initiateAttempts = 3
	// paymentGrace leaves time after the window closes for a transfer made
	// at the last minute to be confirmed.
	paymentGrace = 5 * time.Minute
)

// CateringPaymentHandler records a confirmed catering transfer; handled is
// false when the reference is not a catering payment.
type CateringPaymentHandler func(reference string, amountKobo int64) (handled bool, err error)

type paymentService struct {
	gateway             utils.PaymentGateway
	orderRepo           repository.IOrderRepository
	notificationService NotificationService
	catering            CateringPaymentHandler
	unpaid              unpaidOrderCanceller
	sleep               func(time.Duration)
	now                 func() time.Time
}

func NewPaymentService(gateway utils.PaymentGateway, orderRepo repository.IOrderRepository, notificationService NotificationService, unpaid unpaidOrderCanceller) PaymentService {
	return &paymentService{
		gateway:             gateway,
		orderRepo:           orderRepo,
		notificationService: notificationService,
		unpaid:              unpaid,
		sleep:               time.Sleep,
		now:                 time.Now,
	}
}

func (s *paymentService) SetCateringHandler(handler CateringPaymentHandler) {
	s.catering = handler
}

func (s *paymentService) InitiateTransferPayment(order *models.Order) (*models.Order, error) {
	if order.CustomerEmail == nil || *order.CustomerEmail == "" {
		return nil, fmt.Errorf("a transfer payment requires a customer email")
	}

	// Each attempt uses a new reference, so a charge whose response was lost
	// can't block the retry; the customer only ever sees the one returned.
	var charge *utils.TransferChargeResult
	var err error
	for attempt := 1; attempt <= initiateAttempts; attempt++ {
		reference := utils.GeneratePaymentReference(order.OrderNumber)
		charge, err = s.gateway.InitiateTransfer(*order.CustomerEmail, order.TotalAmount, reference)
		if err == nil {
			break
		}
		if attempt < initiateAttempts {
			s.sleep(time.Duration(attempt) * 500 * time.Millisecond)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to initiate transfer payment after %d attempts: %w", initiateAttempts, err)
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

	if s.catering != nil {
		handled, err := s.catering(event.Data.Reference, event.Data.Amount)
		if handled {
			return err
		}
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

	return s.markOrderPaid(order)
}

// markOrderPaid records a confirmed transfer and confirms the order.
func (s *paymentService) markOrderPaid(order *models.Order) error {
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

// awaitingTransfer reports an online transfer order that is still unpaid and
// open.
func awaitingTransfer(order *models.Order) bool {
	return order.PaymentMethod == models.PaymentMethodTransfer &&
		order.PaymentStatus == models.PaymentStatusPending &&
		order.Status != models.OrderStatusCancelled
}

func (s *paymentService) RenewTransferPayment(order *models.Order) (*models.Order, error) {
	if order == nil || !awaitingTransfer(order) {
		return nil, ErrPaymentNotRenewable
	}
	// Details that are still valid are returned as they are, so pressing the
	// button twice doesn't replace an account the customer may be paying into.
	if order.PaymentAccountNumber != nil && *order.PaymentAccountNumber != "" &&
		order.PaymentExpiresAt != nil && s.now().Before(*order.PaymentExpiresAt) {
		return order, nil
	}
	return s.InitiateTransferPayment(order)
}

func (s *paymentService) CheckTransferPayment(ctx context.Context, orderID int) error {
	order, err := s.orderRepo.GetOrderWithItems(orderID)
	if err != nil {
		return fmt.Errorf("failed to load order %d: %w", orderID, err)
	}
	if order == nil {
		return jobs.Permanent(fmt.Errorf("order %d not found", orderID))
	}
	if !awaitingTransfer(order) {
		return nil
	}

	// The check was queued for the end of the first payment window; a
	// customer who got new details since has until their window closes.
	if order.PaymentExpiresAt != nil {
		if deadline := order.PaymentExpiresAt.Add(paymentGrace); s.now().Before(deadline) {
			return jobs.RunAgainAt(deadline)
		}
	}

	if order.PaymentReference != nil && *order.PaymentReference != "" {
		status, err := s.gateway.VerifyTransaction(*order.PaymentReference)
		if err != nil {
			return fmt.Errorf("failed to verify payment for order #%s: %w", order.OrderNumber, err)
		}
		if status.Paid() {
			expectedKobo := int64(math.Round(order.TotalAmount * 100))
			if status.AmountKobo != expectedKobo {
				log.Printf("[payment check] CRITICAL: order #%s was paid %d kobo, expected %d; left pending for review",
					order.OrderNumber, status.AmountKobo, expectedKobo)
				return nil
			}
			log.Printf("[payment check] order #%s was paid but its webhook never arrived; marking it paid", order.OrderNumber)
			return s.markOrderPaid(order)
		}
	}

	cancelled, err := s.unpaid.CancelUnpaidOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to cancel unpaid order #%s: %w", order.OrderNumber, err)
	}
	if cancelled {
		log.Printf("[payment check] order #%s was not paid in time; cancelled and its stock returned", order.OrderNumber)
	}
	return nil
}
