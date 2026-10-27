package services

import (
	"bytes"
	"errors"
	"log"
	"strings"
	"sync"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

type failingNotificationService struct {
	orderErr       error
	cateringErr    error
	orderCalled    chan struct{}
	cateringCalled chan struct{}
}

func (f *failingNotificationService) SendOrderConfirmation(order *models.Order) error {
	select {
	case f.orderCalled <- struct{}{}:
	default:
	}
	return f.orderErr
}

func (f *failingNotificationService) SendOrderStatusUpdate(orderID int, newStatus string) error {
	return nil
}

func (f *failingNotificationService) SendCateringConfirmation(request *models.CateringRequest) error {
	select {
	case f.cateringCalled <- struct{}{}:
	default:
	}
	return f.cateringErr
}

func (f *failingNotificationService) SendPasswordReset(email, resetToken string) error { return nil }
func (f *failingNotificationService) SendPhoneVerification(phoneNumber, code string) error {
	return nil
}
func (f *failingNotificationService) GetNotificationHistory(userID int, limit int) ([]models.Notification, error) {
	return nil, nil
}
func (f *failingNotificationService) MarkNotificationAsRead(notificationID int) error { return nil }

type lockedBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (l *lockedBuffer) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.Write(p)
}

func (l *lockedBuffer) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.String()
}

func waitForSignal(t *testing.T, ch <-chan struct{}, label string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(300 * time.Millisecond):
		t.Fatalf("timed out waiting for %s", label)
	}
}

func waitForLogContains(t *testing.T, logs *lockedBuffer, want string) {
	t.Helper()
	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		if strings.Contains(logs.String(), want) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("expected log to contain %q, got: %s", want, logs.String())
}

func TestCreateOrderSucceedsWhenOrderConfirmationNotificationFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer db.Close()

	orderRepo := repository.NewOrderRepository(db)
	menu := &fakeMenuReader{items: map[int]*models.MenuItem{
		1: {ID: 1, Name: "Burger", Price: 12.5, IsAvailable: true},
	}}
	stock := &fakeStockService{stock: map[int]int{1: 5}}
	notifier := &failingNotificationService{
		orderErr:    errors.New("ses failure"),
		orderCalled: make(chan struct{}, 1),
	}
	svc := NewOrderService(orderRepo, menu, stock, notifier)

	logs := &lockedBuffer{}
	originalOutput := log.Writer()
	log.SetOutput(logs)
	defer log.SetOutput(originalOutput)

	createdAt := time.Now()
	order := &models.Order{
		OrderNumber:   "FG-2026-04-2001",
		CustomerName:  "Notification Test",
		CustomerPhone: "+15550003333",
		OrderType:     models.OrderTypeDelivery,
		Status:        models.OrderStatusPending,
		TotalAmount:   37.5,
		CreatedAt:     createdAt,
	}
	items := []models.OrderItemRequest{{
		MenuItemID: 1,
		Name:       "Burger",
		Quantity:   3,
		UnitPrice:  12.5,
	}}

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO orders").WillReturnRows(
		sqlmock.NewRows([]string{"id", "created_at"}).AddRow(101, createdAt),
	)
	mock.ExpectQuery("INSERT INTO order_items").WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(201),
	)
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT id, order_number, user_id, customer_name, customer_phone").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "order_number", "user_id", "customer_name", "customer_phone",
			"customer_email", "order_type", "status", "total_amount", "notes",
			"pickup_time", "created_at", "payment_method", "payment_status",
			"payment_reference", "payment_account_number", "payment_account_name",
			"payment_bank_name", "payment_expires_at", "payment_paid_at",
			"delivery_address", "invoice_number", "subtotal_amount", "delivery_fee",
			"tax_amount", "discount_amount", "payment_collected_via",
		}).AddRow(
			101, "FG-2026-04-2001", nil, "Notification Test", "+15550003333", nil,
			string(models.OrderTypeDelivery), string(models.OrderStatusPending), 37.5, nil, nil, createdAt,
			string(models.PaymentMethodCash), string(models.PaymentStatusNotRequired),
			nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil, nil, nil,
		))
	mock.ExpectQuery("SELECT id, order_id, menu_item_id, name, quantity, unit_price, special_instructions").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "order_id", "menu_item_id", "name", "quantity", "unit_price", "special_instructions",
		}).AddRow(201, 101, 1, "Burger", 3, 12.5, nil))

	createdOrder, err := svc.CreateOrder(order, items)
	if err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}
	if createdOrder == nil || createdOrder.ID != 101 {
		t.Fatalf("expected persisted order ID 101, got %+v", createdOrder)
	}
	if createdOrder.TotalAmount != 37.5 {
		t.Fatalf("expected persisted total amount 37.5, got %v", createdOrder.TotalAmount)
	}
	if got := stock.stock[1]; got != 2 {
		t.Fatalf("expected stock deduction to persist despite notification failure, got %d", got)
	}

	waitForSignal(t, notifier.orderCalled, "order confirmation notification attempt")
	waitForLogContains(t, logs, "Failed to send order confirmation notification")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations failed: %v", err)
	}
}

func TestCreateCateringRequestSucceedsWhenNotificationFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	defer db.Close()

	repo := repository.NewCateringRepository(db)
	notifier := &failingNotificationService{
		cateringErr:    errors.New("twilio failure"),
		cateringCalled: make(chan struct{}, 1),
	}
	// A custom request with menu notes needs neither packages nor menu lookups.
	svc := NewCateringService(repo, nil, nil, notifier)

	logs := &lockedBuffer{}
	originalOutput := log.Writer()
	log.SetOutput(logs)
	defer log.SetOutput(originalOutput)

	eventName := "Office Party"
	menuNotes := "Jollof rice, grilled chicken and small chops"
	budget := 500000.0

	eventDate := time.Now().AddDate(0, 0, 10).Format("2006-01-02")
	input := models.CateringRequestInput{
		RequestType:   models.CateringRequestTypeCustom,
		EventName:     &eventName,
		ContactName:   "Planner",
		ContactPhone:  "+15550004444",
		ContactEmail:  "planner@example.com",
		EventType:     "corporate",
		EventDate:     eventDate,
		EventTime:     "18:30",
		DurationHours: 4,
		GuestCount:    25,
		VenueAddress:  "12 Admiralty Way",
		VenueArea:     "Lekki",
		VenueState:    "Lagos",
		ServiceStyle:  models.ServiceStyleBuffet,
		MenuNotes:     &menuNotes,
		Budget:        &budget,
	}

	mock.ExpectQuery("INSERT INTO catering_requests").WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(501),
	)

	createdRequest, err := svc.CreateRequest(input, nil)
	if err != nil {
		t.Fatalf("CreateRequest() error = %v", err)
	}
	if createdRequest == nil || createdRequest.ID != 501 {
		t.Fatalf("expected persisted catering request ID 501, got %+v", createdRequest)
	}
	if createdRequest.ContactName != "Planner" || createdRequest.GuestCount != 25 {
		t.Fatalf("persisted request fields mismatch: %+v", createdRequest)
	}
	if createdRequest.Status != models.CateringStatusPending {
		t.Fatalf("expected pending status, got %s", createdRequest.Status)
	}

	waitForSignal(t, notifier.cateringCalled, "catering confirmation notification attempt")
	waitForLogContains(t, logs, "Failed to send catering confirmation")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations failed: %v", err)
	}
}
