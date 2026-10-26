package services

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"text/tabwriter"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

// InvoicePrefix starts every invoice number; the rest is the order number,
// which is already unique.
const InvoicePrefix = "INV-"

// InvoiceNumberFor returns the invoice number raised for an order.
func InvoiceNumberFor(orderNumber string) string {
	return InvoicePrefix + orderNumber
}

type InvoiceService interface {
	// BuildInvoice returns the invoice for an order, or nil when the order has
	// none (only pay-on-delivery orders are invoiced).
	BuildInvoice(order *models.Order) (*models.Invoice, error)
	// SendInvoice emails the invoice to the customer, if the order has one.
	SendInvoice(order *models.Order) error
}

type invoiceService struct {
	settings SettingsService
	email    utils.EmailService
}

func NewInvoiceService(settings SettingsService, email utils.EmailService) InvoiceService {
	return &invoiceService{settings: settings, email: email}
}

func (s *invoiceService) BuildInvoice(order *models.Order) (*models.Invoice, error) {
	if order == nil || order.InvoiceNumber == nil {
		return nil, nil
	}

	settings, err := s.settings.GetSettings()
	if err != nil {
		return nil, fmt.Errorf("loading business details for invoice: %w", err)
	}

	lines := make([]models.InvoiceLine, 0, len(order.Items))
	itemsTotal := 0.0
	for _, item := range order.Items {
		lineTotal := item.UnitPrice * float64(item.Quantity)
		itemsTotal += lineTotal
		lines = append(lines, models.InvoiceLine{
			Name:      item.Name,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Total:     lineTotal,
		})
	}

	customer := models.InvoiceParty{Name: order.CustomerName, Phone: order.CustomerPhone}
	if order.CustomerEmail != nil {
		customer.Email = *order.CustomerEmail
	}
	if order.DeliveryAddress != nil {
		customer.Address = *order.DeliveryAddress
	}

	return &models.Invoice{
		InvoiceNumber: *order.InvoiceNumber,
		OrderNumber:   order.OrderNumber,
		IssuedAt:      order.CreatedAt,
		Business: models.InvoiceParty{
			Name:    settings.BusinessName,
			Phone:   settings.PhoneNumber,
			Email:   settings.Email,
			Address: joinNonEmpty(", ", settings.Address, settings.City, settings.State, settings.PostalCode, settings.Country),
		},
		Customer:            customer,
		OrderType:           order.OrderType,
		Items:               lines,
		Subtotal:            valueOr(order.SubtotalAmount, itemsTotal),
		DeliveryFee:         valueOr(order.DeliveryFee, 0),
		Tax:                 valueOr(order.TaxAmount, 0),
		Discount:            valueOr(order.DiscountAmount, 0),
		Total:               order.TotalAmount,
		PaymentMethod:       order.PaymentMethod,
		PaymentStatus:       order.PaymentStatus,
		PaymentCollectedVia: order.PaymentCollectedVia,
		PaidAt:              order.PaymentPaidAt,
	}, nil
}

func (s *invoiceService) SendInvoice(order *models.Order) error {
	if order == nil || order.CustomerEmail == nil || strings.TrimSpace(*order.CustomerEmail) == "" {
		return nil
	}

	invoice, err := s.BuildInvoice(order)
	if err != nil || invoice == nil {
		return err
	}

	htmlBody, err := RenderInvoiceHTML(invoice)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("Invoice %s from %s", invoice.InvoiceNumber, invoice.Business.Name)
	return s.email.SendInvoice(*order.CustomerEmail, subject, htmlBody, RenderInvoiceText(invoice))
}

// FormatNaira formats an amount as "₦12,345.00".
func FormatNaira(amount float64) string {
	negative := amount < 0
	if negative {
		amount = -amount
	}
	whole := fmt.Sprintf("%.2f", amount)
	intPart, decimals := whole[:len(whole)-3], whole[len(whole)-3:]

	var grouped strings.Builder
	for i, digit := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			grouped.WriteByte(',')
		}
		grouped.WriteRune(digit)
	}

	sign := ""
	if negative {
		sign = "-"
	}
	return sign + "₦" + grouped.String() + decimals
}

