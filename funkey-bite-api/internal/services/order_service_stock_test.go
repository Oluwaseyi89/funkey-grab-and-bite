package services

import (
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

	svc := NewOrderService(nil, menu, stock, &noopNotificationService{})

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

	_, err := svc.CreateOrder(order, items)
	if err == nil {
		t.Fatal("expected CreateOrder to fail when requested quantity exceeds available stock")
	}

	if got := stock.stock[1]; got != 2 {
		t.Fatalf("stock changed on failed order: got %d want %d", got, 2)
	}
}

func TestCreateOrderDeductsInventoryOnSuccess(t *testing.T) {
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

	svc := NewOrderService(orderRepo, menu, stock, &noopNotificationService{})

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
			"pickup_time", "created_at",
		}).AddRow(
			101, "FG-2026-04-1002", nil, "Test User", "+15550000002", nil,
			string(models.OrderTypeDelivery), string(models.OrderStatusPending), 37.5, nil, nil, createdAt,
		))
	mock.ExpectQuery("SELECT id, order_id, menu_item_id, name, quantity, unit_price, special_instructions").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "order_id", "menu_item_id", "name", "quantity", "unit_price", "special_instructions",
		}).AddRow(201, 101, 1, "Burger", 3, 12.5, nil))

	_, err = svc.CreateOrder(order, items)
	if err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}

	if got := stock.stock[1]; got != 2 {
		t.Fatalf("stock was not deducted after successful order: got %d want %d", got, 2)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations failed: %v", err)
	}
}
