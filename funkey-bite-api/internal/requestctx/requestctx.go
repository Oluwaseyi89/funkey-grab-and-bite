// Package requestctx carries who is making a request and its idempotency key
// from the HTTP layer down to services and repositories, so audit entries and
// idempotency records can be written in the same transaction as the change.
package requestctx

import "context"

// Actor types, matching the audit_log.actor_type check.
const (
	ActorAdmin    = "admin"
	ActorCustomer = "customer"
	ActorSystem   = "system"
	ActorPaystack = "paystack"
)

// Actor is who made a request. ID is empty for guests and the system.
type Actor struct {
	Type      string
	ID        string
	IP        string
	RequestID string
}

// Idempotency identifies a client request made with an Idempotency-Key.
type Idempotency struct {
	Scope string
	Key   string
}

type actorKey struct{}
type idempotencyKey struct{}

// WithActor returns ctx carrying actor.
func WithActor(ctx context.Context, actor Actor) context.Context {
	return context.WithValue(ctx, actorKey{}, actor)
}

// ActorFrom returns the actor in ctx, or the system when there is none, such
// as for background work.
func ActorFrom(ctx context.Context) Actor {
	if actor, ok := ctx.Value(actorKey{}).(Actor); ok && actor.Type != "" {
		return actor
	}
	return Actor{Type: ActorSystem}
}

// WithIdempotency returns ctx carrying the request's idempotency key.
func WithIdempotency(ctx context.Context, idem Idempotency) context.Context {
	return context.WithValue(ctx, idempotencyKey{}, idem)
}

// IdempotencyFrom returns the request's idempotency key, if it has one.
func IdempotencyFrom(ctx context.Context) (Idempotency, bool) {
	idem, ok := ctx.Value(idempotencyKey{}).(Idempotency)
	return idem, ok && idem.Key != ""
}
