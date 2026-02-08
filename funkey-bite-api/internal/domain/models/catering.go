package models

import "time"

type CateringStatus string

const (
	CateringStatusPending   CateringStatus = "pending"
	CateringStatusConfirmed CateringStatus = "confirmed"
	CateringStatusDeclined  CateringStatus = "declined"
	CateringStatusCompleted CateringStatus = "completed"
)

type CateringRequest struct {
	ID              int            `json:"id" db:"id"`
	UserID          *int           `json:"userId,omitempty" db:"user_id"`
	EventName       *string        `json:"eventName,omitempty" db:"event_name"`
	ContactName     string         `json:"contactName" db:"contact_name"`
	ContactPhone    string         `json:"contactPhone" db:"contact_phone"`
	ContactEmail    *string        `json:"contactEmail,omitempty" db:"contact_email"`
	EventDate       string         `json:"eventDate" db:"event_date"`
	EventTime       *string        `json:"eventTime,omitempty" db:"event_time"`
	GuestCount      int            `json:"guestCount" db:"guest_count"`
	EventType       string         `json:"eventType" db:"event_type"`
	Package         *string        `json:"package,omitempty" db:"package"`
	Budget          *float64       `json:"budget,omitempty" db:"budget"`
	SpecialRequests *string        `json:"specialRequests,omitempty" db:"special_requests"`
	Status          CateringStatus `json:"status" db:"status"`
	CreatedAt       time.Time      `json:"createdAt" db:"created_at"`
}

type CateringRequestInput struct {
	EventName       *string  `json:"eventName,omitempty" validate:"omitempty,min=2,max=200"`
	ContactName     string   `json:"contactName" validate:"required,min=2,max=200"`
	ContactPhone    string   `json:"contactPhone" validate:"required,min=10,max=20"`
	ContactEmail    *string  `json:"contactEmail,omitempty" validate:"omitempty,email"`
	EventDate       string   `json:"eventDate" validate:"required"`
	EventTime       *string  `json:"eventTime,omitempty"`
	GuestCount      int      `json:"guestCount" validate:"required,min=1,max=1000"`
	EventType       string   `json:"eventType" validate:"required"`
	Package         *string  `json:"package,omitempty"`
	Budget          *float64 `json:"budget,omitempty" validate:"omitempty,min=0"`
	SpecialRequests *string  `json:"specialRequests,omitempty" validate:"omitempty,max=1000"`
}

type CateringPackage struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	PricePerPerson float64  `json:"pricePerPerson"`
	MinGuests      int      `json:"minGuests"`
	MaxGuests      *int     `json:"maxGuests,omitempty"`
	Includes       []string `json:"includes"`
}