var invoiceTemplate = template.Must(template.New("invoice").Funcs(template.FuncMap{
	"naira": FormatNaira,
}).Parse(`<!doctype html>
<html>
<body style="font-family: Arial, sans-serif; color: #1f2937; max-width: 640px; margin: 0 auto;">
  <h2 style="margin-bottom: 4px;">Invoice {{.InvoiceNumber}}</h2>
  <p style="margin-top: 0; color: #6b7280;">Order {{.OrderNumber}} &middot; {{.IssuedAt.Format "2 Jan 2006, 15:04"}}</p>

  <table style="width: 100%; margin: 16px 0;">
    <tr>
      <td style="vertical-align: top; width: 50%;">
        <strong>From</strong><br>
        {{.Business.Name}}<br>
        {{if .Business.Address}}{{.Business.Address}}<br>{{end}}
        {{.Business.Phone}}{{if .Business.Email}}<br>{{.Business.Email}}{{end}}
      </td>
      <td style="vertical-align: top; width: 50%;">
        <strong>Bill to</strong><br>
        {{.Customer.Name}}<br>
        {{.Customer.Phone}}{{if .Customer.Email}}<br>{{.Customer.Email}}{{end}}
        {{if .Customer.Address}}<br>Deliver to: {{.Customer.Address}}{{end}}
      </td>
    </tr>
  </table>

  <table style="width: 100%; border-collapse: collapse;">
    <thead>
      <tr style="background: #f3f4f6;">
        <th style="text-align: left; padding: 8px;">Item</th>
        <th style="text-align: right; padding: 8px;">Qty</th>
        <th style="text-align: right; padding: 8px;">Price</th>
        <th style="text-align: right; padding: 8px;">Amount</th>
      </tr>
    </thead>
    <tbody>
      {{range .Items}}
      <tr>
        <td style="padding: 8px; border-bottom: 1px solid #e5e7eb;">{{.Name}}</td>
        <td style="padding: 8px; border-bottom: 1px solid #e5e7eb; text-align: right;">{{.Quantity}}</td>
        <td style="padding: 8px; border-bottom: 1px solid #e5e7eb; text-align: right;">{{naira .UnitPrice}}</td>
        <td style="padding: 8px; border-bottom: 1px solid #e5e7eb; text-align: right;">{{naira .Total}}</td>
      </tr>
      {{end}}
    </tbody>
  </table>

  <table style="width: 100%; margin-top: 12px;">
    <tr><td style="text-align: right;">Subtotal</td><td style="text-align: right; width: 140px;">{{naira .Subtotal}}</td></tr>
    {{if .DeliveryFee}}<tr><td style="text-align: right;">Delivery fee</td><td style="text-align: right;">{{naira .DeliveryFee}}</td></tr>{{end}}
    {{if .Tax}}<tr><td style="text-align: right;">Tax</td><td style="text-align: right;">{{naira .Tax}}</td></tr>{{end}}
    {{if .Discount}}<tr><td style="text-align: right;">Discount</td><td style="text-align: right;">-{{naira .Discount}}</td></tr>{{end}}
    <tr><td style="text-align: right; font-weight: bold;">Total due</td><td style="text-align: right; font-weight: bold;">{{naira .Total}}</td></tr>
  </table>

  <p style="margin-top: 24px; padding: 12px; background: #fef3c7; border-radius: 6px;">
    {{if eq .OrderType "delivery"}}Please pay {{naira .Total}} when your order is delivered.{{else}}Please pay {{naira .Total}} when you collect your order.{{end}}
    Cash, POS and bank transfer are accepted.
  </p>
</body>
</html>`))

// RenderInvoiceHTML renders the invoice email body. html/template escapes the
// customer-supplied fields.
func RenderInvoiceHTML(invoice *models.Invoice) (string, error) {
	var buf bytes.Buffer
	if err := invoiceTemplate.Execute(&buf, invoice); err != nil {
		return "", fmt.Errorf("rendering invoice: %w", err)
	}
	return buf.String(), nil
}

// RenderInvoiceText renders the plain-text alternative of the invoice email.
func RenderInvoiceText(invoice *models.Invoice) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Invoice %s\nOrder %s, %s\n\n", invoice.InvoiceNumber, invoice.OrderNumber, invoice.IssuedAt.Format("2 Jan 2006, 15:04"))
	fmt.Fprintf(&buf, "From: %s\n", joinNonEmpty(", ", invoice.Business.Name, invoice.Business.Address, invoice.Business.Phone))
	fmt.Fprintf(&buf, "Bill to: %s\n", joinNonEmpty(", ", invoice.Customer.Name, invoice.Customer.Phone, invoice.Customer.Email))
	if invoice.Customer.Address != "" {
		fmt.Fprintf(&buf, "Deliver to: %s\n", invoice.Customer.Address)
	}
	buf.WriteString("\n")

	table := tabwriter.NewWriter(&buf, 0, 4, 2, ' ', tabwriter.AlignRight)
	for _, line := range invoice.Items {
		fmt.Fprintf(table, "%s\t%dx\t%s\t\n", line.Name, line.Quantity, FormatNaira(line.Total))
	}
	table.Flush()

	fmt.Fprintf(&buf, "\nSubtotal: %s\n", FormatNaira(invoice.Subtotal))
	if invoice.DeliveryFee != 0 {
		fmt.Fprintf(&buf, "Delivery fee: %s\n", FormatNaira(invoice.DeliveryFee))
	}
	if invoice.Tax != 0 {
		fmt.Fprintf(&buf, "Tax: %s\n", FormatNaira(invoice.Tax))
	}
	if invoice.Discount != 0 {
		fmt.Fprintf(&buf, "Discount: -%s\n", FormatNaira(invoice.Discount))
	}
	fmt.Fprintf(&buf, "Total due: %s\n\nPlease pay when you receive your order. Cash, POS and bank transfer are accepted.\n", FormatNaira(invoice.Total))
	return buf.String()
}

func joinNonEmpty(sep string, parts ...string) string {
	kept := parts[:0:0]
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			kept = append(kept, strings.TrimSpace(part))
		}
	}
	return strings.Join(kept, sep)
}

func valueOr(value *float64, fallback float64) float64 {
	if value == nil {
		return fallback
	}
	return *value
}
