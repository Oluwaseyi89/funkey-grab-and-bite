package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

type fakePaymentGateway struct {
	initiateResult *utils.TransferChargeResult
	initiateErr    error
	verifySig      bool

	capturedEmail  string
	capturedAmount float64
	capturedRef    string
	initiateCalls  int
	// failFirst fails this many initiate attempts before succeeding.
	failFirst int

	verifyResult *utils.TransactionStatus
	verifyErr    error
}

func (f *fakePaymentGateway) InitiateTransfer(email string, amountNaira float64, reference string) (*utils.TransferChargeResult, error) {
	f.capturedEmail = email
	f.capturedAmount = amountNaira
	f.capturedRef = reference
	f.initiateCalls++
	if f.initiateCalls <= f.failFirst {
		return nil, errors.New("paystack timeout")
	}
	if f.initiateErr != nil {
		return nil, f.initiateErr
	}
	return f.initiateResult, nil
}

func (f *fakePaymentGateway) VerifyTransaction(reference string) (*utils.TransactionStatus, error) {
	return f.verifyResult, f.verifyErr
}

// fakePaymentStore stands in for the transactions that change an order's payment.
type fakePaymentStore struct {
	applied []repository.TransferPayment
	result  repository.PaymentResult
	err     error
	cancels int
}

func (f *fakePaymentStore) ApplyTransferPayment(ctx context.Context, payment repository.TransferPayment) (repository.PaymentResult, error) {
	f.applied = append(f.applied, payment)
	return f.result, f.err
}

func (f *fakePaymentStore) CancelUnpaidOrder(ctx context.Context, orderID int) (bool, error) {
	f.cancels++
	return true, nil
}

func (f *fakePaymentGateway) VerifyWebhookSignature(rawBody []byte, signature string) bool {
	return f.verifySig
}

