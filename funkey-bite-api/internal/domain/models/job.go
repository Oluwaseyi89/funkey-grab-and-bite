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
