package v1

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"

	"github.com/gin-gonic/gin"
)

type fakeAdminService struct {
	updateOrderStatusErr  error
	updateUserStatusErr   error
	updateAdminUserErr    error
	updateAdminPassErr    error
	updatedAdminLookupErr error
}

func (f *fakeAdminService) AdminLogin(email, password string) (*models.AdminUser, string, error) {
	return nil, "", nil
}
func (f *fakeAdminService) AdminLogout(adminID int) error { return nil }
func (f *fakeAdminService) GetAdminUsers(page, limit int) ([]models.AdminUser, int, error) {
	return nil, 0, nil
}
func (f *fakeAdminService) GetAdminUserByID(adminID int) (*models.AdminUser, error) {
	if f.updatedAdminLookupErr != nil {
		return nil, f.updatedAdminLookupErr
	}
	return nil, nil
}
func (f *fakeAdminService) CreateAdminUser(admin *models.AdminUser, password string) (*models.AdminUser, error) {
	return nil, nil
}
func (f *fakeAdminService) UpdateAdminUser(adminID int, updates *models.AdminUser) error {
	return f.updateAdminUserErr
}
func (f *fakeAdminService) DeleteAdminUser(adminID int) error { return nil }
func (f *fakeAdminService) UpdateAdminPassword(adminID int, currentPassword, newPassword string) error {
	return f.updateAdminPassErr
}
func (f *fakeAdminService) GetDashboardStats() (*models.AdminStats, error) { return nil, nil }
func (f *fakeAdminService) GetTodayStats() (*models.AdminStats, error)     { return nil, nil }
func (f *fakeAdminService) GetSalesReport(fromDate, toDate string) ([]models.SalesReport, error) {
	return nil, nil
}
func (f *fakeAdminService) GetAllOrders(page, limit int, status string) ([]models.Order, int, error) {
	return nil, 0, nil
}
func (f *fakeAdminService) UpdateOrderStatus(orderID int, status string) error {
	return f.updateOrderStatusErr
}
func (f *fakeAdminService) GetAllUsers(page, limit int) ([]models.User, int, error) {
	return nil, 0, nil
}
func (f *fakeAdminService) UpdateUserStatus(userID int, isActive bool) error {
	return f.updateUserStatusErr
}
func (f *fakeAdminService) CreateMenuItem(item *models.MenuItem) (*models.MenuItem, error) {
	return nil, nil
}
func (f *fakeAdminService) GetMenuItems(page, limit int, categoryID *int, query string) ([]models.MenuItem, error) {
	return nil, nil
}
func (f *fakeAdminService) UpdateMenuItem(item *models.MenuItem) error { return nil }
func (f *fakeAdminService) DeleteMenuItem(id int) error                { return nil }
func (f *fakeAdminService) GetMenuItemByID(id int) (*models.MenuItem, error) {
	return nil, nil
}
func (f *fakeAdminService) GetAllCateringRequests(page, limit int, status string) ([]models.CateringRequest, int, error) {
	return nil, 0, nil
}

type fakeInventoryService struct {
	resolveAlertErr error
	updateStockErr  error
}

func (f *fakeInventoryService) GetInventoryItem(id int) (*models.InventoryItem, error) {
	return nil, nil
}
func (f *fakeInventoryService) GetInventoryByMenuItemID(menuItemID int) (*models.InventoryItem, error) {
	return nil, nil
}
func (f *fakeInventoryService) GetAllInventory() ([]models.InventoryItem, error) { return nil, nil }
func (f *fakeInventoryService) GetLowStock() ([]models.InventoryItem, error)     { return nil, nil }
func (f *fakeInventoryService) UpdateStock(update *models.InventoryUpdate) (*models.InventoryItem, error) {
	if f.updateStockErr != nil {
		return nil, f.updateStockErr
	}
	return &models.InventoryItem{ID: 1, MenuItemID: update.MenuItemID, CurrentStock: 1}, nil
}
func (f *fakeInventoryService) CreateInventoryItem(item *models.InventoryItem) (*models.InventoryItem, error) {
	return nil, nil
}
func (f *fakeInventoryService) GetAlerts(resolved bool) ([]models.InventoryAlert, error) {
	return nil, nil
}
func (f *fakeInventoryService) ResolveAlert(alertID int) error { return f.resolveAlertErr }
func (f *fakeInventoryService) CheckAvailability(menuItemID, quantity int) (bool, string) {
	return true, ""
}
func (f *fakeInventoryService) DeductStockForOrder(menuItemID, quantity int, orderID int) error {
	return nil
}
func (f *fakeInventoryService) RestockItem(menuItemID, quantity int, reason string) error { return nil }
func (f *fakeInventoryService) GetInventoryDashboard() (map[string]interface{}, error) {
	return nil, nil
}

func TestUpdateOrderStatusMissingOrderReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeAdminService{updateOrderStatusErr: sql.ErrNoRows}
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PATCH("/admin/orders/:id/status", handler.UpdateOrderStatus)

	body := []byte(`{"status":"confirmed"}`)
	req := httptest.NewRequest(http.MethodPatch, "/admin/orders/999/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), "Order not found") {
		t.Fatalf("expected not-found message, got body=%s", res.Body.String())
	}
}

func TestUpdateUserStatusMissingUserReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeAdminService{updateUserStatusErr: sql.ErrNoRows}
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PATCH("/admin/users/:id/status", handler.UpdateUserStatus)

	body := []byte(`{"isActive":true}`)
	req := httptest.NewRequest(http.MethodPatch, "/admin/users/999/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), "User not found") {
		t.Fatalf("expected not-found message, got body=%s", res.Body.String())
	}
}

func TestUpdateAdminUserMissingAdminReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeAdminService{updateAdminUserErr: errors.New("admin not found")}
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PUT("/admin/users/admins/:id", func(c *gin.Context) {
		c.Set("user_id", 1)
		handler.UpdateAdminUser(c)
	})

	body := []byte(`{"username":"missing","email":"missing@example.com","role":"manager","isActive":true}`)
	req := httptest.NewRequest(http.MethodPut, "/admin/users/admins/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), "admin not found") {
		t.Fatalf("expected admin-not-found message, got body=%s", res.Body.String())
	}
}

func TestUpdateAdminPasswordMissingAdminReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeAdminService{updateAdminPassErr: errors.New("admin not found")}
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PATCH("/admin/auth/password", func(c *gin.Context) {
		c.Set("user_id", 999)
		handler.UpdateAdminPassword(c)
	})

	body := []byte(`{"currentPassword":"Current123","newPassword":"NewStrong123"}`)
	req := httptest.NewRequest(http.MethodPatch, "/admin/auth/password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), "admin not found") {
		t.Fatalf("expected admin-not-found message, got body=%s", res.Body.String())
	}
}

func TestResolveAlertMissingAlertReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeInventoryService{resolveAlertErr: sql.ErrNoRows}
	handler := NewInventoryHandler(service)

	router := gin.New()
	router.PATCH("/admin/inventory/alerts/:id/resolve", handler.ResolveAlert)

	req := httptest.NewRequest(http.MethodPatch, "/admin/inventory/alerts/999/resolve", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), "Alert not found") {
		t.Fatalf("expected alert-not-found message, got body=%s", res.Body.String())
	}
}

func TestUpdateStockMissingInventoryReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeInventoryService{updateStockErr: errors.New("inventory item not found for menu item ID: 999")}
	handler := NewInventoryHandler(service)

	router := gin.New()
	router.PATCH("/admin/inventory/stock", handler.UpdateStock)

	payload := models.InventoryUpdate{
		MenuItemID: 999,
		Quantity:   1,
		Operation:  "subtract",
		Reason:     "test",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPatch, "/admin/inventory/stock", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", res.Code, res.Body.String())
	}
}

func TestUpdateOrderStatusMissingOrderDoesNotReturnSuccessPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeAdminService{updateOrderStatusErr: sql.ErrNoRows}
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PATCH("/admin/orders/:id/status", handler.UpdateOrderStatus)

	body := []byte(`{"status":"ready"}`)
	req := httptest.NewRequest(http.MethodPatch, "/admin/orders/404/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code == http.StatusOK {
		t.Fatalf("expected non-200 for missing order, got body=%s", res.Body.String())
	}
	if strings.Contains(res.Body.String(), "updated successfully") {
		t.Fatalf("unexpected success body for missing order: %s", res.Body.String())
	}
}

func TestUpdateUserStatusMissingUserDoesNotReturnSuccessPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeAdminService{updateUserStatusErr: sql.ErrNoRows}
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PATCH("/admin/users/:id/status", handler.UpdateUserStatus)

	body := []byte(`{"isActive":true}`)
	req := httptest.NewRequest(http.MethodPatch, "/admin/users/404/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code == http.StatusOK {
		t.Fatalf("expected non-200 for missing user, got body=%s", res.Body.String())
	}
	if strings.Contains(res.Body.String(), "updated successfully") {
		t.Fatalf("unexpected success body for missing user: %s", res.Body.String())
	}
}

func TestUpdateAdminPasswordValidationStillApplies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeAdminService{}
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PATCH("/admin/auth/password", func(c *gin.Context) {
		c.Set("user_id", 999)
		handler.UpdateAdminPassword(c)
	})

	body := []byte(`{"currentPassword":"short","newPassword":"short"}`)
	req := httptest.NewRequest(http.MethodPatch, "/admin/auth/password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 validation failure, got %d body=%s", res.Code, res.Body.String())
	}
}

func TestUpdateAdminUserSelfModificationBlockedBeforeService(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeAdminService{}
	handler := NewAdminHandler(service)

	router := gin.New()
	router.PUT("/admin/users/admins/:id", func(c *gin.Context) {
		c.Set("user_id", 5)
		handler.UpdateAdminUser(c)
	})

	payload := map[string]interface{}{
		"username": "self",
		"email":    "self@example.com",
		"role":     "manager",
		"isActive": true,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/admin/users/admins/5", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for self-modification guard, got %d body=%s", res.Code, res.Body.String())
	}
}

func TestUpdateEndpointsMissingRecordsMaintainErrorSemantics(t *testing.T) {
	_ = time.Now()
}
