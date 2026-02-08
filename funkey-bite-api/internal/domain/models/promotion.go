package models

import (
	"time"
)

type PromotionType string

const (
	PromotionTypePercentage PromotionType = "percentage"
	PromotionTypeFixed      PromotionType = "fixed"
	PromotionTypeBOGO       PromotionType = "bogo"
)

type PromotionStatus string

const (
	PromotionStatusActive   PromotionStatus = "active"
	PromotionStatusInactive PromotionStatus = "inactive"
	PromotionStatusExpired  PromotionStatus = "expired"
)

type Promotion struct {
	ID             int           `json:"id" db:"id"`
	Code           string        `json:"code" db:"code"`
	Title          string        `json:"title" db:"title"`
	Description    string        `json:"description" db:"description"`
	PromotionType  PromotionType `json:"promotionType" db:"promotion_type"`
	DiscountValue  float64       `json:"discountValue" db:"discount_value"`
	MaxDiscount    *float64      `json:"maxDiscount,omitempty" db:"max_discount"`
	MinOrderAmount *float64      `json:"minOrderAmount,omitempty" db:"min_order_amount"`
	ValidFrom      time.Time     `json:"validFrom" db:"valid_from"`
	ValidUntil     time.Time     `json:"validUntil" db:"valid_until"`
	UsageLimit     *int          `json:"usageLimit,omitempty" db:"usage_limit"`
	UsedCount      int           `json:"usedCount" db:"used_count"`
	IsActive       bool          `json:"isActive" db:"is_active"`
	CreatedAt      time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time     `json:"updatedAt" db:"updated_at"`
}

type PromotionCreate struct {
	Code           string        `json:"code" validate:"required,min=3,max=50"`
	Title          string        `json:"title" validate:"required,min=3,max=200"`
	Description    string        `json:"description" validate:"max=500"`
	PromotionType  PromotionType `json:"promotionType" validate:"required,oneof=percentage fixed bogo"`
	DiscountValue  float64       `json:"discountValue" validate:"required,min=0"`
	MaxDiscount    *float64      `json:"maxDiscount,omitempty" validate:"omitempty,min=0"`
	MinOrderAmount *float64      `json:"minOrderAmount,omitempty" validate:"omitempty,min=0"`
	ValidFrom      time.Time     `json:"validFrom" validate:"required"`
	ValidUntil     time.Time     `json:"validUntil" validate:"required"`
	UsageLimit     *int          `json:"usageLimit,omitempty" validate:"omitempty,min=1"`
	IsActive       bool          `json:"isActive"`
}

type PromotionUpdate struct {
	Title          *string        `json:"title,omitempty" validate:"omitempty,min=3,max=200"`
	Description    *string        `json:"description,omitempty" validate:"omitempty,max=500"`
	PromotionType  *PromotionType `json:"promotionType,omitempty" validate:"omitempty,oneof=percentage fixed bogo"`
	DiscountValue  *float64       `json:"discountValue,omitempty" validate:"omitempty,min=0"`
	MaxDiscount    *float64       `json:"maxDiscount,omitempty" validate:"omitempty,min=0"`
	MinOrderAmount *float64       `json:"minOrderAmount,omitempty" validate:"omitempty,min=0"`
	ValidUntil     *time.Time     `json:"validUntil,omitempty"`
	UsageLimit     *int           `json:"usageLimit,omitempty" validate:"omitempty,min=1"`
	IsActive       *bool          `json:"isActive,omitempty"`
}

type PromotionValidation struct {
	IsValid     bool    `json:"isValid"`
	Message     string  `json:"message,omitempty"`
	Discount    float64 `json:"discount,omitempty"`
	PromotionID int     `json:"promotionId,omitempty"`
}

type PromotionUsage struct {
	ID              int       `json:"id" db:"id"`
	PromotionID     int       `json:"promotionId" db:"promotion_id"`
	OrderID         int       `json:"orderId" db:"order_id"`
	CustomerID      *int      `json:"customerId,omitempty" db:"customer_id"`
	DiscountApplied float64   `json:"discountApplied" db:"discount_applied"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
}
