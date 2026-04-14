package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
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
	trackByPhone *models.Order
	trackErr     error
	history      []models.Order
	historyErr   error
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
	if f.historyErr != nil {
		return nil, f.historyErr
	}
	return f.history, nil
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
	if f.trackErr != nil {
		return nil, f.trackErr
	}
	return f.trackByPhone, nil
}
func (f *fakeOrderService) CancelOrder(id int, userID int) error { return nil }

type fakeAuthService struct {
	authUser *models.User
	authErr  error
}

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
	if f.authErr != nil {
		return nil, f.authErr
	}
	if f.authUser != nil {
		return f.authUser, nil
	}
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

type fakePromotionService struct {
	validation     *models.PromotionValidation
	validateErr    error
	applyDiscount  float64
	applyErr       error
	appliedCode    string
	appliedOrderID int
}

func (f *fakePromotionService) CreatePromotion(promotion *models.PromotionCreate) (*models.Promotion, error) {
	return nil, nil
}
func (f *fakePromotionService) GetPromotionByID(id int) (*models.Promotion, error) {
	return nil, nil
}
func (f *fakePromotionService) GetPromotionByCode(code string) (*models.Promotion, error) {
	return nil, nil
}
func (f *fakePromotionService) GetAllPromotions(page, limit int, status string) ([]models.Promotion, int, error) {
	return nil, 0, nil
}
func (f *fakePromotionService) UpdatePromotion(id int, updates *models.PromotionUpdate) (*models.Promotion, error) {
	return nil, nil
}
func (f *fakePromotionService) DeletePromotion(id int) error { return nil }
func (f *fakePromotionService) ValidatePromotion(code string, orderAmount float64, customerID *int) (*models.PromotionValidation, error) {
	if f.validateErr != nil {
		return nil, f.validateErr
	}
	if f.validation != nil {
		return f.validation, nil
	}
	return &models.PromotionValidation{IsValid: false, Message: "invalid promotion"}, nil
}
func (f *fakePromotionService) ApplyPromotion(promotionID, orderID int, customerID *int, orderAmount float64) (float64, error) {
	return f.applyDiscount, f.applyErr
}
func (f *fakePromotionService) ApplyPromotionByCode(code string, orderID int, customerID *int, orderAmount float64) (float64, error) {
	f.appliedCode = code
	f.appliedOrderID = orderID
	return f.applyDiscount, f.applyErr
}
func (f *fakePromotionService) GetActivePromotions() ([]models.Promotion, error) { return nil, nil }

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

func runCreateOrderUnauthenticated(t *testing.T, h *OrderHandler, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/orders", h.CreateOrder)

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

	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, settingsSvc, &fakePromotionService{})
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

	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, settingsSvc, &fakePromotionService{})
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

	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, settingsSvc, &fakePromotionService{})
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

	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, settingsSvc, &fakePromotionService{})
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

func TestCreateOrderAppliesPromotionDiscountToPersistedTotal(t *testing.T) {
	orderSvc := &fakeOrderService{
		subtotal: 20,
		response: &models.Order{ID: 41, OrderNumber: "FG-2026-04-0041", TotalAmount: 18.6},
	}
	settingsSvc := &fakeSettingsService{
		minOrderOK:  true,
		orderTimeOK: true,
		acceptOK:    true,
		deliveryFee: 2,
		taxAmount:   1.6,
		total:       23.6,
	}
	promoSvc := &fakePromotionService{
		validation: &models.PromotionValidation{IsValid: true, PromotionID: 99, Discount: 5},
	}

	code := "SAVE5"
	req := models.OrderWithAuth{
		CustomerName:  "Test User",
		CustomerPhone: "+15550001111",
		OrderType:     models.OrderTypeDelivery,
		PromotionCode: &code,
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

	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, settingsSvc, promoSvc)
	rec := runCreateOrder(t, handler, body)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}
	if orderSvc.captured == nil {
		t.Fatal("expected CreateOrder to be called")
	}
	if orderSvc.captured.TotalAmount != 18.6 {
		t.Fatalf("expected discounted total amount 18.6, got %v", orderSvc.captured.TotalAmount)
	}
	if promoSvc.appliedCode != "SAVE5" {
		t.Fatalf("expected promotion code SAVE5 to be applied, got %q", promoSvc.appliedCode)
	}
}

func TestCreateOrderRejectsInvalidPromotion(t *testing.T) {
	orderSvc := &fakeOrderService{subtotal: 20}
	settingsSvc := &fakeSettingsService{
		minOrderOK:  true,
		orderTimeOK: true,
		acceptOK:    true,
		total:       23.6,
	}
	promoSvc := &fakePromotionService{
		validation: &models.PromotionValidation{IsValid: false, Message: "Promotion code not found"},
	}

	code := "BADCODE"
	req := models.OrderWithAuth{
		CustomerName:  "Test User",
		CustomerPhone: "+15550001111",
		OrderType:     models.OrderTypeDelivery,
		PromotionCode: &code,
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

	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, settingsSvc, promoSvc)
	rec := runCreateOrder(t, handler, body)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Promotion code not found") {
		t.Fatalf("expected promotion validation message, got body: %s", rec.Body.String())
	}
	if orderSvc.captured != nil {
		t.Fatal("CreateOrder should not be called when promotion validation fails")
	}
}

