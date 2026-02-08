package models

import "time"

type Notification struct {
	ID            int        `json:"id" db:"id"`
	UserID        int        `json:"userId" db:"user_id"`
	Type          string     `json:"type" db:"type"`
	Title         string     `json:"title" db:"title"`
	Message       string     `json:"message" db:"message"`
	IsRead        bool       `json:"isRead" db:"is_read"`
	ReferenceID   *int       `json:"referenceId,omitempty" db:"reference_id"`
	ReferenceType string     `json:"referenceType,omitempty" db:"reference_type"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	ReadAt        *time.Time `json:"readAt,omitempty" db:"read_at"`
}
