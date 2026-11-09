package models

// IdempotencyOutcome is what to do with a request carrying an Idempotency-Key.
type IdempotencyOutcome int

const (
	// IdempotencyNew: first time this key is seen (or a crashed attempt is
	// being taken over), so the request should run.
	IdempotencyNew IdempotencyOutcome = iota
	// IdempotencyReplay: the request already completed; return its result.
	IdempotencyReplay
	// IdempotencyInProgress: another attempt with this key is still running.
	IdempotencyInProgress
	// IdempotencyMismatch: the key was used for a different request body.
	IdempotencyMismatch
)

// IdempotencyRecord is a stored request result. EntityType and EntityID name
// what the request created or changed; ResponseBody is the response as first
// sent, when it was saved.
type IdempotencyRecord struct {
	Scope          string
	Key            string
	RequestHash    string
	Status         string
	EntityType     string
	EntityID       string
	ResponseStatus int
	ResponseBody   []byte
}