type fakePaymentOrderRepo struct {
	getWithItemsResult *models.Order

	setInitiatedErr   error
	setInitiatedCalls int

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
func (f *fakePaymentOrderRepo) SetOrderPaymentInitiated(orderID int, reference, accountNumber, accountName, bankName string, expiresAt *time.Time) error {
	f.setInitiatedCalls++
	return f.setInitiatedErr
}
func (f *fakePaymentOrderRepo) MarkPaymentCollected(orderID int, via string, paidAt time.Time) (bool, error) {
	return false, nil
}

type fakePaymentNotificationService struct {
	statusUpdateCalls int
	lastOrderID       int
	lastStatus        string
}

func (f *fakePaymentNotificationService) SendOrderConfirmation(order *models.Order) error { return nil }
func (f *fakePaymentNotificationService) SendOrderStatusEmail(orderID int, newStatus string) error {
	return nil
}
func (f *fakePaymentNotificationService) SendOrderStatusSMS(orderID int, newStatus string) error {
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
	svc := NewPaymentService(gateway, repo, nil)

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
	svc := NewPaymentService(gateway, repo, nil)

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
	svc := NewPaymentService(gateway, repo, nil)

	_, err := svc.HandleWebhook(context.Background(), []byte(`{}`), "bad-signature")
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
	store := &fakePaymentStore{}
	svc := NewPaymentService(gateway, &fakePaymentOrderRepo{}, store)

	body := webhookPayload(t, "charge.failed", "PSK-FG-1-abc", 500000)
	if _, err := svc.HandleWebhook(context.Background(), body, "sig"); err != nil {
		t.Fatalf("HandleWebhook() error = %v, want nil (ignored event)", err)
	}
	if len(store.applied) != 0 {
		t.Errorf("expected no payment recorded for a non-success event, got %d", len(store.applied))
	}
}

func TestHandleWebhookRecordsPaymentThroughTheStore(t *testing.T) {
	store := &fakePaymentStore{result: repository.PaymentResult{Outcome: repository.PaymentApplied, OrderID: 7, OrderNumber: "FG-7"}}
	svc := NewPaymentService(&fakePaymentGateway{verifySig: true}, &fakePaymentOrderRepo{}, store)

	result, err := svc.HandleWebhook(context.Background(), webhookPayload(t, "charge.success", "PSK-FG-7-abc", 250000), "sig")
	if err != nil {
		t.Fatal(err)
	}
	if len(store.applied) != 1 {
		t.Fatalf("expected the payment recorded once, got %d", len(store.applied))
	}
	got := store.applied[0]
	if got.Event != "charge.success" || got.Reference != "PSK-FG-7-abc" || got.AmountKobo != 250000 || got.Currency != "NGN" {
		t.Fatalf("expected the event's details passed on, got %+v", got)
	}
	if result.Outcome != repository.PaymentApplied || result.OrderID != 7 {
		t.Fatalf("expected the store's result returned, got %+v", result)
	}
}

func TestHandleWebhookAsksForRedeliveryWhenRecordingFails(t *testing.T) {
	store := &fakePaymentStore{err: errors.New("database unavailable")}
	svc := NewPaymentService(&fakePaymentGateway{verifySig: true}, &fakePaymentOrderRepo{}, store)

	_, err := svc.HandleWebhook(context.Background(), webhookPayload(t, "charge.success", "PSK-FG-7-abc", 250000), "sig")
	if err == nil || errors.Is(err, ErrInvalidWebhookSignature) {
		t.Fatalf("a recording failure must be an error so Paystack delivers again, got %v", err)
	}
}

func TestHandleWebhookLeavesCateringPaymentsToCatering(t *testing.T) {
	store := &fakePaymentStore{}
	svc := NewPaymentService(&fakePaymentGateway{verifySig: true}, &fakePaymentOrderRepo{}, store)
	svc.SetCateringHandler(func(reference string, amountKobo int64) (bool, error) { return true, nil })

	if _, err := svc.HandleWebhook(context.Background(), webhookPayload(t, "charge.success", "CATPAY-ABC-1", 500000), "sig"); err != nil {
		t.Fatal(err)
	}
	if len(store.applied) != 0 {
		t.Fatal("a catering payment must not be recorded against an order")
	}
}

func newTestPaymentService(gateway *fakePaymentGateway, repo *fakePaymentOrderRepo, store *fakePaymentStore, now time.Time) *paymentService {
	svc := NewPaymentService(gateway, repo, store).(*paymentService)
	svc.sleep = func(time.Duration) {}
	svc.now = func() time.Time { return now }
	return svc
}

func TestInitiateTransferPaymentRetriesWithFreshReferences(t *testing.T) {
	gateway := &fakePaymentGateway{failFirst: 2, initiateResult: &utils.TransferChargeResult{Reference: "PSK-3", AccountNumber: "0123456789"}}
	repo := &fakePaymentOrderRepo{getWithItemsResult: &models.Order{ID: 1}}
	svc := newTestPaymentService(gateway, repo, nil, time.Now())

	if _, err := svc.InitiateTransferPayment(&models.Order{ID: 1, OrderNumber: "FG-1", CustomerEmail: emailPtr("a@b.c")}); err != nil {
		t.Fatalf("expected success on the third attempt, got %v", err)
	}
	if gateway.initiateCalls != 3 || repo.setInitiatedCalls != 1 {
		t.Fatalf("expected 3 attempts and one save, got %d and %d", gateway.initiateCalls, repo.setInitiatedCalls)
	}
}

func TestInitiateTransferPaymentGivesUpAfterThreeAttempts(t *testing.T) {
	gateway := &fakePaymentGateway{failFirst: 10}
	repo := &fakePaymentOrderRepo{}
	svc := newTestPaymentService(gateway, repo, nil, time.Now())

	if _, err := svc.InitiateTransferPayment(&models.Order{ID: 1, OrderNumber: "FG-1", CustomerEmail: emailPtr("a@b.c")}); err == nil {
		t.Fatal("expected an error after the attempts ran out")
	}
	if gateway.initiateCalls != initiateAttempts || repo.setInitiatedCalls != 0 {
		t.Fatalf("expected %d attempts and no save, got %d and %d", initiateAttempts, gateway.initiateCalls, repo.setInitiatedCalls)
	}
}

func awaitingOrder(expiresAt time.Time) *models.Order {
	ref, account := "PSK-1", "0123456789"
	return &models.Order{
		ID: 7, OrderNumber: "FG-7", TotalAmount: 2500, Status: models.OrderStatusPending,
		PaymentMethod: models.PaymentMethodTransfer, PaymentStatus: models.PaymentStatusPending,
		PaymentReference: &ref, PaymentAccountNumber: &account, PaymentExpiresAt: &expiresAt,
		CustomerEmail: emailPtr("a@b.c"),
	}
}

func TestCheckTransferPaymentWaitsUntilWindowCloses(t *testing.T) {
	now := time.Now()
	expires := now.Add(10 * time.Minute)
	store := &fakePaymentStore{}
	svc := newTestPaymentService(&fakePaymentGateway{}, &fakePaymentOrderRepo{getWithItemsResult: awaitingOrder(expires)}, store, now)

	err := svc.CheckTransferPayment(context.Background(), 7)
	if err == nil || store.cancels != 0 {
		t.Fatalf("expected the check to wait, got %v and %d cancellations", err, store.cancels)
	}
	if !strings.Contains(err.Error(), expires.Add(paymentGrace).Format(time.RFC3339)) {
		t.Fatalf("expected to run again when the window plus grace ends, got %v", err)
	}
}

func TestCheckTransferPaymentMarksPaidWhenWebhookWasMissed(t *testing.T) {
	now := time.Now()
	gateway := &fakePaymentGateway{verifyResult: &utils.TransactionStatus{Status: "success", AmountKobo: 250000}}
	repo := &fakePaymentOrderRepo{getWithItemsResult: awaitingOrder(now.Add(-time.Hour))}
	store := &fakePaymentStore{}
	svc := newTestPaymentService(gateway, repo, store, now)

	if err := svc.CheckTransferPayment(context.Background(), 7); err != nil {
		t.Fatal(err)
	}
	if len(store.applied) != 1 || store.applied[0].Event != "transaction.verify" || store.applied[0].AmountKobo != 250000 || store.cancels != 0 {
		t.Fatalf("expected the verified payment recorded like a webhook, not cancelled; got %+v and %d cancelled", store.applied, store.cancels)
	}
}

func TestCheckTransferPaymentCancelsUnpaidOrder(t *testing.T) {
	now := time.Now()
	gateway := &fakePaymentGateway{verifyResult: &utils.TransactionStatus{Status: "abandoned"}}
	repo := &fakePaymentOrderRepo{getWithItemsResult: awaitingOrder(now.Add(-time.Hour))}
	store := &fakePaymentStore{}
	svc := newTestPaymentService(gateway, repo, store, now)

	if err := svc.CheckTransferPayment(context.Background(), 7); err != nil {
		t.Fatal(err)
	}
	if store.cancels != 1 || len(store.applied) != 0 {
		t.Fatalf("expected the unpaid order cancelled, got %d cancellations and %d payments", store.cancels, len(store.applied))
	}
}

func TestCheckTransferPaymentRetriesWhenPaystackCantBeReached(t *testing.T) {
	now := time.Now()
	gateway := &fakePaymentGateway{verifyErr: errors.New("paystack unreachable")}
	store := &fakePaymentStore{}
	svc := newTestPaymentService(gateway, &fakePaymentOrderRepo{getWithItemsResult: awaitingOrder(now.Add(-time.Hour))}, store, now)

	if err := svc.CheckTransferPayment(context.Background(), 7); err == nil || store.cancels != 0 {
		t.Fatalf("an unknown payment state must be retried, never cancelled; got %v and %d cancellations", err, store.cancels)
	}
}

func TestCheckTransferPaymentLeavesSettledOrdersAlone(t *testing.T) {
	now := time.Now()
	paid := awaitingOrder(now.Add(-time.Hour))
	paid.PaymentStatus = models.PaymentStatusPaid
	store := &fakePaymentStore{}
	svc := newTestPaymentService(&fakePaymentGateway{}, &fakePaymentOrderRepo{getWithItemsResult: paid}, store, now)

	if err := svc.CheckTransferPayment(context.Background(), 7); err != nil || store.cancels != 0 {
		t.Fatalf("a paid order needs no check, got %v and %d cancellations", err, store.cancels)
	}
}

func TestRenewTransferPayment(t *testing.T) {
	now := time.Now()
	newAccount := &utils.TransferChargeResult{Reference: "PSK-2", AccountNumber: "9876543210"}

	t.Run("valid details are kept", func(t *testing.T) {
		gateway := &fakePaymentGateway{initiateResult: newAccount}
		svc := newTestPaymentService(gateway, &fakePaymentOrderRepo{}, nil, now)
		order := awaitingOrder(now.Add(5 * time.Minute))
		got, err := svc.RenewTransferPayment(order)
		if err != nil || got != order || gateway.initiateCalls != 0 {
			t.Fatalf("expected the current details back without a new charge, got %v, %d charges", err, gateway.initiateCalls)
		}
	})
	t.Run("expired details are replaced", func(t *testing.T) {
		gateway := &fakePaymentGateway{initiateResult: newAccount}
		repo := &fakePaymentOrderRepo{getWithItemsResult: &models.Order{ID: 7}}
		svc := newTestPaymentService(gateway, repo, nil, now)
		if _, err := svc.RenewTransferPayment(awaitingOrder(now.Add(-time.Minute))); err != nil || gateway.initiateCalls != 1 || repo.setInitiatedCalls != 1 {
			t.Fatalf("expected new details, got %v, %d charges, %d saves", err, gateway.initiateCalls, repo.setInitiatedCalls)
		}
	})
	t.Run("settled orders are refused", func(t *testing.T) {
		order := awaitingOrder(now.Add(-time.Minute))
		order.Status = models.OrderStatusCancelled
		svc := newTestPaymentService(&fakePaymentGateway{}, &fakePaymentOrderRepo{}, nil, now)
		if _, err := svc.RenewTransferPayment(order); !errors.Is(err, ErrPaymentNotRenewable) {
			t.Fatalf("expected ErrPaymentNotRenewable, got %v", err)
		}
	})
}
