package models

import (
	"encoding/json"
	"time"
)

// OnlinePayment is an online (transfer or card) order as the admin's
// payments page shows it.
type OnlinePayment struct {
	OrderID          int        `json:"orderId"`
	OrderNumber      string     `json:"orderNumber"`
	CustomerName     string     `json:"customerName"`
	CustomerPhone    string     `json:"customerPhone"`
	TotalAmount      float64    `json:"totalAmount"`
	Status           string     `json:"status"`
	PaymentMethod    string     `json:"paymentMethod"`
	PaymentStatus    string     `json:"paymentStatus"`
	PaymentReference *string    `json:"paymentReference,omitempty"`
	ReceivedAmount   *float64   `json:"receivedAmount,omitempty"`
	ReviewReason     *string    `json:"reviewReason,omitempty"`
	ReferenceCount   int        `json:"referenceCount"`
	ExpiresAt        *time.Time `json:"expiresAt,omitempty"`
	PaidAt           *time.Time `json:"paidAt,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
}

// OnlinePaymentFilter narrows the payments list; empty fields match all.
// Search matches the order number, customer name, phone or a reference.
type OnlinePaymentFilter struct {
	PaymentStatus string
	Search        string
}

// PaymentEventRecord is a Paystack event as received, for payments that
// matched no order.
type PaymentEventRecord struct {
	ID         int64     `json:"id"`
	Event      string    `json:"event"`
	Reference  string    `json:"reference"`
	Amount     *float64  `json:"amount,omitempty"`
	Currency   *string   `json:"currency,omitempty"`
	Outcome    string    `json:"outcome"`
	ReceivedAt time.Time `json:"receivedAt"`
}

// Payment review decisions the admin can make.
const (
	ReviewDecisionAccept   = "accept"
	ReviewDecisionRefund   = "refund"
	ReviewDecisionRefunded = "refunded"
)

// AuditRecord is an audit log entry as the admin sees it.
type AuditRecord struct {
	ID             int64           `json:"id"`
	OccurredAt     time.Time       `json:"occurredAt"`
	ActorType      string          `json:"actorType"`
	ActorID        *string         `json:"actorId,omitempty"`
	Action         string          `json:"action"`
	EntityType     string          `json:"entityType"`
	EntityID       *string         `json:"entityId,omitempty"`
	RequestID      *string         `json:"requestId,omitempty"`
	IdempotencyKey *string         `json:"idempotencyKey,omitempty"`
	IPAddress      *string         `json:"ipAddress,omitempty"`
	Before         json.RawMessage `json:"before,omitempty"`
	After          json.RawMessage `json:"after,omitempty"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
}

// AuditFilter narrows the audit log; empty fields match all.
type AuditFilter struct {
	EntityType string
	EntityID   string
	Action     string
	ActorType  string
}
