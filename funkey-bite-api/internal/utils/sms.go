package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

type TwilioSMSService struct {
	accountSID string
	authToken  string
	fromNumber string
	client     *http.Client
}

func NewSMSService() SMSService {
	if os.Getenv("ENVIRONMENT") == "development" || os.Getenv("TWILIO_ACCOUNT_SID") == "" {
		return &MockSMSService{}
	}

	return &TwilioSMSService{
		accountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
		authToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
		fromNumber: os.Getenv("TWILIO_PHONE_NUMBER"),
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *TwilioSMSService) SendOrderConfirmation(phoneNumber, orderNumber string, totalAmount float64) error {
	message := fmt.Sprintf(
		"Thank you for your order at Funkey Grab & Bite! Order #%s for $%.2f. We'll notify you when it's ready.",
		orderNumber, totalAmount,
	)
	return s.sendSMS(phoneNumber, message)
}

func (s *TwilioSMSService) SendOrderStatusUpdate(phoneNumber, orderNumber, status string) error {
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

func (s *TwilioSMSService) SendCateringConfirmation(phoneNumber, requestID string) error {
	message := fmt.Sprintf(
		"Thank you for your catering request with Funkey Grab & Bite! Request #%s. We'll contact you within 24 hours.",
		requestID,
	)
	return s.sendSMS(phoneNumber, message)
}

func (s *TwilioSMSService) SendVerificationCode(phoneNumber, code string) error {
	message := fmt.Sprintf(
		"Your Funkey Grab & Bite verification code is: %s. Valid for 10 minutes.",
		code,
	)
	return s.sendSMS(phoneNumber, message)
}

func (s *TwilioSMSService) sendSMS(to, body string) error {
	to = strings.ReplaceAll(to, " ", "")
	if !strings.HasPrefix(to, "+") {
		to = "+" + to
	}

	data := url.Values{}
	data.Set("To", to)
	data.Set("From", s.fromNumber)
	data.Set("Body", body)

	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s.accountSID)
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(s.accountSID, s.authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		return fmt.Errorf("SMS failed with status %d: %v", resp.StatusCode, result)
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
