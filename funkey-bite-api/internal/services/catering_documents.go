package services

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

// cateringDocument is the data for a catering quote, invoice or receipt email.
type cateringDocument struct {
	Kind      string // quote, invoice or receipt
	Title     string
	Number    string
	Business  models.InvoiceParty
	Request   *models.CateringRequest
	EventDate string
	Quote     *models.CateringQuote
	Lines     []models.CateringQuoteLine
	Total     float64
	Deposit   float64
	Paid      float64
	Balance   float64
	Payments  []models.CateringPayment
	Link      string
}

var cateringDocumentTemplate = template.Must(template.New("catering").Funcs(template.FuncMap{
	"naira":  FormatNaira,
	"method": func(m string) string { return strings.ReplaceAll(m, "_", " ") },
}).Parse(`<!doctype html>
<html>
<body style="font-family: Arial, sans-serif; color: #1f2937; max-width: 640px; margin: 0 auto;">
  <h2 style="margin-bottom: 4px;">{{.Title}}{{if .Number}} {{.Number}}{{end}}</h2>
  <p style="margin-top: 0; color: #6b7280;">Catering request {{.Request.Reference}} &middot; {{.EventDate}} &middot; {{.Request.GuestCount}} guests</p>

  <table style="width: 100%; margin: 16px 0;">
    <tr>
      <td style="vertical-align: top; width: 50%;">
        <strong>From</strong><br>
        {{.Business.Name}}<br>
        {{if .Business.Address}}{{.Business.Address}}<br>{{end}}
        {{.Business.Phone}}{{if .Business.Email}}<br>{{.Business.Email}}{{end}}
      </td>
      <td style="vertical-align: top; width: 50%;">
        <strong>For</strong><br>
        {{.Request.ContactName}}{{if .Request.Organization}} ({{.Request.Organization}}){{end}}<br>
        {{.Request.ContactPhone}}{{if .Request.ContactEmail}}<br>{{.Request.ContactEmail}}{{end}}
        {{if .Request.VenueAddress}}<br>Venue: {{.Request.VenueAddress}}{{end}}
      </td>
    </tr>
  </table>

  {{if .Lines}}
  <table style="width: 100%; border-collapse: collapse;">
    <thead>
      <tr style="background: #f3f4f6;">
        <th style="text-align: left; padding: 8px;">Description</th>
        <th style="text-align: right; padding: 8px;">Amount</th>
      </tr>
    </thead>
    <tbody>
      {{range .Lines}}
      <tr>
        <td style="padding: 8px; border-bottom: 1px solid #e5e7eb;">{{.Description}}</td>
        <td style="padding: 8px; border-bottom: 1px solid #e5e7eb; text-align: right;">{{naira .Amount}}</td>
      </tr>
      {{end}}
    </tbody>
  </table>
  {{end}}

  <table style="width: 100%; margin-top: 12px;">
    <tr><td style="text-align: right; font-weight: bold;">Total</td><td style="text-align: right; width: 140px; font-weight: bold;">{{naira .Total}}</td></tr>
    <tr><td style="text-align: right;">Deposit to confirm ({{.Quote.DepositPercent}}%)</td><td style="text-align: right;">{{naira .Deposit}}</td></tr>
    {{if ne .Kind "quote"}}
    <tr><td style="text-align: right;">Paid</td><td style="text-align: right;">{{naira .Paid}}</td></tr>
    <tr><td style="text-align: right; font-weight: bold;">Balance</td><td style="text-align: right; font-weight: bold;">{{naira .Balance}}</td></tr>
    {{end}}
  </table>

  {{if .Payments}}
  <h3 style="margin-top: 24px;">Payments</h3>
  <table style="width: 100%; border-collapse: collapse;">
    {{range .Payments}}
    <tr>
      <td style="padding: 6px 8px; border-bottom: 1px solid #e5e7eb;">{{.PaidAt.Format "2 Jan 2006"}}</td>
      <td style="padding: 6px 8px; border-bottom: 1px solid #e5e7eb; text-transform: capitalize;">{{method .Method}}</td>
      <td style="padding: 6px 8px; border-bottom: 1px solid #e5e7eb; text-align: right;">{{naira .Amount}}</td>
    </tr>
    {{end}}
  </table>
  {{end}}

  {{if eq .Kind "quote"}}
  {{if .Quote.Message}}<p style="margin-top: 24px;">{{.Quote.Message}}</p>{{end}}
  <p style="margin-top: 24px; padding: 12px; background: #fef3c7; border-radius: 6px;">
    This quote is valid until {{.Quote.ValidUntil}}. Accepting it books your event; a deposit of {{naira .Deposit}} confirms the booking.
  </p>
  <h3>Terms</h3>
  <p style="white-space: pre-line; font-size: 13px; color: #4b5563;">{{.Quote.Terms}}</p>
  {{else if eq .Kind "invoice"}}
  <p style="margin-top: 24px; padding: 12px; background: #fef3c7; border-radius: 6px;">
    {{if lt .Paid .Deposit}}Please pay at least {{naira .Deposit}} in total to confirm your booking.{{else}}Your booking is confirmed. Please pay the balance of {{naira .Balance}} before the event.{{end}}
    You can pay by bank transfer online, or by cash or transfer to us directly.
  </p>
  {{else}}
  <p style="margin-top: 24px; padding: 12px; background: #d1fae5; border-radius: 6px;">
    Thank you. Your booking is paid in full.
  </p>
  {{end}}

  {{if .Link}}
  <p style="margin-top: 24px;"><a href="{{.Link}}" style="background: #ea580c; color: #fff; padding: 10px 18px; border-radius: 6px; text-decoration: none;">{{if eq .Kind "quote"}}Review and accept the quote{{else}}View your booking{{end}}</a></p>
  {{end}}
</body>
</html>`))

