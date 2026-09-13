package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

type fakePaymentGateway struct {
	initiateResult *utils.TransferChargeResult
	initiateErr    error
	verifySig      bool

	capturedEmail  string
	capturedAmount float64
	capturedRef    string
}

func (f *fakePaymentGateway) InitiateTransfer(email string, amountNaira float64, reference string) (*utils.TransferChargeResult, error) {
	f.capturedEmail = email
	f.capturedAmount = amountNaira
	f.capturedRef = reference
	if f.initiateErr != nil {
		return nil, f.initiateErr
	}
	return f.initiateResult, nil
}

func (f *fakePaymentGateway) VerifyWebhookSignature(rawBody []byte, signature string) bool {
	return f.verifySig
}

type fakePaymentOrderRepo struct {
	byReference    *models.Order
	byReferenceErr error

	getWithItemsResult *models.Order

	setInitiatedErr   error
	setInitiatedCalls int

	markPaidErr   error
	markPaidCalls int

	updateStatusErr   error
	updateStatusCalls int
	lastStatus        string
}

func (f *fakePaymentOrderRepo) Create(order *models.Order) (*models.Order, error) { return order, nil }
func (f *fakePaymentOrderRepo) CreateOrderItem(item *models.OrderItem) (*models.OrderItem, error) {
	return item, nil
}
func (f *fakePaymentOrderRepo) GetOrderWithItems(id int) (*models.Order, error) {
	return f.getWithItemsResult, nil
}
func (f *fakePaymentOrderRepo) GetOrdersByUserID(userID int) ([]models.Order, error) { return nil, nil }
func (f *fakePaymentOrderRepo) UpdateOrderStatus(id int, status string) error {
	f.updateStatusCalls++
	f.lastStatus = status
	return f.updateStatusErr
}
func (f *fakePaymentOrderRepo) GetOrderByOrderNumber(orderNumber string) (*models.Order, error) {
	return nil, nil
}
func (f *fakePaymentOrderRepo) GetOrderByPhoneAndOrderNumber(phone, orderNumber string) (*models.Order, error) {
	return nil, nil
}
func (f *fakePaymentOrderRepo) CancelOrder(id int) error           { return nil }
func (f *fakePaymentOrderRepo) BeginTransaction() (*sql.Tx, error) { return nil, nil }
func (f *fakePaymentOrderRepo) CreateOrderWithTransaction(tx *sql.Tx, order *models.Order) (*models.Order, error) {
	return order, nil
}
func (f *fakePaymentOrderRepo) CreateOrderItemWithTransaction(tx *sql.Tx, item *models.OrderItem) (*models.OrderItem, error) {
	return item, nil
}
func (f *fakePaymentOrderRepo) GetOrderByPaymentReference(reference string) (*models.Order, error) {
	return f.byReference, f.byReferenceErr
}
func (f *fakePaymentOrderRepo) SetOrderPaymentInitiated(orderID int, reference, accountNumber, accountName, bankName string, expiresAt *time.Time) error {
	f.setInitiatedCalls++
	return f.setInitiatedErr
}
func (f *fakePaymentOrderRepo) MarkOrderPaymentPaid(orderID int, paidAt time.Time) error {
	f.markPaidCalls++
	return f.markPaidErr
}

type fakePaymentNotificationService struct {
	statusUpdateCalls int
	lastOrderID       int
	lastStatus        string
}

func (f *fakePaymentNotificationService) SendOrderConfirmation(order *models.Order) error { return nil }
func (f *fakePaymentNotificationService) SendOrderStatusUpdate(orderID int, newStatus string) error {
	f.statusUpdateCalls++
	f.lastOrderID = orderID
	f.lastStatus = newStatus
	return nil
}
func (f *fakePaymentNotificationService) SendCateringConfirmation(request *models.CateringRequest) error {
	return nil
}
func (f *fakePaymentNotificationService) SendPasswordReset(email, resetToken string) error {
	return nil
}
func (f *fakePaymentNotificationService) SendPhoneVerification(phoneNumber, code string) error {
	return nil
}
func (f *fakePaymentNotificationService) GetNotificationHistory(userID int, limit int) ([]models.Notification, error) {
	return nil, nil
}
func (f *fakePaymentNotificationService) MarkNotificationAsRead(notificationID int) error { return nil }

