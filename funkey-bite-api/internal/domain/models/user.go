package models

import (
	"time"
)

type User struct {
	ID           int        `json:"id" db:"id"`
	Phone        string     `json:"phone" db:"phone"`
	Email        *string    `json:"email,omitempty" db:"email"`
	FullName     string     `json:"fullName" db:"full_name"`
	PasswordHash string     `json:"-" db:"password_hash"`
	IsVerified   bool       `json:"isVerified" db:"is_verified"`
	IsActive     bool       `json:"isActive" db:"is_active"`
	LastLogin    *time.Time `json:"lastLogin,omitempty" db:"last_login"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
}

// AdminCustomer is a customer as the admin sees them, with their order
// totals. Cancelled orders are left out; totalSpent counts paid orders only.
type AdminCustomer struct {
	User
	OrderCount       int        `json:"orderCount"`
	TotalSpent       float64    `json:"totalSpent"`
	FirstOrderAt     *time.Time `json:"firstOrderAt,omitempty"`
	LastOrderAt      *time.Time `json:"lastOrderAt,omitempty"`
	FavoriteCategory *string    `json:"favoriteCategory,omitempty"`
}

// Customer list status filters.
const (
	CustomerStatusActive     = "active"
	CustomerStatusInactive   = "inactive"
	CustomerStatusVerified   = "verified"
	CustomerStatusUnverified = "unverified"
)

// CustomerListFilter narrows the admin customer list; empty fields match
// everything. Search matches name, phone or email.
type CustomerListFilter struct {
	Search string
	Status string
}

type UserRegistration struct {
	Phone string `json:"phone" validate:"required,min=10,max=20"`
	// Email is needed to confirm the account and to reset its password.
	Email    string `json:"email" validate:"required,email,max=200"`
	FullName string `json:"fullName" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegistrationConfirmation completes a registration for a number with guest
// orders, with the code emailed to the address on those orders.
type RegistrationConfirmation struct {
	UserRegistration
	Code string `json:"code" validate:"required,len=6,numeric"`
}

// PasswordResetRequest asks for a reset code by email.
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email,max=200"`
}

// PasswordReset sets a new password with an emailed code.
type PasswordReset struct {
	Email    string `json:"email" validate:"required,email,max=200"`
	Code     string `json:"code" validate:"required,len=6,numeric"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserLogin struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
	IsNewUser    bool   `json:"isNewUser,omitempty"`
}

type ProfileUpdate struct {
	FullName *string `json:"fullName,omitempty"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
}
