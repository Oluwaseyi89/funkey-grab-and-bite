package services

import (
	"context"
	"database/sql"
	"sync"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/repository"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

type fakeMenuReader struct {
	items map[int]*models.MenuItem
}

func (f *fakeMenuReader) GetByID(id int) (*models.MenuItem, error) {
	if item, ok := f.items[id]; ok {
		return item, nil
	}
	return nil, nil
}

type fakeStockService struct {
	mu    sync.Mutex
	stock map[int]int
}

func (f *fakeStockService) CheckAvailability(menuItemID, quantity int) (bool, string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	current := f.stock[menuItemID]
	if current < quantity {
		return false, "Insufficient stock"
	}
	return true, ""
}

func (f *fakeStockService) DeductStockForOrder(menuItemID, quantity int, orderID int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.stock[menuItemID] < quantity {
		return sql.ErrNoRows
	}
	f.stock[menuItemID] -= quantity
	return nil
}

func (f *fakeStockService) RestockItem(menuItemID, quantity int, reason string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.stock[menuItemID] += quantity
	return nil
}

// fakeCheckoutStore stands in for the transaction that saves an order.
type fakeCheckoutStore struct {
	orderID int
	err     error
	placed  *models.Checkout
}

func (f *fakeCheckoutStore) PlaceOrder(ctx context.Context, checkout *models.Checkout) (int, error) {
	f.placed = checkout
	return f.orderID, f.err
}

type noopNotificationService struct{}

func (n *noopNotificationService) SendOrderConfirmation(order *models.Order) error { return nil }
func (n *noopNotificationService) SendOrderStatusUpdate(orderID int, newStatus string) error {
	return nil
}
func (n *noopNotificationService) SendCateringConfirmation(request *models.CateringRequest) error {
	return nil
}
func (n *noopNotificationService) SendPasswordReset(email, resetToken string) error { return nil }
func (n *noopNotificationService) SendPhoneVerification(phoneNumber, code string) error {
	return nil
}
func (n *noopNotificationService) GetNotificationHistory(userID int, limit int) ([]models.Notification, error) {
	return nil, nil
}
func (n *noopNotificationService) MarkNotificationAsRead(notificationID int) error { return nil }

func TestCreateOrderRejectsWhenInventoryIsInsufficientAndLeavesStockUnchanged(t *testing.T) {
	menu := &fakeMenuReader{items: map[int]*models.MenuItem{
		1: {ID: 1, Name: "Burger", Price: 12.5, IsAvailable: true},
	}}
	stock := &fakeStockService{stock: map[int]int{1: 2}}

	checkout := &fakeCheckoutStore{}
	svc := NewOrderService(nil, menu, stock, &noopNotificationService{}, checkout)

	order := &models.Order{
		OrderNumber:   "FG-2026-04-1001",
		CustomerName:  "Test User",
		CustomerPhone: "+15550000001",
		OrderType:     models.OrderTypeDelivery,
		Status:        models.OrderStatusPending,
		TotalAmount:   37.5,
		CreatedAt:     time.Now(),
	}
	items := []models.OrderItemRequest{{
		MenuItemID: 1,
		Name:       "Burger",
		Quantity:   3,
		UnitPrice:  12.5,
	}}

	_, err := svc.CreateOrder(context.Background(), &models.Checkout{Order: order, Items: items})
	if err == nil {
		t.Fatal("expected CreateOrder to fail when requested quantity exceeds available stock")
	}

	if checkout.placed != nil {
		t.Fatal("an order that fails the stock check must not be saved")
	}
}

func TestCreateOrderSavesCheckoutOnSuccess(t *testing.T) {
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

	checkout := &fakeCheckoutStore{orderID: 101}
	svc := NewOrderService(orderRepo, menu, stock, &noopNotificationService{}, checkout)

	createdAt := time.Now()
	order := &models.Order{
		OrderNumber:   "FG-2026-04-1002",
		CustomerName:  "Test User",
		CustomerPhone: "+15550000002",
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
			101, "FG-2026-04-1002", nil, "Test User", "+15550000002", nil,
			string(models.OrderTypeDelivery), string(models.OrderStatusPending), 37.5, nil, nil, createdAt,
			string(models.PaymentMethodCash), string(models.PaymentStatusNotRequired),
			nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil, nil, nil,
		))
	mock.ExpectQuery("SELECT id, order_id, menu_item_id, name, quantity, unit_price, special_instructions").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "order_id", "menu_item_id", "name", "quantity", "unit_price", "special_instructions",
		}).AddRow(201, 101, 1, "Burger", 3, 12.5, nil))

	_, err = svc.CreateOrder(context.Background(), &models.Checkout{Order: order, Items: items})
	if err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}

	if checkout.placed == nil || checkout.placed.Order != order || len(checkout.placed.Items) != 1 {
		t.Fatalf("expected the order and its items to be saved, got %+v", checkout.placed)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations failed: %v", err)
	}
}

func TestCreateOrderReportsStockAndPromotionThatRanOutWhileSaving(t *testing.T) {
	menu := &fakeMenuReader{items: map[int]*models.MenuItem{
		1: {ID: 1, Name: "Burger", Price: 12.5, IsAvailable: true},
	}}
	items := []models.OrderItemRequest{{MenuItemID: 1, Name: "Burger", Quantity: 3, UnitPrice: 12.5}}

	cases := map[string]struct {
		err  error
		want string
	}{
		"some left":   {&repository.InsufficientStockError{MenuItemID: 1, Available: 1, Requested: 3}, "Burger has only 1 left"},
		"none left":   {&repository.InsufficientStockError{MenuItemID: 1, Requested: 3}, "Burger is out of stock"},
		"not tracked": {&repository.InsufficientStockError{MenuItemID: 1, Requested: 3, Untracked: true}, "Burger is not available right now"},
		"promotion":   {repository.ErrPromotionUnavailable, ErrPromotionUnavailable.Error()},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			stock := &fakeStockService{stock: map[int]int{1: 5}}
			svc := NewOrderService(nil, menu, stock, &noopNotificationService{}, &fakeCheckoutStore{err: tc.err})

			_, err := svc.CreateOrder(context.Background(), &models.Checkout{Order: &models.Order{}, Items: items})
			if err == nil || err.Error() != tc.want {
				t.Fatalf("expected %q, got %v", tc.want, err)
			}
		})
	}
}
