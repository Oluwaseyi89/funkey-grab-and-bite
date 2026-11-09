package models

// AuditEntry records one change: who made it, what changed and the state
// before and after. Fields left empty are filled from the request where
// possible (actor, IP, request ID and idempotency key).
type AuditEntry struct {
	ActorType      string
	ActorID        string
	Action         string
	EntityType     string
	EntityID       string
	RequestID      string
	IdempotencyKey string
	IPAddress      string
	// Before, After and Metadata are stored as JSON; nil leaves them empty.
	Before   any
	After    any
	Metadata any
}
