package models

import (
	"time"
)

type User struct {
	ID           int        `json:"id" db:"id"`
	Phone        string     `json:"phone" db:"phone"`
	Email        *string    `json:"email,omitempty" db:"email"`
	FullName     string     `json:"full_name" db:"full_name"`
	PasswordHash string     `json:"-" db:"password_hash"`
	IsVerified   bool       `json:"is_verified" db:"is_verified"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	LastLogin    *time.Time `json:"last_login,omitempty" db:"last_login"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type UserRegistration struct {
	Phone    string `json:"phone" validate:"required,min=10,max=20"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserLogin struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IsNewUser    bool   `json:"is_new_user,omitempty"`
}

type ProfileUpdate struct {
	FullName *string `json:"full_name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
}
