package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type SMSService interface {
	SendOrderConfirmation(phoneNumber, orderNumber string, totalAmount float64) error
	SendOrderStatusUpdate(phoneNumber, orderNumber, status string) error
	SendCateringConfirmation(phoneNumber, requestID string) error
	SendVerificationCode(phoneNumber, code string) error
}

// TermiiSMSService sends SMS via Termii (https://developers.termii.com/messaging-api),
// which routes to Nigerian carriers more cheaply and reliably than international
// gateways like Twilio.
type TermiiSMSService struct {
	apiKey   string
	senderID string
	baseURL  string
	channel  string
	client   *http.Client
}

func NewSMSService() SMSService {
	if os.Getenv("ENVIRONMENT") == "development" || os.Getenv("TERMII_API_KEY") == "" {
		return &MockSMSService{}
	}

	return &TermiiSMSService{
		apiKey:   os.Getenv("TERMII_API_KEY"),
		senderID: os.Getenv("TERMII_SENDER_ID"),
		baseURL:  envOrDefault("TERMII_BASE_URL", "https://api.ng.termii.com"),
		// "dnd" is Termii's transactional route. Every message this service sends
		// (order/catering confirmations, status updates, OTP codes) is transactional,
		// never promotional — the "generic" route is for promotional campaigns and
		// explicitly should not be used for OTP/transactional messages, since it gets
		// filtered for numbers on Nigeria's Do-Not-Disturb list.
		channel: envOrDefault("TERMII_CHANNEL", "dnd"),
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func (s *TermiiSMSService) SendOrderConfirmation(phoneNumber, orderNumber string, totalAmount float64) error {
	message := fmt.Sprintf(
		"Thank you for your order at Funkey Grab & Bite! Order #%s for $%.2f. We'll notify you when it's ready.",
		orderNumber, totalAmount,
	)
	return s.sendSMS(phoneNumber, message)
}

func (s *TermiiSMSService) SendOrderStatusUpdate(phoneNumber, orderNumber, status string) error {
	statusMessages := map[string]string{
		"confirmed": "Your order #%s has been confirmed and is being prepared.",
		"preparing": "Your order #%s is now being prepared.",
		"ready":     "Your order #%s is ready for pickup!",
		"completed": "Your order #%s has been completed. Thank you!",
		"cancelled": "Your order #%s has been cancelled.",
	}

	messageTemplate, ok := statusMessages[status]
	if !ok {
		messageTemplate = "Your order #%s status has been updated to: %s"
	}

	message := fmt.Sprintf(messageTemplate, orderNumber, status)
	return s.sendSMS(phoneNumber, message)
}

func (s *TermiiSMSService) SendCateringConfirmation(phoneNumber, requestID string) error {
	message := fmt.Sprintf(
		"Thank you for your catering request with Funkey Grab & Bite! Request #%s. We'll contact you within 24 hours.",
		requestID,
	)
	return s.sendSMS(phoneNumber, message)
}

func (s *TermiiSMSService) SendVerificationCode(phoneNumber, code string) error {
	message := fmt.Sprintf(
		"Your Funkey Grab & Bite verification code is: %s. Valid for 10 minutes.",
		code,
	)
	return s.sendSMS(phoneNumber, message)
}

type termiiSendRequest struct {
	To      string `json:"to"`
	From    string `json:"from"`
	SMS     string `json:"sms"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	APIKey  string `json:"api_key"`
}

type termiiSendResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	MessageID string `json:"message_id"`
}

func (s *TermiiSMSService) sendSMS(to, body string) error {
	// Termii expects digits only (country code + number), unlike Twilio's "+"-prefixed E.164.
	to = strings.TrimPrefix(strings.ReplaceAll(to, " ", ""), "+")

	payload := termiiSendRequest{
		To:      to,
		From:    s.senderID,
		SMS:     body,
		Type:    "plain",
		Channel: s.channel,
		APIKey:  s.apiKey,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode Termii request: %w", err)
	}

	urlStr := fmt.Sprintf("%s/api/sms/send", strings.TrimRight(s.baseURL, "/"))
	req, err := http.NewRequest("POST", urlStr, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS via Termii: %w", err)
	}
	defer resp.Body.Close()

	var result termiiSendResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode Termii response (status %d): %w", resp.StatusCode, err)
	}

	if resp.StatusCode >= 300 || (result.Code != "" && result.Code != "ok") {
		return fmt.Errorf("Termii SMS failed (status %d): %s", resp.StatusCode, result.Message)
	}

	return nil
}

type MockSMSService struct{}

func (m *MockSMSService) SendOrderConfirmation(phoneNumber, orderNumber string, totalAmount float64) error {
	log.Printf("[MOCK SMS] Order confirmation sent to %s: Order %s - $%.2f", phoneNumber, orderNumber, totalAmount)
	return nil
}

func (m *MockSMSService) SendOrderStatusUpdate(phoneNumber, orderNumber, status string) error {
	log.Printf("[MOCK SMS] Order status update sent to %s: Order %s - Status: %s", phoneNumber, orderNumber, status)
	return nil
}

func (m *MockSMSService) SendCateringConfirmation(phoneNumber, requestID string) error {
	log.Printf("[MOCK SMS] Catering confirmation sent to %s: Request %s", phoneNumber, requestID)
	return nil
}

func (m *MockSMSService) SendVerificationCode(phoneNumber, code string) error {
	log.Printf("[MOCK SMS] Verification code sent to %s: Code: %s", phoneNumber, code)
	return nil
}