func TestCreateOrderFailsWhenPromotionApplyFailsAfterCreate(t *testing.T) {
	orderSvc := &fakeOrderService{
		subtotal: 20,
		response: &models.Order{ID: 77, OrderNumber: "FG-2026-04-0077", TotalAmount: 20.6},
	}
	settingsSvc := &fakeSettingsService{
		minOrderOK:  true,
		orderTimeOK: true,
		acceptOK:    true,
		total:       23.6,
	}
	promoSvc := &fakePromotionService{
		validation:    &models.PromotionValidation{IsValid: true, PromotionID: 199, Discount: 3},
		applyErr:      fmt.Errorf("promotion usage limit reached"),
		applyDiscount: 3,
	}

	code := "FLASH3"
	req := models.OrderWithAuth{
		CustomerName:  "Test User",
		CustomerPhone: "+15550001111",
		OrderType:     models.OrderTypeDelivery,
		PromotionCode: &code,
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

	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, settingsSvc, promoSvc)
	rec := runCreateOrder(t, handler, body)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusInternalServerError, rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Failed to apply promotion") {
		t.Fatalf("expected apply promotion error, got body: %s", rec.Body.String())
	}
}

func TestCreateOrderUnauthenticatedWithNewPhoneRequiresAccount(t *testing.T) {
	orderSvc := &fakeOrderService{subtotal: 20}
	settingsSvc := &fakeSettingsService{
		minOrderOK:  true,
		orderTimeOK: true,
		acceptOK:    true,
		total:       23.6,
	}
	authSvc := &fakeAuthService{authErr: errors.New("user_not_found")}

	handler := NewOrderHandler(orderSvc, authSvc, settingsSvc, &fakePromotionService{})
	rec := runCreateOrderUnauthenticated(t, handler, newOrderRequestBody(t, models.OrderTypeDelivery))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "User not found. Please create an account first.") {
		t.Fatalf("expected new-user account-required message, got body: %s", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "USER_NOT_FOUND") {
		t.Fatalf("expected USER_NOT_FOUND code, got body: %s", rec.Body.String())
	}
	if orderSvc.captured != nil {
		t.Fatal("CreateOrder should not be called for unauthenticated new-phone checkout")
	}
}

func TestCreateOrderUnauthenticatedExistingPhoneWithoutPasswordIsRejected(t *testing.T) {
	orderSvc := &fakeOrderService{subtotal: 20}
	settingsSvc := &fakeSettingsService{
		minOrderOK:  true,
		orderTimeOK: true,
		acceptOK:    true,
		total:       23.6,
	}
	authSvc := &fakeAuthService{authErr: errors.New("password_required")}

	handler := NewOrderHandler(orderSvc, authSvc, settingsSvc, &fakePromotionService{})
	rec := runCreateOrderUnauthenticated(t, handler, newOrderRequestBody(t, models.OrderTypeDelivery))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusUnauthorized, rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Password required for existing user") {
		t.Fatalf("expected password-required message, got body: %s", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "PASSWORD_REQUIRED") {
		t.Fatalf("expected PASSWORD_REQUIRED code, got body: %s", rec.Body.String())
	}
	if orderSvc.captured != nil {
		t.Fatal("CreateOrder should not be called when existing user omits password")
	}
}

func TestTrackOrderPublicReturnsLimitedFieldsOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)

	email := "customer@example.com"
	pickupTime := time.Now().Add(20 * time.Minute)
	orderSvc := &fakeOrderService{
		trackByPhone: &models.Order{
			OrderNumber:   "FG-2026-04-1234",
			Status:        models.OrderStatusPreparing,
			OrderType:     models.OrderTypeDelivery,
			TotalAmount:   49.99,
			CustomerName:  "Sensitive Name",
			CustomerPhone: "+15551234567",
			CustomerEmail: &email,
			PickupTime:    &pickupTime,
			Items: []models.OrderItem{{
				MenuItemID: 1,
				Name:       "Secret Item",
				Quantity:   2,
				UnitPrice:  24.99,
			}},
		},
	}

	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, &fakeSettingsService{}, &fakePromotionService{})
	router := gin.New()
	router.GET("/order/track/:phone/:orderNumber", handler.TrackOrderPublic)

	req := httptest.NewRequest(http.MethodGet, "/order/track/+15551234567/FG-2026-04-1234", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected %d got %d body=%s", http.StatusOK, res.Code, res.Body.String())
	}

	body := res.Body.String()
	if strings.Contains(body, "items") {
		t.Fatalf("public tracking response leaked items: %s", body)
	}
	if strings.Contains(body, "totalAmount") {
		t.Fatalf("public tracking response leaked total amount: %s", body)
	}
	if strings.Contains(body, "Sensitive Name") || strings.Contains(body, "customerPhone") || strings.Contains(body, "customerEmail") {
		t.Fatalf("public tracking response leaked customer details: %s", body)
	}
	if !strings.Contains(body, "Limited order information available") {
		t.Fatalf("expected limited-info message, got: %s", body)
	}
}

