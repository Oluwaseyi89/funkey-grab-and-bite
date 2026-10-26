package services

import (
	"strings"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type invoiceSettingsStub struct {
	SettingsService
	settings *models.BusinessSettings
}

func (s invoiceSettingsStub) GetSettings() (*models.BusinessSettings, error) { return s.settings, nil }

type capturedEmail struct{ to, subject, html, text string }

type invoiceEmailStub struct {
	sent []capturedEmail
}

func (e *invoiceEmailStub) SendOrderConfirmation(string, string, string, float64) error { return nil }
func (e *invoiceEmailStub) SendOrderStatusUpdate(string, string, string, string) error  { return nil }
func (e *invoiceEmailStub) SendCateringConfirmation(string, string, string, string) error {
	return nil
}
func (e *invoiceEmailStub) SendPasswordReset(string, string) error { return nil }
func (e *invoiceEmailStub) SendInvoice(to, subject, html, text string) error {
	e.sent = append(e.sent, capturedEmail{to, subject, html, text})
	return nil
}

func invoicedOrder() *models.Order {
	invoiceNumber, email, address := "INV-FG-2026-10-9", "ada@example.com", "12 Admiralty Way, Lekki"
	subtotal, fee, tax, discount := 2000.0, 500.0, 150.0, 100.0
	return &models.Order{
		ID: 9, OrderNumber: "FG-2026-10-9", CustomerName: "Ada <script>alert(1)</script>",
		CustomerPhone: "08012345678", CustomerEmail: &email, DeliveryAddress: &address,
		OrderType: models.OrderTypeDelivery, TotalAmount: 2550, CreatedAt: time.Date(2026, 10, 26, 12, 0, 0, 0, time.UTC),
		PaymentMethod: models.PaymentMethodCash, PaymentStatus: models.PaymentStatusPending,
		InvoiceNumber: &invoiceNumber, SubtotalAmount: &subtotal, DeliveryFee: &fee, TaxAmount: &tax, DiscountAmount: &discount,
		Items: []models.OrderItem{{Name: "Shawarma", Quantity: 2, UnitPrice: 1000}},
	}
}

func newInvoiceTestService(email *invoiceEmailStub) InvoiceService {
	return NewInvoiceService(invoiceSettingsStub{settings: &models.BusinessSettings{
		BusinessName: "Funkey Grab & Bite", PhoneNumber: "+2348000000000", Email: "hello@example.com",
		Address: "1 Market Road", City: "Ikeja", State: "Lagos", Country: "Nigeria",
	}}, email)
}

func TestFormatNaira(t *testing.T) {
	for amount, want := range map[float64]string{
		0:         "₦0.00",
		999.5:     "₦999.50",
		1234567.8: "₦1,234,567.80",
		-2500:     "-₦2,500.00",
	} {
		if got := FormatNaira(amount); got != want {
			t.Errorf("FormatNaira(%v) = %q, want %q", amount, got, want)
		}
	}
}

func TestBuildInvoice(t *testing.T) {
	svc := newInvoiceTestService(&invoiceEmailStub{})

	invoice, err := svc.BuildInvoice(invoicedOrder())
	if err != nil {
		t.Fatal(err)
	}
	if invoice.InvoiceNumber != "INV-FG-2026-10-9" || invoice.Total != 2550 || invoice.Subtotal != 2000 || invoice.Discount != 100 {
		t.Fatalf("unexpected invoice: %+v", invoice)
	}
	if invoice.Business.Address != "1 Market Road, Ikeja, Lagos, Nigeria" {
		t.Fatalf("unexpected business address %q", invoice.Business.Address)
	}
	if invoice.Customer.Address != "12 Admiralty Way, Lekki" || len(invoice.Items) != 1 || invoice.Items[0].Total != 2000 {
		t.Fatalf("unexpected customer or lines: %+v", invoice)
	}

	online := invoicedOrder()
	online.InvoiceNumber = nil
	if invoice, err := svc.BuildInvoice(online); err != nil || invoice != nil {
		t.Fatalf("orders without an invoice number have no invoice, got %+v, %v", invoice, err)
	}
}

func TestSendInvoiceEscapesCustomerInput(t *testing.T) {
	email := &invoiceEmailStub{}
	if err := newInvoiceTestService(email).SendInvoice(invoicedOrder()); err != nil {
		t.Fatal(err)
	}
	if len(email.sent) != 1 {
		t.Fatalf("expected one invoice email, got %d", len(email.sent))
	}

	sent := email.sent[0]
	if sent.to != "ada@example.com" || !strings.Contains(sent.subject, "INV-FG-2026-10-9") {
		t.Fatalf("unexpected recipient or subject: %q %q", sent.to, sent.subject)
	}
	if strings.Contains(sent.html, "<script>") {
		t.Fatal("customer input must be escaped in the HTML invoice")
	}
	for _, want := range []string{"₦2,550.00", "Please pay ₦2,550.00 when your order is delivered", "12 Admiralty Way"} {
		if !strings.Contains(sent.html, want) {
			t.Errorf("HTML invoice missing %q", want)
		}
	}
	if !strings.Contains(sent.text, "Total due: ₦2,550.00") {
		t.Errorf("text invoice missing total: %s", sent.text)
	}
}

func TestSendInvoiceSkipsOrdersWithoutEmail(t *testing.T) {
	email := &invoiceEmailStub{}
	order := invoicedOrder()
	order.CustomerEmail = nil
	if err := newInvoiceTestService(email).SendInvoice(order); err != nil || len(email.sent) != 0 {
		t.Fatalf("expected no email, got %d, %v", len(email.sent), err)
	}
}