func emailPtr(s string) *string { return &s }

func TestInitiateTransferPaymentRequiresEmail(t *testing.T) {
	gateway := &fakePaymentGateway{}
	repo := &fakePaymentOrderRepo{}
	svc := NewPaymentService(gateway, repo, &fakePaymentNotificationService{})

	order := &models.Order{ID: 1, OrderNumber: "FG-1", TotalAmount: 100}
	_, err := svc.InitiateTransferPayment(order)

	if err == nil {
		t.Fatal("expected an error when the order has no customer email")
	}
	if repo.setInitiatedCalls != 0 {
		t.Errorf("expected SetOrderPaymentInitiated not to be called, got %d calls", repo.setInitiatedCalls)
	}
}

func TestInitiateTransferPaymentPersistsChargeDetails(t *testing.T) {
	expiresAt := time.Now().Add(15 * time.Minute)
	gateway := &fakePaymentGateway{
		initiateResult: &utils.TransferChargeResult{
			Reference:     "PSK-FG-1-abc123",
			AccountNumber: "1234567890",
			AccountName:   "Funkey Grab & Bite",
			BankName:      "Test Bank",
			ExpiresAt:     &expiresAt,
		},
	}
	finalOrder := &models.Order{ID: 1, OrderNumber: "FG-1", PaymentStatus: models.PaymentStatusPending}
	repo := &fakePaymentOrderRepo{getWithItemsResult: finalOrder}
	svc := NewPaymentService(gateway, repo, &fakePaymentNotificationService{})

	order := &models.Order{ID: 1, OrderNumber: "FG-1", TotalAmount: 5000, CustomerEmail: emailPtr("customer@example.com")}
	result, err := svc.InitiateTransferPayment(order)
	if err != nil {
		t.Fatalf("InitiateTransferPayment() error = %v", err)
	}

	if gateway.capturedEmail != "customer@example.com" {
		t.Errorf("gateway called with email %q, want customer@example.com", gateway.capturedEmail)
	}
	if gateway.capturedAmount != 5000 {
		t.Errorf("gateway called with amount %v, want 5000", gateway.capturedAmount)
	}
	if repo.setInitiatedCalls != 1 {
		t.Errorf("expected SetOrderPaymentInitiated to be called once, got %d", repo.setInitiatedCalls)
	}
	if result != finalOrder {
		t.Error("expected the re-fetched order (with payment fields) to be returned")
	}
}

func TestHandleWebhookRejectsInvalidSignature(t *testing.T) {
	gateway := &fakePaymentGateway{verifySig: false}
	repo := &fakePaymentOrderRepo{}
	svc := NewPaymentService(gateway, repo, &fakePaymentNotificationService{})

	err := svc.HandleWebhook([]byte(`{}`), "bad-signature")
	if !errors.Is(err, ErrInvalidWebhookSignature) {
		t.Fatalf("HandleWebhook() error = %v, want ErrInvalidWebhookSignature", err)
	}
}

