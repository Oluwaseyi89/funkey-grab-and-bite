package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestTermiiService(t *testing.T, handler http.HandlerFunc) *TermiiSMSService {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return &TermiiSMSService{
		apiKey:   "test-api-key",
		senderID: "FunkeyBite",
		baseURL:  server.URL,
		channel:  "dnd",
		client:   &http.Client{Timeout: 5 * time.Second},
	}
}

func TestTermiiSendSMSRequestShape(t *testing.T) {
	var captured termiiSendRequest
	var gotPath, gotContentType string

	svc := newTestTermiiService(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotContentType = r.Header.Get("Content-Type")
		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(termiiSendResponse{Code: "ok", Message: "Successfully Sent"})
	})

	if err := svc.SendVerificationCode("+234 801 234 5678", "123456"); err != nil {
		t.Fatalf("SendVerificationCode() error = %v", err)
	}

	if gotPath != "/api/sms/send" {
		t.Errorf("path = %q, want /api/sms/send", gotPath)
	}
	if gotContentType != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", gotContentType)
	}
	if captured.To != "2348012345678" {
		t.Errorf("to = %q, want digits-only with no '+' or spaces (2348012345678)", captured.To)
	}
	if captured.From != "FunkeyBite" {
		t.Errorf("from = %q, want FunkeyBite", captured.From)
	}
	if captured.Channel != "dnd" {
		t.Errorf("channel = %q, want dnd (transactional route)", captured.Channel)
	}
	if captured.Type != "plain" {
		t.Errorf("type = %q, want plain", captured.Type)
	}
	if captured.APIKey != "test-api-key" {
		t.Errorf("api_key = %q, want test-api-key", captured.APIKey)
	}
	if captured.SMS == "" {
		t.Error("sms body was empty")
	}
}

func TestTermiiSendSMSErrorOnNonOKCode(t *testing.T) {
	svc := newTestTermiiService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(termiiSendResponse{Code: "error", Message: "Insufficient balance"})
	})

	err := svc.SendOrderConfirmation("2348012345678", "ORD-1", 12.5)
	if err == nil {
		t.Fatal("expected an error when Termii returns a non-ok code, got nil")
	}
}

func TestTermiiSendSMSErrorOnHTTPFailureStatus(t *testing.T) {
	svc := newTestTermiiService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(termiiSendResponse{Message: "Invalid API key"})
	})

	err := svc.SendCateringConfirmation("2348012345678", "REQ-1")
	if err == nil {
		t.Fatal("expected an error on HTTP 401, got nil")
	}
}
