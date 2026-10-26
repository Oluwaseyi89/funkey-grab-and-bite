package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/services"

	"github.com/gin-gonic/gin"
)

func acceptingSettings() *fakeSettingsService {
	return &fakeSettingsService{minOrderOK: true, orderTimeOK: true, acceptOK: true, deliveryFee: 500, taxAmount: 75, total: 1575}
}

func orderBody(t *testing.T, mutate func(*models.OrderWithAuth)) []byte {
	t.Helper()
	var req models.OrderWithAuth
	if err := json.Unmarshal(newOrderRequestBody(t, models.OrderTypeDelivery), &req); err != nil {
		t.Fatal(err)
	}
	mutate(&req)
	body, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func TestCreateOrderGuestCheckoutLinksCustomerByPhone(t *testing.T) {
	orderSvc := &fakeOrderService{subtotal: 1000}
	authSvc := &fakeAuthService{authUser: &models.User{ID: 7}}
	handler := NewOrderHandler(orderSvc, authSvc, acceptingSettings(), &fakePromotionService{}, &fakePaymentService{}, &fakeInvoiceService{})

	body := orderBody(t, func(r *models.OrderWithAuth) { r.CustomerPhone = " +234 (801) 234-5678 " })
	rec := runCreateOrderUnauthenticated(t, handler, body)

	if rec.Code != http.StatusCreated {
		t.Fatalf("guest checkout must succeed, got %d: %s", rec.Code, rec.Body.String())
	}
	if authSvc.resolvedPhone != "+2348012345678" {
		t.Fatalf("expected cleaned phone to be resolved, got %q", authSvc.resolvedPhone)
	}
	if orderSvc.captured.UserID == nil || *orderSvc.captured.UserID != 7 {
		t.Fatalf("expected order linked to customer 7, got %v", orderSvc.captured.UserID)
	}
	if orderSvc.captured.CustomerPhone != "+2348012345678" {
		t.Fatalf("expected cleaned phone stored, got %q", orderSvc.captured.CustomerPhone)
	}
}

func TestCreateOrderGuestCheckoutSurvivesCustomerRecordFailure(t *testing.T) {
	orderSvc := &fakeOrderService{subtotal: 1000}
	authSvc := &fakeAuthService{authErr: errors.New("database unavailable")}
	handler := NewOrderHandler(orderSvc, authSvc, acceptingSettings(), &fakePromotionService{}, &fakePaymentService{}, &fakeInvoiceService{})

	rec := runCreateOrderUnauthenticated(t, handler, newOrderRequestBody(t, models.OrderTypePickup))

	if rec.Code != http.StatusCreated {
		t.Fatalf("a customer-record failure must not lose the order, got %d: %s", rec.Code, rec.Body.String())
	}
	if orderSvc.captured.UserID != nil {
		t.Fatalf("expected unlinked order, got user %d", *orderSvc.captured.UserID)
	}
}

func TestCreateOrderPayOnDeliveryRaisesInvoice(t *testing.T) {
	orderSvc := &fakeOrderService{subtotal: 1000}
	invoices := &fakeInvoiceService{sent: make(chan *models.Order, 1)}
	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, acceptingSettings(), &fakePromotionService{}, &fakePaymentService{}, invoices)

	rec := runCreateOrder(t, handler, newOrderRequestBody(t, models.OrderTypeDelivery))
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	order := orderSvc.captured
	if order.PaymentStatus != models.PaymentStatusPending {
		t.Fatalf("pay on delivery must await payment, got %q", order.PaymentStatus)
	}
	if order.InvoiceNumber == nil || *order.InvoiceNumber != services.InvoiceNumberFor(order.OrderNumber) {
		t.Fatalf("expected invoice number for %s, got %v", order.OrderNumber, order.InvoiceNumber)
	}
	if order.DeliveryAddress == nil || *order.DeliveryAddress != testDeliveryAddress {
		t.Fatalf("expected delivery address stored, got %v", order.DeliveryAddress)
	}
	if *order.SubtotalAmount != 1000 || *order.DeliveryFee != 500 || *order.TaxAmount != 75 || *order.DiscountAmount != 0 {
		t.Fatalf("expected amount breakdown stored, got %v/%v/%v/%v", *order.SubtotalAmount, *order.DeliveryFee, *order.TaxAmount, *order.DiscountAmount)
	}

	select {
	case sent := <-invoices.sent:
		if sent.InvoiceNumber == nil {
			t.Fatal("expected the invoiced order to be emailed")
		}
	case <-time.After(time.Second):
		t.Fatal("expected the invoice to be sent")
	}
}

