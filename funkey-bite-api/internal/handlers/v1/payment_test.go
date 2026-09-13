package v1

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"funkey-grab-and-bite/funkey-bite-api/internal/services"
)

var errBoom = errors.New("boom")

func runWebhook(t *testing.T, h *PaymentHandler, body []byte, signature string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/payments/webhook", h.Webhook)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/payments/webhook", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if signature != "" {
		req.Header.Set("x-paystack-signature", signature)
	}
	router.ServeHTTP(recorder, req)

	return recorder
}

func TestPaymentWebhookSucceeds(t *testing.T) {
	svc := &fakePaymentService{}
	handler := NewPaymentHandler(svc)

	body := []byte(`{"event":"charge.success","data":{"reference":"PSK-FG-1-abc"}}`)
	rec := runWebhook(t, handler, body, "valid-signature")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	if !svc.webhookCalled {
		t.Fatal("expected HandleWebhook to be called")
	}
	if !bytes.Equal(svc.capturedBody, body) {
		t.Errorf("HandleWebhook received body %s, want %s", svc.capturedBody, body)
	}
	if svc.capturedSignature != "valid-signature" {
		t.Errorf("HandleWebhook received signature %q, want valid-signature", svc.capturedSignature)
	}
}

func TestPaymentWebhookRejectsInvalidSignature(t *testing.T) {
	svc := &fakePaymentService{handleWebhookErr: services.ErrInvalidWebhookSignature}
	handler := NewPaymentHandler(svc)

	rec := runWebhook(t, handler, []byte(`{}`), "bad-signature")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401; body = %s", rec.Code, rec.Body.String())
	}
}

func TestPaymentWebhookReturns500OnOtherFailures(t *testing.T) {
	svc := &fakePaymentService{handleWebhookErr: errBoom}
	handler := NewPaymentHandler(svc)

	rec := runWebhook(t, handler, []byte(`{}`), "some-signature")

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500; body = %s", rec.Code, rec.Body.String())
	}
}
