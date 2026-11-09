package models

import (
	"encoding/json"
	"time"
)

// NewJob is a background job to queue.
type NewJob struct {
	Kind    string
	Payload any
	// RunAt delays the first attempt; zero means as soon as possible.
	RunAt time.Time
	// MaxAttempts bounds retries; zero uses the default of 5.
	MaxAttempts int
	// DedupeKey, when set, stops the same job being queued twice.
	DedupeKey string
	// Sensitive jobs have their payload wiped once they succeed.
	Sensitive bool
}

// JobSummary is a job as the admin sees it. The payload is never included.
type JobSummary struct {
	ID          int64      `json:"id"`
	Kind        string     `json:"kind"`
	Status      string     `json:"status"`
	Attempts    int        `json:"attempts"`
	MaxAttempts int        `json:"maxAttempts"`
	RunAt       time.Time  `json:"runAt"`
	LastError   *string    `json:"lastError,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}

// EmailJobPayload is a rendered email to send.
type EmailJobPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	HTML    string `json:"html"`
	Text    string `json:"text"`
}

// Job is a queued job claimed by a worker. Attempts includes the current one.
type Job struct {
	ID          int64
	Kind        string
	Payload     json.RawMessage
	Attempts    int
	MaxAttempts int
}

// OrderJobPayload identifies the order a job is for.
type OrderJobPayload struct {
	OrderID int `json:"orderId"`
}