func TestCreateOrderOnlinePaymentHasNoInvoice(t *testing.T) {
	orderSvc := &fakeOrderService{subtotal: 1000}
	handler := NewOrderHandler(orderSvc, &fakeAuthService{}, acceptingSettings(), &fakePromotionService{}, &fakePaymentService{}, &fakeInvoiceService{})

	body := orderBody(t, func(r *models.OrderWithAuth) { r.PaymentMethod = models.PaymentMethodTransfer })
	rec := runCreateOrder(t, handler, body)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	if orderSvc.captured.InvoiceNumber != nil {
		t.Fatalf("online payments are not invoiced, got %v", *orderSvc.captured.InvoiceNumber)
	}
}

func TestCreateOrderRespectsPayOnDeliverySwitch(t *testing.T) {
	settings := acceptingSettings()
	settings.payOnDeliveryDisabled = true

	t.Run("pay on delivery rejected", func(t *testing.T) {
		orderSvc := &fakeOrderService{subtotal: 1000}
		handler := NewOrderHandler(orderSvc, &fakeAuthService{}, settings, &fakePromotionService{}, &fakePaymentService{}, &fakeInvoiceService{})
		rec := runCreateOrder(t, handler, newOrderRequestBody(t, models.OrderTypeDelivery))

		if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "Pay on delivery is not available") {
			t.Fatalf("expected pay-on-delivery rejection, got %d: %s", rec.Code, rec.Body.String())
		}
		if orderSvc.captured != nil {
			t.Fatal("order must not be created")
		}
	})

	t.Run("online payment still accepted", func(t *testing.T) {
		orderSvc := &fakeOrderService{subtotal: 1000}
		handler := NewOrderHandler(orderSvc, &fakeAuthService{}, settings, &fakePromotionService{}, &fakePaymentService{}, &fakeInvoiceService{})
		body := orderBody(t, func(r *models.OrderWithAuth) { r.PaymentMethod = models.PaymentMethodTransfer })
		rec := runCreateOrder(t, handler, body)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
		}
	})
}

func TestCreateOrderValidatesCustomerDetails(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*models.OrderWithAuth)
		wantErr string
	}{
		{"delivery needs an address", func(r *models.OrderWithAuth) { r.DeliveryAddress = nil }, "delivery address is required"},
		{"blank address", func(r *models.OrderWithAuth) { blank := "   "; r.DeliveryAddress = &blank }, "delivery address is required"},
		{"email is required", func(r *models.OrderWithAuth) { r.CustomerEmail = nil }, "CustomerEmail"},
		{"invalid phone", func(r *models.OrderWithAuth) { r.CustomerPhone = "12-345" }, "valid phone number"},
		{"phone with no digits", func(r *models.OrderWithAuth) { r.CustomerPhone = "call me" }, "CustomerPhone"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			orderSvc := &fakeOrderService{subtotal: 1000}
			handler := NewOrderHandler(orderSvc, &fakeAuthService{}, acceptingSettings(), &fakePromotionService{}, &fakePaymentService{}, &fakeInvoiceService{})
			rec := runCreateOrder(t, handler, orderBody(t, tc.mutate))

			if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), tc.wantErr) {
				t.Fatalf("expected 400 containing %q, got %d: %s", tc.wantErr, rec.Code, rec.Body.String())
			}
			if orderSvc.captured != nil {
				t.Fatal("order must not be created")
			}
		})
	}

	t.Run("pickup drops a stray address", func(t *testing.T) {
		orderSvc := &fakeOrderService{subtotal: 1000}
		handler := NewOrderHandler(orderSvc, &fakeAuthService{}, acceptingSettings(), &fakePromotionService{}, &fakePaymentService{}, &fakeInvoiceService{})
		body := orderBody(t, func(r *models.OrderWithAuth) { r.OrderType = models.OrderTypePickup })
		rec := runCreateOrder(t, handler, body)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
		}
		if orderSvc.captured.DeliveryAddress != nil {
			t.Fatalf("pickup orders store no address, got %q", *orderSvc.captured.DeliveryAddress)
		}
	})
}