func renderCateringDocumentText(doc *cateringDocument) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s\n", doc.Title, doc.Number)
	fmt.Fprintf(&b, "Catering request %s, %s, %d guests\n\n", doc.Request.Reference, doc.EventDate, doc.Request.GuestCount)
	for _, line := range doc.Lines {
		fmt.Fprintf(&b, "%s: %s\n", line.Description, FormatNaira(line.Amount))
	}
	fmt.Fprintf(&b, "Total: %s\n", FormatNaira(doc.Total))
	fmt.Fprintf(&b, "Deposit to confirm (%d%%): %s\n", doc.Quote.DepositPercent, FormatNaira(doc.Deposit))
	if doc.Kind != "quote" {
		fmt.Fprintf(&b, "Paid: %s\nBalance: %s\n", FormatNaira(doc.Paid), FormatNaira(doc.Balance))
	}
	switch doc.Kind {
	case "quote":
		if doc.Quote.Message != nil {
			fmt.Fprintf(&b, "\n%s\n", *doc.Quote.Message)
		}
		fmt.Fprintf(&b, "\nValid until %s.\n\nTerms:\n%s\n", doc.Quote.ValidUntil, doc.Quote.Terms)
	case "receipt":
		b.WriteString("\nThank you. Your booking is paid in full.\n")
	}
	if doc.Link != "" {
		fmt.Fprintf(&b, "\n%s\n", doc.Link)
	}
	fmt.Fprintf(&b, "\n%s, %s\n", doc.Business.Name, doc.Business.Phone)
	return b.String()
}

func (s *cateringContractService) newDocument(kind string, request *models.CateringRequest, quote *models.CateringQuote) *cateringDocument {
	doc := &cateringDocument{
		Kind: kind, Request: request, Quote: quote, Lines: quote.LineItems,
		Total: quote.Amount, Deposit: depositFor(quote), Balance: quote.Amount,
		EventDate: request.EventDate,
	}
	if settings, err := s.settings.GetSettings(); err == nil && settings != nil {
		doc.Business = models.InvoiceParty{
			Name:    settings.BusinessName,
			Phone:   settings.PhoneNumber,
			Email:   settings.Email,
			Address: joinNonEmpty(", ", settings.Address, settings.City, settings.State, settings.PostalCode, settings.Country),
		}
	}
	return doc
}

// customerLink is the page where the customer reviews their booking: the
// quote link when there is a token, else the status page.
func (s *cateringContractService) customerLink(token string) string {
	if s.webURL == "" {
		return ""
	}
	if token != "" {
		return s.webURL + "/catering/quote?token=" + token
	}
	return s.webURL + "/catering/request"
}

func (s *cateringContractService) deliver(request *models.CateringRequest, subject string, doc *cateringDocument) {
	if request.ContactEmail == nil || strings.TrimSpace(*request.ContactEmail) == "" {
		return
	}
	var html bytes.Buffer
	if err := cateringDocumentTemplate.Execute(&html, doc); err != nil {
		log.Printf("Failed to render catering %s for %s: %v", doc.Kind, request.Reference, err)
		return
	}
	if doc.Business.Name != "" {
		subject += " from " + doc.Business.Name
	}
	if err := s.email.SendInvoice(*request.ContactEmail, subject, html.String(), renderCateringDocumentText(doc)); err != nil {
		log.Printf("Failed to email catering %s for %s: %v", doc.Kind, request.Reference, err)
	}
}

func (s *cateringContractService) sendQuoteEmail(request *models.CateringRequest, quote *models.CateringQuote, token string) {
	doc := s.newDocument("quote", request, quote)
	doc.Title = "Quote"
	doc.Number = fmt.Sprintf("%s v%d", request.Reference, quote.Version)
	doc.Link = s.customerLink(token)
	s.deliver(request, "Your catering quote "+request.Reference, doc)
}

func (s *cateringContractService) sendInvoiceEmail(contract *models.CateringContract, invoice *models.CateringInvoice) {
	if contract.Quote == nil {
		return
	}
	doc := s.newDocument("invoice", contract.Request, contract.Quote)
	doc.Title = "Invoice"
	doc.Number = invoice.InvoiceNumber
	doc.Paid = invoice.AmountPaid
	doc.Balance = invoice.Balance
	doc.Payments = contract.Payments
	doc.Link = s.customerLink("")
	s.deliver(contract.Request, fmt.Sprintf("Invoice %s", invoice.InvoiceNumber), doc)
}

func (s *cateringContractService) sendReceiptEmail(contract *models.CateringContract, receipt *models.CateringReceipt) {
	if contract.Quote == nil {
		return
	}
	doc := s.newDocument("receipt", contract.Request, contract.Quote)
	doc.Title = "Receipt"
	doc.Number = receipt.ReceiptNumber
	doc.Paid = receipt.TotalPaid
	doc.Balance = 0
	doc.Payments = contract.Payments
	doc.Link = s.customerLink("")
	s.deliver(contract.Request, fmt.Sprintf("Receipt %s", receipt.ReceiptNumber), doc)
}
