package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

// PaymentGateway is the payment-provider boundary, mirroring the EmailService/
// SMSService pattern elsewhere in this package: a real implementation (Paystack)
// behind an interface, with a mock used in development.
type PaymentGateway interface {
	// InitiateTransfer starts a Paystack "Pay with Transfer" charge and returns
	// the dedicated account details the customer should transfer to.
	InitiateTransfer(email string, amountNaira float64, reference string) (*TransferChargeResult, error)
	// VerifyWebhookSignature checks the x-paystack-signature header against the
	// raw (unparsed) request body. Must be called with the exact bytes Paystack
	// sent — re-serializing the parsed JSON will not reproduce the same signature.
	VerifyWebhookSignature(rawBody []byte, signature string) bool
}

type TransferChargeResult struct {
	Reference     string
	Status        string
	AccountNumber string
	AccountName   string
	BankName      string
	ExpiresAt     *time.Time
}

type PaystackGateway struct {
	secretKey string
	baseURL   string
	client    *http.Client
}

func NewPaymentGateway() PaymentGateway {
	if os.Getenv("ENVIRONMENT") == "development" || os.Getenv("PAYSTACK_SECRET_KEY") == "" {
		return &MockPaymentGateway{}
	}

	return &PaystackGateway{
		secretKey: os.Getenv("PAYSTACK_SECRET_KEY"),
		baseURL:   envOrDefault("PAYSTACK_BASE_URL", "https://api.paystack.co"),
		client:    &http.Client{Timeout: 15 * time.Second},
	}
}

type paystackBankTransferOpts struct {
	AccountExpiresAt string `json:"account_expires_at,omitempty"`
}

type paystackChargeRequest struct {
	Email        string                   `json:"email"`
	Amount       int64                    `json:"amount"`
	Currency     string                   `json:"currency"`
	Reference    string                   `json:"reference"`
	BankTransfer paystackBankTransferOpts `json:"bank_transfer"`
}

type paystackChargeResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Status           string `json:"status"`
		Reference        string `json:"reference"`
		AccountName      string `json:"account_name"`
		AccountNumber    string `json:"account_number"`
		AccountExpiresAt string `json:"account_expires_at"`
		Bank             struct {
			Name string `json:"name"`
		} `json:"bank"`
	} `json:"data"`
}

func (g *PaystackGateway) InitiateTransfer(email string, amountNaira float64, reference string) (*TransferChargeResult, error) {
	// Paystack accepts a caller-supplied window (defaults itself if we don't set
	// one, but its default range is 15min-8h); 15 minutes is a predictable,
	// deliberately short window for a POD-style checkout.
	expiresAt := time.Now().Add(15 * time.Minute).UTC().Format(time.RFC3339)

	payload := paystackChargeRequest{
		Email:     email,
		Amount:    int64(math.Round(amountNaira * 100)), // Naira -> kobo
		Currency:  "NGN",
		Reference: reference,
		BankTransfer: paystackBankTransferOpts{
			AccountExpiresAt: expiresAt,
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode Paystack charge request: %w", err)
	}

	req, err := http.NewRequest("POST", g.baseURL+"/charge", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+g.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reach Paystack: %w", err)
	}
	defer resp.Body.Close()

	var result paystackChargeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Paystack response (status %d): %w", resp.StatusCode, err)
	}

	if resp.StatusCode >= 300 || !result.Status {
		return nil, fmt.Errorf("Paystack charge failed (status %d): %s", resp.StatusCode, result.Message)
	}

	charge := &TransferChargeResult{
		Reference:     result.Data.Reference,
		Status:        result.Data.Status,
		AccountNumber: result.Data.AccountNumber,
		AccountName:   result.Data.AccountName,
		BankName:      result.Data.Bank.Name,
	}

	if result.Data.AccountExpiresAt != "" {
		if parsed, err := time.Parse(time.RFC3339Nano, result.Data.AccountExpiresAt); err == nil {
			charge.ExpiresAt = &parsed
		} else {
			log.Printf("[Paystack] could not parse account_expires_at %q: %v", result.Data.AccountExpiresAt, err)
		}
	}

	return charge, nil
}

func (g *PaystackGateway) VerifyWebhookSignature(rawBody []byte, signature string) bool {
	if signature == "" {
		return false
	}

	mac := hmac.New(sha512.New, []byte(g.secretKey))
	mac.Write(rawBody)
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}

// MockPaymentGateway simulates a Paystack transfer charge without any network
// call, used in development the same way MockEmailService/MockSMSService are.
type MockPaymentGateway struct{}

func (m *MockPaymentGateway) InitiateTransfer(email string, amountNaira float64, reference string) (*TransferChargeResult, error) {
	log.Printf("[MOCK PAYSTACK] Transfer charge initiated: ref=%s amount=NGN %.2f email=%s", reference, amountNaira, email)
	expiresAt := time.Now().Add(15 * time.Minute)
	return &TransferChargeResult{
		Reference:     reference,
		Status:        "pending_bank_transfer",
		AccountNumber: "0000000000",
		AccountName:   "Funkey Grab & Bite (Mock)",
		BankName:      "Mock Test Bank",
		ExpiresAt:     &expiresAt,
	}, nil
}

func (m *MockPaymentGateway) VerifyWebhookSignature(rawBody []byte, signature string) bool {
	log.Printf("[MOCK PAYSTACK] Skipping webhook signature verification in development")
	return true
}