func runLookup(t *testing.T, handler *OrderHandler, body string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/orders/lookup", handler.LookupOrder)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/orders/lookup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)
	return rec
}

func TestLookupOrder(t *testing.T) {
	invoiceNumber := "INV-FG-2026-10-1"
	order := &models.Order{
		ID: 1, OrderNumber: "FG-2026-10-1", CustomerPhone: "+2348012345678",
		TotalAmount: 1575, InvoiceNumber: &invoiceNumber,
	}

	t.Run("matching phone in another format returns order and invoice", func(t *testing.T) {
		handler := NewOrderHandler(&fakeOrderService{byNumber: order}, &fakeAuthService{}, acceptingSettings(), &fakePromotionService{}, &fakePaymentService{}, &fakeInvoiceService{})
		rec := runLookup(t, handler, `{"orderNumber":"FG-2026-10-1","phone":"0801 234 5678"}`)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
		}
		var payload struct {
			Order   models.Order    `json:"order"`
			Invoice *models.Invoice `json:"invoice"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
			t.Fatal(err)
		}
		if payload.Order.OrderNumber != order.OrderNumber || payload.Invoice == nil || payload.Invoice.InvoiceNumber != invoiceNumber {
			t.Fatalf("unexpected payload: %s", rec.Body.String())
		}
	})

	for name, tc := range map[string]struct {
		found *models.Order
		body  string
	}{
		"wrong phone":   {order, `{"orderNumber":"FG-2026-10-1","phone":"08099999999"}`},
		"unknown order": {nil, `{"orderNumber":"FG-404","phone":"08012345678"}`},
	} {
		t.Run(name+" is a generic not found", func(t *testing.T) {
			handler := NewOrderHandler(&fakeOrderService{byNumber: tc.found}, &fakeAuthService{}, acceptingSettings(), &fakePromotionService{}, &fakePaymentService{}, &fakeInvoiceService{})
			rec := runLookup(t, handler, tc.body)

			if rec.Code != http.StatusNotFound || !strings.Contains(rec.Body.String(), "Order not found") {
				t.Fatalf("expected generic 404, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestMarkPaymentCollected(t *testing.T) {
	run := func(orderSvc *fakeOrderService, body string) *httptest.ResponseRecorder {
		gin.SetMode(gin.TestMode)
		handler := NewOrderHandler(orderSvc, &fakeAuthService{}, acceptingSettings(), &fakePromotionService{}, &fakePaymentService{}, &fakeInvoiceService{})
		router := gin.New()
		router.PATCH("/admin/orders/:id/payment", handler.MarkPaymentCollected)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/admin/orders/5/payment", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rec, req)
		return rec
	}

	t.Run("records how the payment was collected", func(t *testing.T) {
		orderSvc := &fakeOrderService{collected: &models.Order{ID: 5, PaymentStatus: models.PaymentStatusPaid}}
		rec := run(orderSvc, `{"collectedVia":"pos"}`)

		if rec.Code != http.StatusOK || orderSvc.collectedVia != "pos" {
			t.Fatalf("expected 200 with pos recorded, got %d (%q): %s", rec.Code, orderSvc.collectedVia, rec.Body.String())
		}
	})

	cases := map[string]struct {
		svc  *fakeOrderService
		body string
		want int
	}{
		"unknown method":      {&fakeOrderService{}, `{"collectedVia":"cheque"}`, http.StatusBadRequest},
		"order not found":     {&fakeOrderService{collectErr: services.ErrOrderNotFound}, `{"collectedVia":"cash"}`, http.StatusNotFound},
		"already paid/online": {&fakeOrderService{collectErr: services.ErrPaymentNotCollectable}, `{"collectedVia":"cash"}`, http.StatusConflict},
		"unexpected failure":  {&fakeOrderService{collectErr: errors.New("db down")}, `{"collectedVia":"cash"}`, http.StatusInternalServerError},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if rec := run(tc.svc, tc.body); rec.Code != tc.want {
				t.Fatalf("expected %d, got %d: %s", tc.want, rec.Code, rec.Body.String())
			}
		})
	}
}
