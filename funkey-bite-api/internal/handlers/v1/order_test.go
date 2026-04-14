package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"

	"github.com/gin-gonic/gin"
)

type fakeOrderService struct {
	subtotal     float64
	calculateErr error
	createErr    error
	captured     *models.Order
	response     *models.Order
}

func (f *fakeOrderService) CreateOrder(order *models.Order, items []models.OrderItemRequest) (*models.Order, error) {
	if f.createErr != nil {
		return nil, f.createErr
	}
	copyOrder := *order
	f.captured = &copyOrder
	if f.response != nil {
		return f.response, nil
	}
	return order, nil
}

func (f *fakeOrderService) GetOrderByID(id int) (*models.Order, error) { return nil, nil }
func (f *fakeOrderService) GetOrdersByUserID(userID int) ([]models.Order, error) {
	return nil, nil
}
func (f *fakeOrderService) UpdateOrderStatus(id int, status string) error { return nil }
func (f *fakeOrderService) CalculateOrderTotal(items []models.OrderItemRequest) (float64, error) {
	if f.calculateErr != nil {
		return 0, f.calculateErr
	}
	return f.subtotal, nil
}
func (f *fakeOrderService) GetOrderByOrderNumber(orderNumber string) (*models.Order, error) {
	return nil, nil
}
func (f *fakeOrderService) GetOrderByPhoneAndOrderNumber(phone, orderNumber string) (*models.Order, error) {
	return nil, nil
}
func (f *fakeOrderService) CancelOrder(id int, userID int) error { return nil }

type fakeAuthService struct{}

func (f *fakeAuthService) Register(userData models.UserRegistration) (*models.AuthResponse, error) {
	return nil, nil
}
func (f *fakeAuthService) Login(phone, password string) (*models.AuthResponse, error) {
	return nil, nil
}
func (f *fakeAuthService) CheckUserExists(phone, email string) (*models.User, bool, error) {
	return nil, false, nil
}
func (f *fakeAuthService) AuthenticateOrder(orderData models.OrderWithAuth) (*models.User, error) {
	uid := 99
	return &models.User{ID: uid}, nil
}

type fakeSettingsService struct {
	minOrderOK   bool
	minOrderMsg  string
	orderTimeOK  bool
	orderTimeMsg string
	acceptOK     bool
	acceptMsg    string
	deliveryFee  float64
	taxAmount    float64
	total        float64
}

func (f *fakeSettingsService) GetSettings() (*models.BusinessSettings, error) { return nil, nil }
func (f *fakeSettingsService) UpdateSettings(updates *models.BusinessSettingsUpdate) (*models.BusinessSettings, error) {
	return nil, nil
}
func (f *fakeSettingsService) GetOpeningHours() ([]models.OpeningHours, error) { return nil, nil }
func (f *fakeSettingsService) IsBusinessOpen() (bool, error)                   { return true, nil }
func (f *fakeSettingsService) ValidateOrderTime(orderType string, requestedTime *time.Time) (bool, string) {
	return f.orderTimeOK, f.orderTimeMsg
}
func (f *fakeSettingsService) CanAcceptOrders(orderType string) (bool, string) {
	return f.acceptOK, f.acceptMsg
}
func (f *fakeSettingsService) CalculateOrderFees(subtotal float64, orderType string) (deliveryFee, taxAmount, total float64) {
	return f.deliveryFee, f.taxAmount, f.total
}
func (f *fakeSettingsService) ValidateMinimumOrder(subtotal float64, orderType string) (bool, string) {
	return f.minOrderOK, f.minOrderMsg
}

func newOrderRequestBody(t *testing.T, orderType models.OrderType) []byte {
	t.Helper()

	req := models.OrderWithAuth{
		CustomerName:  "Test User",
		CustomerPhone: "+15550001111",
		OrderType:     orderType,
		Items: []models.OrderItemRequest{{
			MenuItemID: 1,
			Name:       "Burger",
			Quantity:   1,
			UnitPrice:  8,
		}},
	}

	body, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}
	return body
}

func runCreateOrder(t *testing.T, h *OrderHandler, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/orders", func(c *gin.Context) {
		c.Set("user_id", 99)
		h.CreateOrder(c)
	})

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, req)

	return recorder
}

func TestCreateOrderRejectsBelowMinimumOrder(t *testing.T) {
	orderSvc := &fakeOrderService{subtotal: 8}
	settingsSvc := &fakeSettingsService{
		minOrderOK:  false,
		minOrderMsg: "Minimum order amount is $10.00",
		orderTimeOK: true,
		acceptOK:    true,
		total:       8,
	}

	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, settingsSvc)
	rec := runCreateOrder(t, handler, newOrderRequestBody(t, models.OrderTypePickup))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Minimum order amount") {
		t.Fatalf("expected minimum order error, got body: %s", rec.Body.String())
	}
	if orderSvc.captured != nil {
		t.Fatal("CreateOrder should not be called when minimum order validation fails")
	}
}

func TestCreateOrderRejectsOutsideBusinessHours(t *testing.T) {
	orderSvc := &fakeOrderService{subtotal: 18}
	settingsSvc := &fakeSettingsService{
		minOrderOK:   true,
		orderTimeOK:  false,
		orderTimeMsg: "Business is closed at 23:30. Hours: 08:00 - 22:00",
		acceptOK:     true,
		total:        19.44,
	}

	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, settingsSvc)
	rec := runCreateOrder(t, handler, newOrderRequestBody(t, models.OrderTypePickup))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Business is closed") {
		t.Fatalf("expected business-hours error, got body: %s", rec.Body.String())
	}
	if orderSvc.captured != nil {
		t.Fatal("CreateOrder should not be called when order-time validation fails")
	}
}

func TestCreateOrderRejectsWhenDeliveryIsDisabled(t *testing.T) {
	orderSvc := &fakeOrderService{subtotal: 22}
	settingsSvc := &fakeSettingsService{
		minOrderOK:  true,
		orderTimeOK: true,
		acceptOK:    false,
		acceptMsg:   "Delivery service is currently unavailable",
		total:       26.87,
	}

	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, settingsSvc)
	rec := runCreateOrder(t, handler, newOrderRequestBody(t, models.OrderTypeDelivery))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Delivery service is currently unavailable") {
		t.Fatalf("expected delivery-disabled error, got body: %s", rec.Body.String())
	}
	if orderSvc.captured != nil {
		t.Fatal("CreateOrder should not be called when order acceptance validation fails")
	}
}

func TestCreateOrderPersistsTotalWithDeliveryFeeAndTax(t *testing.T) {
	orderSvc := &fakeOrderService{
		subtotal: 20,
		response: &models.Order{ID: 1, OrderNumber: "FG-2026-04-0001", TotalAmount: 23.6},
	}
	settingsSvc := &fakeSettingsService{
		minOrderOK:  true,
		orderTimeOK: true,
		acceptOK:    true,
		deliveryFee: 2,
		taxAmount:   1.6,
		total:       23.6,
	}

	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, settingsSvc)
	rec := runCreateOrder(t, handler, newOrderRequestBody(t, models.OrderTypeDelivery))

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}
	if orderSvc.captured == nil {
		t.Fatal("expected CreateOrder to be called")
	}
	if orderSvc.captured.TotalAmount != 23.6 {
		t.Fatalf("expected persisted total amount 23.6, got %v", orderSvc.captured.TotalAmount)
	}
}
