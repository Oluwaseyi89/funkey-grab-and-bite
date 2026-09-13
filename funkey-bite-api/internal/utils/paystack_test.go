package utils

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestPaystackGateway(t *testing.T, handler http.HandlerFunc) *PaystackGateway {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return &PaystackGateway{
		secretKey: "sk_test_secret",
		baseURL:   server.URL,
		client:    &http.Client{Timeout: 5 * time.Second},
	}
}

func TestPaystackInitiateTransferRequestShape(t *testing.T) {
	var captured paystackChargeRequest
	var gotPath, gotAuth string

	gateway := newTestPaystackGateway(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		resp := paystackChargeResponse{Status: true, Message: "Charge attempted"}
		resp.Data.Status = "pending_bank_transfer"
		resp.Data.Reference = captured.Reference
		resp.Data.AccountNumber = "1231084927"
		resp.Data.AccountName = "PAYSTACK CHECKOUT"
		resp.Data.AccountExpiresAt = "2026-09-12T13:10:55.000Z"
		resp.Data.Bank.Name = "Test Bank"

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})

	result, err := gateway.InitiateTransfer("customer@example.com", 5000, "PSK-FG-1-abc123")
	if err != nil {
		t.Fatalf("InitiateTransfer() error = %v", err)
	}

	if gotPath != "/charge" {
		t.Errorf("path = %q, want /charge", gotPath)
	}
	if gotAuth != "Bearer sk_test_secret" {
		t.Errorf("Authorization = %q, want Bearer sk_test_secret", gotAuth)
	}
	if captured.Email != "customer@example.com" {
		t.Errorf("email = %q, want customer@example.com", captured.Email)
	}
	if captured.Amount != 500000 {
		t.Errorf("amount = %d, want 500000 (5000 NGN in kobo)", captured.Amount)
	}
	if captured.Currency != "NGN" {
		t.Errorf("currency = %q, want NGN", captured.Currency)
	}
	if captured.BankTransfer.AccountExpiresAt == "" {
		t.Error("expected an account_expires_at to be sent")
	}

	if result.AccountNumber != "1231084927" {
		t.Errorf("AccountNumber = %q, want 1231084927", result.AccountNumber)
	}
	if result.BankName != "Test Bank" {
		t.Errorf("BankName = %q, want Test Bank", result.BankName)
	}
	if result.ExpiresAt == nil {
		t.Error("expected ExpiresAt to be parsed from account_expires_at")
	}
}

func TestPaystackInitiateTransferErrorOnFailedStatus(t *testing.T) {
	gateway := newTestPaystackGateway(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(paystackChargeResponse{Status: false, Message: "Insufficient permissions"})
	})

	_, err := gateway.InitiateTransfer("customer@example.com", 5000, "PSK-FG-1-abc123")
	if err == nil {
		t.Fatal("expected an error when Paystack returns status: false")
	}
}

func TestPaystackVerifyWebhookSignature(t *testing.T) {
	gateway := &PaystackGateway{secretKey: "sk_test_secret"}
	body := []byte(`{"event":"charge.success","data":{"reference":"PSK-FG-1-abc123"}}`)

	mac := hmac.New(sha512.New, []byte("sk_test_secret"))
	mac.Write(body)
	validSignature := hex.EncodeToString(mac.Sum(nil))

	if !gateway.VerifyWebhookSignature(body, validSignature) {
		t.Error("expected a correctly computed signature to verify")
	}
	if gateway.VerifyWebhookSignature(body, "0000000000000000") {
		t.Error("expected an incorrect signature to be rejected")
	}
	if gateway.VerifyWebhookSignature(body, "") {
		t.Error("expected an empty signature to be rejected")
	}
	// A signature computed over different bytes (e.g. a re-serialized copy)
	// must not verify — this is the exact mistake the handler must avoid.
	tamperedBody := []byte(`{"event":"charge.success","data":{"reference":"PSK-FG-1-different"}}`)
	if gateway.VerifyWebhookSignature(tamperedBody, validSignature) {
		t.Error("expected a signature computed over different bytes to fail verification")
	}
}
