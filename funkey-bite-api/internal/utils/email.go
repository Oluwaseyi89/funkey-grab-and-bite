package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type EmailService interface {
	SendOrderConfirmation(toEmail, toName, orderNumber string, totalAmount float64) error
	SendOrderStatusUpdate(toEmail, toName, orderNumber, status string) error
	SendCateringConfirmation(toEmail, toName, requestID, eventName string) error
	SendPasswordReset(toEmail, resetToken string) error
}

type AWSEmailService struct {
	sesClient *ses.Client
	sender    string
}

func NewEmailService() EmailService {
	if os.Getenv("ENVIRONMENT") == "development" {
		return &MockEmailService{}
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Printf("Warning: Failed to load AWS config: %v. Using mock email service.", err)
		return &MockEmailService{}
	}

	return &AWSEmailService{
		sesClient: ses.NewFromConfig(cfg),
		sender:    os.Getenv("SES_SENDER_EMAIL"),
	}
}

func (s *AWSEmailService) SendOrderConfirmation(toEmail, toName, orderNumber string, totalAmount float64) error {
	subject := fmt.Sprintf("Order Confirmation - %s", orderNumber)
	htmlBody := fmt.Sprintf(`
		<html>
		<body>
			<h2>Order Confirmation</h2>
			<p>Hello %s,</p>
			<p>Thank you for your order at Funkey Grab & Bite!</p>
			<p><strong>Order Number:</strong> %s</p>
			<p><strong>Total Amount:</strong> $%.2f</p>
			<p>We'll notify you when your order is ready for pickup/delivery.</p>
			<p>Thank you for choosing Funkey Grab & Bite!</p>
		</body>
		</html>
	`, toName, orderNumber, totalAmount)

	textBody := fmt.Sprintf(
		"Hello %s,\n\nThank you for your order at Funkey Grab & Bite!\n\nOrder Number: %s\nTotal Amount: $%.2f\n\nWe'll notify you when your order is ready.\n\nThank you!",
		toName, orderNumber, totalAmount,
	)

	return s.sendEmail(toEmail, subject, htmlBody, textBody)
}

func (s *AWSEmailService) SendOrderStatusUpdate(toEmail, toName, orderNumber, status string) error {
	statusMap := map[string]string{
		"confirmed": "Confirmed",
		"preparing": "Preparing",
		"ready":     "Ready for Pickup",
		"completed": "Completed",
		"cancelled": "Cancelled",
	}

	statusText, ok := statusMap[status]
	if !ok {
		statusText = status
	}

	subject := fmt.Sprintf("Order Update - %s", orderNumber)
	htmlBody := fmt.Sprintf(`
		<html>
		<body>
			<h2>Order Status Update</h2>
			<p>Hello %s,</p>
			<p>Your order status has been updated:</p>
			<p><strong>Order Number:</strong> %s</p>
			<p><strong>New Status:</strong> %s</p>
			<p>Thank you for choosing Funkey Grab & Bite!</p>
		</body>
		</html>
	`, toName, orderNumber, statusText)

	textBody := fmt.Sprintf(
		"Hello %s,\n\nYour order status has been updated.\n\nOrder Number: %s\nNew Status: %s\n\nThank you!",
		toName, orderNumber, statusText,
	)

	return s.sendEmail(toEmail, subject, htmlBody, textBody)
}

func (s *AWSEmailService) SendCateringConfirmation(toEmail, toName, requestID, eventName string) error {
	subject := "Catering Request Received"
	htmlBody := fmt.Sprintf(`
		<html>
		<body>
			<h2>Catering Request Received</h2>
			<p>Hello %s,</p>
			<p>Thank you for your catering request with Funkey Grab & Bite!</p>
			<p><strong>Request ID:</strong> %s</p>
			<p><strong>Event:</strong> %s</p>
			<p>Our team will review your request and contact you within 24 hours to discuss the details.</p>
			<p>Thank you for choosing Funkey Grab & Bite for your event!</p>
		</body>
		</html>
	`, toName, requestID, eventName)

	textBody := fmt.Sprintf(
		"Hello %s,\n\nThank you for your catering request!\n\nRequest ID: %s\nEvent: %s\n\nOur team will contact you within 24 hours.\n\nThank you!",
		toName, requestID, eventName,
	)

	return s.sendEmail(toEmail, subject, htmlBody, textBody)
}

func (s *AWSEmailService) SendPasswordReset(toEmail, resetToken string) error {
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("APP_URL"), resetToken)

	subject := "Password Reset Request"
	htmlBody := fmt.Sprintf(`
		<html>
		<body>
			<h2>Password Reset</h2>
			<p>You have requested to reset your password for Funkey Grab & Bite.</p>
			<p><a href="%s">Click here to reset your password</a></p>
			<p>This link will expire in 1 hour.</p>
			<p>If you didn't request this, please ignore this email.</p>
		</body>
		</html>
	`, resetLink)

	textBody := fmt.Sprintf(
		"Password Reset Request\n\nClick this link to reset your password:\n%s\n\nThis link expires in 1 hour.\n\nIf you didn't request this, please ignore this email.",
		resetLink,
	)

	return s.sendEmail(toEmail, subject, htmlBody, textBody)
}

func (s *AWSEmailService) sendEmail(toEmail, subject, htmlBody, textBody string) error {
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{toEmail},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Data: aws.String(htmlBody),
				},
				Text: &types.Content{
					Data: aws.String(textBody),
				},
			},
			Subject: &types.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String(s.sender),
	}

	ctx := context.Background()
	_, err := s.sesClient.SendEmail(ctx, input)
	return err
}

type MockEmailService struct{}

func (m *MockEmailService) SendOrderConfirmation(toEmail, toName, orderNumber string, totalAmount float64) error {
	log.Printf("[MOCK EMAIL] Order confirmation sent to %s: Order %s - $%.2f", toEmail, orderNumber, totalAmount)
	return nil
}

func (m *MockEmailService) SendOrderStatusUpdate(toEmail, toName, orderNumber, status string) error {
	log.Printf("[MOCK EMAIL] Order status update sent to %s: Order %s - Status: %s", toEmail, orderNumber, status)
	return nil
}

func (m *MockEmailService) SendCateringConfirmation(toEmail, toName, requestID, eventName string) error {
	log.Printf("[MOCK EMAIL] Catering confirmation sent to %s: Request %s - Event: %s", toEmail, requestID, eventName)
	return nil
}

func (m *MockEmailService) SendPasswordReset(toEmail, resetToken string) error {
	log.Printf("[MOCK EMAIL] Password reset sent to %s: Token: %s", toEmail, resetToken)
	return nil
}