func TestTrackOrderPublicMissingOrderReturnsGenericNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	orderSvc := &fakeOrderService{trackByPhone: nil}
	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, &fakeSettingsService{}, &fakePromotionService{})
	router := gin.New()
	router.GET("/order/track/:phone/:orderNumber", handler.TrackOrderPublic)

	req := httptest.NewRequest(http.MethodGet, "/order/track/+15550000000/FG-404", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected %d got %d body=%s", http.StatusNotFound, res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), "Order not found") {
		t.Fatalf("expected generic not found message, got body=%s", res.Body.String())
	}
	if strings.Contains(res.Body.String(), "phone number") {
		t.Fatalf("response should not reveal phone mismatch detail: %s", res.Body.String())
	}
}

func buildLargeOrderHistory(totalOrders int, itemsPerOrder int) []models.Order {
	orders := make([]models.Order, 0, totalOrders)
	for i := 0; i < totalOrders; i++ {
		items := make([]models.OrderItem, 0, itemsPerOrder)
		for j := 0; j < itemsPerOrder; j++ {
			items = append(items, models.OrderItem{
				MenuItemID: j + 1,
				Name:       fmt.Sprintf("Item-%d-%d", i+1, j+1),
				Quantity:   1,
				UnitPrice:  9.99,
			})
		}

		orders = append(orders, models.Order{
			ID:            i + 1,
			OrderNumber:   fmt.Sprintf("FG-2026-04-%04d", i+1),
			CustomerName:  "History User",
			CustomerPhone: "+15550001111",
			OrderType:     models.OrderTypeDelivery,
			Status:        models.OrderStatusCompleted,
			TotalAmount:   49.95,
			CreatedAt:     time.Now().Add(-time.Duration(i) * time.Minute),
			Items:         items,
		})
	}
	return orders
}

func TestGetUserOrdersPaginationCorrectnessWithLargeHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	history := buildLargeOrderHistory(1000, 2)
	orderSvc := &fakeOrderService{history: history}
	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, &fakeSettingsService{}, &fakePromotionService{})

	router := gin.New()
	router.GET("/auth/orders", func(c *gin.Context) {
		c.Set("user_id", 99)
		handler.GetUserOrders(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/auth/orders?page=3&limit=25", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected %d got %d body=%s", http.StatusOK, res.Code, res.Body.String())
	}

	var payload struct {
		Orders     []models.Order `json:"orders"`
		Pagination struct {
			Page       int `json:"page"`
			Limit      int `json:"limit"`
			Total      int `json:"total"`
			TotalPages int `json:"totalPages"`
		} `json:"pagination"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(payload.Orders) != 25 {
		t.Fatalf("expected 25 orders on page, got %d", len(payload.Orders))
	}
	if payload.Pagination.Page != 3 || payload.Pagination.Limit != 25 {
		t.Fatalf("unexpected pagination metadata: %+v", payload.Pagination)
	}
	if payload.Pagination.Total != 1000 || payload.Pagination.TotalPages != 40 {
		t.Fatalf("unexpected totals: %+v", payload.Pagination)
	}

	offset := (3 - 1) * 25
	if payload.Orders[0].ID != history[offset].ID {
		t.Fatalf("first order on page mismatch: got ID=%d want ID=%d", payload.Orders[0].ID, history[offset].ID)
	}
	if payload.Orders[24].ID != history[offset+24].ID {
		t.Fatalf("last order on page mismatch: got ID=%d want ID=%d", payload.Orders[24].ID, history[offset+24].ID)
	}
}

func TestGetUserOrdersLargeHistoryLatencyAndMemoryUnderLoad(t *testing.T) {
	gin.SetMode(gin.TestMode)

	history := buildLargeOrderHistory(5000, 6)
	orderSvc := &fakeOrderService{history: history}
	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, &fakeSettingsService{}, &fakePromotionService{})

	router := gin.New()
	router.GET("/auth/orders", func(c *gin.Context) {
		c.Set("user_id", 99)
		handler.GetUserOrders(c)
	})

	runtime.GC()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)

	start := time.Now()
	req := httptest.NewRequest(http.MethodGet, "/auth/orders?page=1&limit=20", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	elapsed := time.Since(start)

	var after runtime.MemStats
	runtime.ReadMemStats(&after)
	allocDelta := after.Alloc - before.Alloc

	if res.Code != http.StatusOK {
		t.Fatalf("expected %d got %d body=%s", http.StatusOK, res.Code, res.Body.String())
	}

	var payload struct {
		Orders []models.Order `json:"orders"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(payload.Orders) != 20 {
		t.Fatalf("expected 20 orders on first page, got %d", len(payload.Orders))
	}

	// Guardrails to catch severe regressions under high-load histories.
	if elapsed > 2*time.Second {
		t.Fatalf("GetUserOrders exceeded latency budget: %v", elapsed)
	}
	if allocDelta > 150*1024*1024 {
		t.Fatalf("GetUserOrders exceeded memory budget: alloc delta %d bytes", allocDelta)
	}
}
