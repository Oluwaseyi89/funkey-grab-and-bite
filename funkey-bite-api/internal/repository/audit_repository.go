package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/requestctx"
)

// Execer is satisfied by both *sql.DB and *sql.Tx, so a write can join the
// caller's transaction.
type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// AuditRepository appends to the audit log. Entries are never updated or
// deleted.
type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// Record appends entry using exec, which should be the transaction making the
// change so the log can never disagree with the data; nil writes on its own.
// The actor, IP, request ID and idempotency key come from ctx when the entry
// leaves them empty.
func (r *AuditRepository) Record(ctx context.Context, exec Execer, entry models.AuditEntry) error {
	if exec == nil {
		exec = r.db
	}
	if entry.Action == "" || entry.EntityType == "" {
		return fmt.Errorf("audit entry needs an action and an entity type")
	}

	actor := requestctx.ActorFrom(ctx)
	if entry.ActorType == "" {
		entry.ActorType = actor.Type
		if entry.ActorID == "" {
			entry.ActorID = actor.ID
		}
	}
	if entry.IPAddress == "" {
		entry.IPAddress = actor.IP
	}
	if entry.RequestID == "" {
		entry.RequestID = actor.RequestID
	}
	if entry.IdempotencyKey == "" {
		if idem, ok := requestctx.IdempotencyFrom(ctx); ok {
			entry.IdempotencyKey = idem.Key
		}
	}

	before, err := jsonOrNull(entry.Before)
	if err != nil {
		return fmt.Errorf("encode audit before state: %w", err)
	}
	after, err := jsonOrNull(entry.After)
	if err != nil {
		return fmt.Errorf("encode audit after state: %w", err)
	}
	metadata, err := jsonOrNull(entry.Metadata)
	if err != nil {
		return fmt.Errorf("encode audit metadata: %w", err)
	}

	_, err = exec.ExecContext(ctx, `
		INSERT INTO audit_log (
			actor_type, actor_id, action, entity_type, entity_id, request_id,
			idempotency_key, ip_address, before_state, after_state, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`,
		entry.ActorType, nullIfEmpty(entry.ActorID), entry.Action, entry.EntityType,
		nullIfEmpty(entry.EntityID), nullIfEmpty(entry.RequestID), nullIfEmpty(entry.IdempotencyKey),
		nullIfEmpty(entry.IPAddress), before, after, metadata,
	)
	if err != nil {
		return fmt.Errorf("write audit entry: %w", err)
	}
	return nil
}

func jsonOrNull(value any) (any, error) {
	if value == nil {
		return nil, nil
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return string(encoded), nil
}

func nullIfEmpty(value string) any {
	if value == "" {
		return nil
	}
	return value
}
