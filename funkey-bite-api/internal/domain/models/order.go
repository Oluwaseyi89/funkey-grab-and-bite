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

type Order struct {
	ID                 int         `json:"id" db:"id"`
	OrderNumber        string      `json:"orderNumber" db:"order_number"`
	UserID             *int        `json:"userId,omitempty" db:"user_id"`
	CustomerID         string      `json:"customerId,omitempty" db:"customer_id"`
	CustomerName       string      `json:"customerName" db:"customer_name"`
	CustomerPhone      string      `json:"customerPhone" db:"customer_phone"`
	CustomerEmail      *string     `json:"customerEmail,omitempty" db:"customer_email"`
	OrderType          OrderType   `json:"orderType" db:"order_type"`
	Status             OrderStatus `json:"status" db:"status"`
	TotalAmount        float64     `json:"totalAmount" db:"total_amount"`
	Notes              *string     `json:"notes,omitempty" db:"notes"`
	PickupTime         *time.Time  `json:"pickupTime,omitempty" db:"pickup_time"`
	EstimatedReadyTime *time.Time  `json:"estimatedReadyTime,omitempty" db:"estimated_ready_time"`
	CreatedAt          time.Time   `json:"createdAt" db:"created_at"`
	Items              []OrderItem `json:"items"`
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
	PromotionCode *string            `json:"promotionCode,omitempty" validate:"omitempty,min=3,max=50"`
	Notes         *string            `json:"notes,omitempty" db:"notes"`
	PickupTime    *time.Time         `json:"pickupTime,omitempty" db:"pickup_time"`
	Items         []OrderItemRequest `json:"items" validate:"required,min=1"`
}

type OrderWithAuth struct {
	CustomerName  string             `json:"customerName" validate:"required"`
	CustomerPhone string             `json:"customerPhone" validate:"required"`
	CustomerEmail *string            `json:"customerEmail,omitempty" validate:"omitempty,email"`
	OrderType     OrderType          `json:"orderType" validate:"required,oneof=pickup delivery catering"`
	PromotionCode *string            `json:"promotionCode,omitempty" validate:"omitempty,min=3,max=50"`
	Notes         *string            `json:"notes,omitempty"`
	PickupTime    *time.Time         `json:"pickupTime,omitempty"`
	Items         []OrderItemRequest `json:"items" validate:"required,min=1"`
	Password      *string            `json:"password,omitempty" validate:"omitempty,min=8"`
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
