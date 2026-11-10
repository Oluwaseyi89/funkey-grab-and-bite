package v1

import (
	"encoding/json"
	"net/http"
	"testing"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"

	"github.com/gin-gonic/gin"
)

func userOrderRouter(h *OrderHandler, userID int) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/auth/orders/:id", func(c *gin.Context) { c.Set("user_id", userID) }, h.GetUserOrder)
	return router
}

func TestGetUserOrderReturnsOwnOrderWithInvoice(t *testing.T) {
	owner := 4
	invoiceNumber := "INV-FG-9"
	order := &models.Order{ID: 9, OrderNumber: "FG-9", UserID: &owner, TotalAmount: 2500,
		PaymentMethod: models.PaymentMethodCash, InvoiceNumber: &invoiceNumber}
	orders := &fakeOrderService{byID: func(id int) *models.Order {
		if id == 9 {
			return order
		}
		return nil
	}}
	handler := NewOrderHandler(orders, &fakeAuthService{}, acceptingSettings(), &fakePromotionService{}, &fakePaymentService{}, &fakeInvoiceService{})

	rec := call(userOrderRouter(handler, owner), http.MethodGet, "/auth/orders/9", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("the owner can open their order, got %d: %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Order   models.Order    `json:"order"`
		Invoice *models.Invoice `json:"invoice"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Order.OrderNumber != "FG-9" || body.Invoice == nil || body.Invoice.InvoiceNumber != invoiceNumber {
		t.Fatalf("expected the order and its invoice, got %s", rec.Body.String())
	}
}

func TestGetUserOrderHidesOtherCustomersOrders(t *testing.T) {
	owner := 4
	mine := &models.Order{ID: 9, OrderNumber: "FG-9", UserID: &owner}
	guest := &models.Order{ID: 10, OrderNumber: "FG-10"}
	orders := &fakeOrderService{byID: func(id int) *models.Order {
		switch id {
		case 9:
			return mine
		case 10:
			return guest
		}
		return nil
	}}
	handler := NewOrderHandler(orders, &fakeAuthService{}, acceptingSettings(), &fakePromotionService{}, &fakePaymentService{}, &fakeInvoiceService{})

	// Another customer's order, an unlinked guest order and a missing one all
	// answer the same, so order IDs can't be probed.
	for _, tc := range []struct {
		user int
		path string
	}{{5, "/auth/orders/9"}, {owner, "/auth/orders/10"}, {owner, "/auth/orders/11"}} {
		if rec := call(userOrderRouter(handler, tc.user), http.MethodGet, tc.path, ""); rec.Code != http.StatusNotFound {
			t.Fatalf("user %d opening %s: expected 404, got %d", tc.user, tc.path, rec.Code)
		}
	}
	if rec := call(userOrderRouter(handler, owner), http.MethodGet, "/auth/orders/abc", ""); rec.Code != http.StatusBadRequest {
		t.Fatalf("a non-numeric ID is a bad request, got %d", rec.Code)
	}
}