func webhookPayload(t *testing.T, event, reference string, amountKobo int64) []byte {
	t.Helper()
	payload := map[string]any{
		"event": event,
		"data": map[string]any{
			"reference": reference,
			"amount":    amountKobo,
			"status":    "success",
			"currency":  "NGN",
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal webhook payload: %v", err)
	}
	return body
}

func TestHandleWebhookIgnoresNonChargeSuccessEvents(t *testing.T) {
	gateway := &fakePaymentGateway{verifySig: true}
	repo := &fakePaymentOrderRepo{}
	svc := NewPaymentService(gateway, repo, &fakePaymentNotificationService{})

	body := webhookPayload(t, "charge.failed", "PSK-FG-1-abc", 500000)
	if err := svc.HandleWebhook(body, "sig"); err != nil {
		t.Fatalf("HandleWebhook() error = %v, want nil (ignored event)", err)
	}
	if repo.markPaidCalls != 0 {
		t.Errorf("expected no payment to be marked paid for a non-success event, got %d calls", repo.markPaidCalls)
	}
}

func TestHandleWebhookMarksOrderPaidAndConfirmed(t *testing.T) {
	gateway := &fakePaymentGateway{verifySig: true}
	order := &models.Order{
		ID:            7,
		OrderNumber:   "FG-7",
		TotalAmount:   5000, // NGN 5000.00 == 500000 kobo
		PaymentStatus: models.PaymentStatusPending,
	}
	repo := &fakePaymentOrderRepo{byReference: order}
	notifier := &fakePaymentNotificationService{}
	svc := NewPaymentService(gateway, repo, notifier)

	body := webhookPayload(t, "charge.success", "PSK-FG-7-abc", 500000)
	if err := svc.HandleWebhook(body, "sig"); err != nil {
		t.Fatalf("HandleWebhook() error = %v", err)
	}

	if repo.markPaidCalls != 1 {
		t.Errorf("expected MarkOrderPaymentPaid to be called once, got %d", repo.markPaidCalls)
	}
	if repo.updateStatusCalls != 1 || repo.lastStatus != string(models.OrderStatusConfirmed) {
		t.Errorf("expected UpdateOrderStatus(confirmed) once, got %d calls with status %q", repo.updateStatusCalls, repo.lastStatus)
	}
	if notifier.statusUpdateCalls != 1 || notifier.lastOrderID != 7 {
		t.Errorf("expected a confirmed-status notification for order 7, got %d calls for order %d", notifier.statusUpdateCalls, notifier.lastOrderID)
	}
}

func TestHandleWebhookIsIdempotentWhenAlreadyPaid(t *testing.T) {
	gateway := &fakePaymentGateway{verifySig: true}
	order := &models.Order{
		ID:            7,
		OrderNumber:   "FG-7",
		TotalAmount:   5000,
		PaymentStatus: models.PaymentStatusPaid,
	}
	repo := &fakePaymentOrderRepo{byReference: order}
	svc := NewPaymentService(gateway, repo, &fakePaymentNotificationService{})

	body := webhookPayload(t, "charge.success", "PSK-FG-7-abc", 500000)
	if err := svc.HandleWebhook(body, "sig"); err != nil {
		t.Fatalf("HandleWebhook() error = %v", err)
	}

	if repo.markPaidCalls != 0 {
		t.Errorf("expected a duplicate webhook delivery not to re-mark payment, got %d calls", repo.markPaidCalls)
	}
	if repo.updateStatusCalls != 0 {
		t.Errorf("expected a duplicate webhook delivery not to re-confirm the order, got %d calls", repo.updateStatusCalls)
	}
}

func TestHandleWebhookRejectsAmountMismatch(t *testing.T) {
	gateway := &fakePaymentGateway{verifySig: true}
	order := &models.Order{
		ID:            7,
		OrderNumber:   "FG-7",
		TotalAmount:   5000, // expects 500000 kobo
		PaymentStatus: models.PaymentStatusPending,
	}
	repo := &fakePaymentOrderRepo{byReference: order}
	svc := NewPaymentService(gateway, repo, &fakePaymentNotificationService{})

	body := webhookPayload(t, "charge.success", "PSK-FG-7-abc", 100000) // wrong amount
	if err := svc.HandleWebhook(body, "sig"); err != nil {
		t.Fatalf("HandleWebhook() error = %v, want nil (logged, not surfaced)", err)
	}

	if repo.markPaidCalls != 0 {
		t.Errorf("expected an amount mismatch not to mark the order paid, got %d calls", repo.markPaidCalls)
	}
}
