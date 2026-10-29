package models

import "time"

type OrderStatus string
type OrderType string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

const (
	OrderTypePickup   OrderType = "pickup"
	OrderTypeDelivery OrderType = "delivery"
	OrderTypeCatering OrderType = "catering"
)

type PaymentMethod string
type PaymentStatus string

const (
	// PaymentMethodCash is pay on delivery (or at pickup): the customer pays
	// when they receive the order, and the admin records the payment.
	PaymentMethodCash     PaymentMethod = "cash"
	PaymentMethodTransfer PaymentMethod = "transfer"
)

// How a pay-on-delivery payment was collected, recorded by the admin.
const (
	PaymentCollectedCash     = "cash"
	PaymentCollectedPOS      = "pos"
	PaymentCollectedTransfer = "transfer"
)

const (
	// PaymentStatusNotRequired was used for cash orders before pay on delivery
	// was tracked; new cash orders start as pending until the admin records payment.
	PaymentStatusNotRequired PaymentStatus = "not_required"
	PaymentStatusPending     PaymentStatus = "pending"
	PaymentStatusPaid        PaymentStatus = "paid"
	PaymentStatusFailed      PaymentStatus = "failed"
)

type Order struct {
	ID                   int           `json:"id" db:"id"`
	OrderNumber          string        `json:"orderNumber" db:"order_number"`
	UserID               *int          `json:"userId,omitempty" db:"user_id"`
	CustomerID           string        `json:"customerId,omitempty" db:"customer_id"`
	CustomerName         string        `json:"customerName" db:"customer_name"`
	CustomerPhone        string        `json:"customerPhone" db:"customer_phone"`
	CustomerEmail        *string       `json:"customerEmail,omitempty" db:"customer_email"`
	OrderType            OrderType     `json:"orderType" db:"order_type"`
	Status               OrderStatus   `json:"status" db:"status"`
	TotalAmount          float64       `json:"totalAmount" db:"total_amount"`
	Notes                *string       `json:"notes,omitempty" db:"notes"`
	PickupTime           *time.Time    `json:"pickupTime,omitempty" db:"pickup_time"`
	EstimatedReadyTime   *time.Time    `json:"estimatedReadyTime,omitempty" db:"estimated_ready_time"`
	CreatedAt            time.Time     `json:"createdAt" db:"created_at"`
	PaymentMethod        PaymentMethod `json:"paymentMethod" db:"payment_method"`
	PaymentStatus        PaymentStatus `json:"paymentStatus" db:"payment_status"`
	PaymentReference     *string       `json:"paymentReference,omitempty" db:"payment_reference"`
	PaymentAccountNumber *string       `json:"paymentAccountNumber,omitempty" db:"payment_account_number"`
	PaymentAccountName   *string       `json:"paymentAccountName,omitempty" db:"payment_account_name"`
	PaymentBankName      *string       `json:"paymentBankName,omitempty" db:"payment_bank_name"`
	PaymentExpiresAt     *time.Time    `json:"paymentExpiresAt,omitempty" db:"payment_expires_at"`
	PaymentPaidAt        *time.Time    `json:"paymentPaidAt,omitempty" db:"payment_paid_at"`
	PaymentCollectedVia  *string       `json:"paymentCollectedVia,omitempty" db:"payment_collected_via"`
	DeliveryAddress      *string       `json:"deliveryAddress,omitempty" db:"delivery_address"`
	InvoiceNumber        *string       `json:"invoiceNumber,omitempty" db:"invoice_number"`
	SubtotalAmount       *float64      `json:"subtotalAmount,omitempty" db:"subtotal_amount"`
	DeliveryFee          *float64      `json:"deliveryFee,omitempty" db:"delivery_fee"`
	TaxAmount            *float64      `json:"taxAmount,omitempty" db:"tax_amount"`
	DiscountAmount       *float64      `json:"discountAmount,omitempty" db:"discount_amount"`
	Items                []OrderItem   `json:"items"`
}

type OrderItem struct {
	ID                  int     `json:"-" db:"id"`
	OrderID             int     `json:"-" db:"order_id"`
	MenuItemID          int     `json:"menuItemId" db:"menu_item_id"`
	Name                string  `json:"name" db:"name"`
	Quantity            int     `json:"quantity" db:"quantity"`
	UnitPrice           float64 `json:"unitPrice" db:"unit_price"`
	SpecialInstructions *string `json:"specialInstructions,omitempty" db:"special_instructions"`
}

type OrderItemRequest struct {
	MenuItemID          int     `json:"menuItemId" validate:"required"`
	Name                string  `json:"name" validate:"required"`
	Quantity            int     `json:"quantity" validate:"required,min=1"`
	UnitPrice           float64 `json:"unitPrice" validate:"required,min=0"`
	SpecialInstructions *string `json:"specialInstructions,omitempty"`
}

