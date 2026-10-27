package models

import "time"

// DefaultCateringDepositPercent is the share of a quote paid to confirm a booking.
const DefaultCateringDepositPercent = 70

// Quote statuses. A quote is "sent" until the customer accepts it, the admin
// sends a newer version (superseded), or the request is closed (withdrawn).
const (
	QuoteStatusSent       = "sent"
	QuoteStatusAccepted   = "accepted"
	QuoteStatusSuperseded = "superseded"
	QuoteStatusWithdrawn  = "withdrawn"
)

// How a customer accepted a quote.
const (
	AcceptedViaEmailLink  = "email_link"
	AcceptedViaStatusPage = "status_page"
)

// Catering payment methods. Cash and bank transfers are recorded by the
// admin; Paystack payments are recorded by the webhook.
const (
	CateringPaymentCash             = "cash"
	CateringPaymentBankTransfer     = "bank_transfer"
	CateringPaymentPaystackTransfer = "paystack_transfer"
	CateringPaymentCard             = "card"
)

type CateringQuoteLine struct {
	Description string  `json:"description" validate:"required,max=200"`
	Amount      float64 `json:"amount" validate:"min=0"`
}

// CateringQuote is a priced offer for a request; the customer accepts one.
type CateringQuote struct {
	ID             int                 `json:"id"`
	RequestID      int                 `json:"requestId"`
	Version        int                 `json:"version"`
	Amount         float64             `json:"amount"`
	LineItems      []CateringQuoteLine `json:"lineItems"`
	DepositPercent int                 `json:"depositPercent"`
	DepositAmount  float64             `json:"depositAmount"`
	Terms          string              `json:"terms"`
	Message        *string             `json:"message,omitempty"`
	ValidUntil     string              `json:"validUntil"`
	Status         string              `json:"status"`
	IsExpired      bool                `json:"isExpired"`
	AcceptedAt     *time.Time          `json:"acceptedAt,omitempty"`
	AcceptedVia    *string             `json:"acceptedVia,omitempty"`
	CreatedAt      time.Time           `json:"createdAt"`
}

// CateringQuoteInput is the admin's quote; omitted options take defaults
// (70% deposit, the default terms, valid for 7 days).
type CateringQuoteInput struct {
	Amount         float64             `json:"amount" validate:"gt=0"`
	LineItems      []CateringQuoteLine `json:"lineItems" validate:"max=30,dive"`
	DepositPercent *int                `json:"depositPercent,omitempty" validate:"omitempty,min=1,max=100"`
	Terms          *string             `json:"terms,omitempty" validate:"omitempty,max=5000"`
	Message        *string             `json:"message,omitempty" validate:"omitempty,max=1000"`
	ValidDays      *int                `json:"validDays,omitempty" validate:"omitempty,min=1,max=60"`
}

type CateringPayment struct {
	ID         int       `json:"id"`
	RequestID  int       `json:"requestId"`
	Amount     float64   `json:"amount"`
	Method     string    `json:"method"`
	Reference  *string   `json:"reference,omitempty"`
	Note       *string   `json:"note,omitempty"`
	RecordedBy *int      `json:"recordedBy,omitempty"`
	PaidAt     time.Time `json:"paidAt"`
	CreatedAt  time.Time `json:"createdAt"`
}

// CateringPaymentInput is a payment the admin received in person or by
// direct bank transfer.
type CateringPaymentInput struct {
	Amount    float64 `json:"amount" validate:"gt=0"`
	Method    string  `json:"method" validate:"required,oneof=cash bank_transfer"`
	Reference *string `json:"reference,omitempty" validate:"omitempty,max=100"`
	Note      *string `json:"note,omitempty" validate:"omitempty,max=500"`
	PaidOn    *string `json:"paidOn,omitempty"` // YYYY-MM-DD; today when omitted
}

// CateringInvoice is one issued version of a booking's invoice.
type CateringInvoice struct {
	ID            int       `json:"id"`
	RequestID     int       `json:"requestId"`
	QuoteID       int       `json:"quoteId"`
	InvoiceNumber string    `json:"invoiceNumber"`
	Version       int       `json:"version"`
	Reason        string    `json:"reason"`
	Total         float64   `json:"total"`
	DepositDue    float64   `json:"depositDue"`
	AmountPaid    float64   `json:"amountPaid"`
	Balance       float64   `json:"balance"`
	IssuedAt      time.Time `json:"issuedAt"`
}

type CateringReceipt struct {
	ID            int       `json:"id"`
	RequestID     int       `json:"requestId"`
	ReceiptNumber string    `json:"receiptNumber"`
	TotalPaid     float64   `json:"totalPaid"`
	IssuedAt      time.Time `json:"issuedAt"`
}

type CateringStatusChange struct {
	ID         int       `json:"id"`
	RequestID  int       `json:"requestId"`
	FromStatus *string   `json:"fromStatus,omitempty"`
	ToStatus   string    `json:"toStatus"`
	ActorType  string    `json:"actorType"`
	ActorID    *int      `json:"actorId,omitempty"`
	Note       *string   `json:"note,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

// CateringPaymentCharge is a Paystack dedicated account awaiting a transfer.
type CateringPaymentCharge struct {
	ID            int        `json:"-"`
	RequestID     int        `json:"-"`
	Reference     string     `json:"reference"`
	Amount        float64    `json:"amount"`
	AccountNumber *string    `json:"accountNumber,omitempty"`
	AccountName   *string    `json:"accountName,omitempty"`
	BankName      *string    `json:"bankName,omitempty"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"createdAt"`
}

// CateringAccount summarises the money on an accepted booking, always
// calculated from the accepted quote and the payment ledger.
type CateringAccount struct {
	Total            float64 `json:"total"`
	DepositDue       float64 `json:"depositDue"`
	AmountPaid       float64 `json:"amountPaid"`
	Balance          float64 `json:"balance"`
	DepositRemaining float64 `json:"depositRemaining"`
}

// CateringContract is a request with everything that has happened to it.
type CateringContract struct {
	Request       *CateringRequest       `json:"request"`
	Quote         *CateringQuote         `json:"quote,omitempty"`
	Quotes        []CateringQuote        `json:"quotes,omitempty"`
	Account       *CateringAccount       `json:"account,omitempty"`
	Payments      []CateringPayment      `json:"payments"`
	Invoice       *CateringInvoice       `json:"invoice,omitempty"`
	Invoices      []CateringInvoice      `json:"invoices,omitempty"`
	Receipt       *CateringReceipt       `json:"receipt,omitempty"`
	PendingCharge *CateringPaymentCharge `json:"pendingCharge,omitempty"`
	History       []CateringStatusChange `json:"history,omitempty"`
}

// CateringCustomerAuth identifies the customer: by the one-time link token
// from their quote email, or by their reference and phone number.
type CateringCustomerAuth struct {
	Token     string `json:"token,omitempty" validate:"omitempty,max=100"`
	Reference string `json:"reference,omitempty" validate:"omitempty,max=30"`
	Phone     string `json:"phone,omitempty" validate:"omitempty,max=30"`
}

type CateringAcceptInput struct {
	CateringCustomerAuth
	AgreeToTerms bool `json:"agreeToTerms"`
}

type CateringPayOnlineInput struct {
	CateringCustomerAuth
	Portion string `json:"portion" validate:"required,oneof=deposit balance"`
}