type OrderRequest struct {
	CustomerName  string             `json:"customerName" validate:"required" db:"customer_name"`
	CustomerPhone string             `json:"customerPhone" validate:"required" db:"customer_phone"`
	CustomerEmail *string            `json:"customerEmail,omitempty" validate:"omitempty,email" db:"customer_email"`
	OrderType     OrderType          `json:"orderType" validate:"required,oneof=pickup delivery catering" db:"order_type"`
	PaymentMethod PaymentMethod      `json:"paymentMethod" validate:"required,oneof=cash transfer" db:"payment_method"`
	PromotionCode *string            `json:"promotionCode,omitempty" validate:"omitempty,min=3,max=50"`
	Notes         *string            `json:"notes,omitempty" db:"notes"`
	PickupTime    *time.Time         `json:"pickupTime,omitempty" db:"pickup_time"`
	Items         []OrderItemRequest `json:"items" validate:"required,min=1"`
}

// OrderWithAuth is the checkout request. Customers check out as guests: the
// contact details are required on every order and link it to a customer record
// by phone number; a signed-in customer's token links it to their account.
type OrderWithAuth struct {
	CustomerName    string             `json:"customerName" validate:"required,max=200"`
	CustomerPhone   string             `json:"customerPhone" validate:"required,max=30"`
	CustomerEmail   *string            `json:"customerEmail" validate:"required,email,max=200"`
	OrderType       OrderType          `json:"orderType" validate:"required,oneof=pickup delivery catering"`
	PaymentMethod   PaymentMethod      `json:"paymentMethod" validate:"required,oneof=cash transfer"`
	DeliveryAddress *string            `json:"deliveryAddress,omitempty" validate:"omitempty,max=500"`
	PromotionCode   *string            `json:"promotionCode,omitempty" validate:"omitempty,min=3,max=50"`
	Notes           *string            `json:"notes,omitempty" validate:"omitempty,max=1000"`
	PickupTime      *time.Time         `json:"pickupTime,omitempty"`
	Items           []OrderItemRequest `json:"items" validate:"required,min=1"`
}

// OrderLookupRequest identifies an order to the customer who placed it.
type OrderLookupRequest struct {
	OrderNumber string `json:"orderNumber" validate:"required,max=50"`
	Phone       string `json:"phone" validate:"required,max=30"`
}

// PaymentCollectionRequest records how a pay-on-delivery payment was received.
type PaymentCollectionRequest struct {
	CollectedVia string `json:"collectedVia" validate:"required,oneof=cash pos transfer"`
}

// OrderListFilter narrows the admin order list; empty fields match everything.
type OrderListFilter struct {
	Status        string
	PaymentStatus string
	PaymentMethod string
	UserID        *int
}

// Invoice is raised for pay-on-delivery orders so the customer knows what to
// pay on delivery and the admin has a record to collect against.
type Invoice struct {
	InvoiceNumber       string        `json:"invoiceNumber"`
	OrderNumber         string        `json:"orderNumber"`
	IssuedAt            time.Time     `json:"issuedAt"`
	Business            InvoiceParty  `json:"business"`
	Customer            InvoiceParty  `json:"customer"`
	OrderType           OrderType     `json:"orderType"`
	Items               []InvoiceLine `json:"items"`
	Subtotal            float64       `json:"subtotal"`
	DeliveryFee         float64       `json:"deliveryFee"`
	Tax                 float64       `json:"tax"`
	Discount            float64       `json:"discount"`
	Total               float64       `json:"total"`
	PaymentMethod       PaymentMethod `json:"paymentMethod"`
	PaymentStatus       PaymentStatus `json:"paymentStatus"`
	PaymentCollectedVia *string       `json:"paymentCollectedVia,omitempty"`
	PaidAt              *time.Time    `json:"paidAt,omitempty"`
}

type InvoiceParty struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email,omitempty"`
	Address string `json:"address,omitempty"`
}

type InvoiceLine struct {
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unitPrice"`
	Total     float64 `json:"total"`
}

func (o *Order) ToFrontendFormat() map[string]interface{} {
	return map[string]interface{}{
		"id":            o.ID,
		"orderNumber":   o.OrderNumber,
		"customerName":  o.CustomerName,
		"customerPhone": o.CustomerPhone,
		"customerEmail": o.CustomerEmail,
		"orderType":     string(o.OrderType),
		"status":        string(o.Status),
		"totalAmount":   o.TotalAmount,
		"notes":         o.Notes,
		"pickupTime":    o.PickupTime,
		"createdAt":     o.CreatedAt,
		"items":         o.Items,
	}
}

func ParseOrderStatus(status string) OrderStatus {
	switch status {
	case "confirmed":
		return OrderStatusConfirmed
	case "preparing":
		return OrderStatusPreparing
	case "ready":
		return OrderStatusReady
	case "completed":
		return OrderStatusCompleted
	case "cancelled":
		return OrderStatusCancelled
	default:
		return OrderStatusPending
	}
}

func ParseOrderType(orderType string) OrderType {
	switch orderType {
	case "delivery":
		return OrderTypeDelivery
	case "catering":
		return OrderTypeCatering
	default:
		return OrderTypePickup
	}
}
