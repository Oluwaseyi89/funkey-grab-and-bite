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

// List returns audit entries matching filter, newest first, and how many
// match.
func (r *AuditRepository) List(ctx context.Context, filter models.AuditFilter, limit, offset int) ([]models.AuditRecord, int, error) {
	where := `
		WHERE ($1 = '' OR entity_type = $1)
		  AND ($2 = '' OR entity_id = $2)
		  AND ($3 = '' OR action = $3)
		  AND ($4 = '' OR actor_type = $4)
	`
	args := []any{filter.EntityType, filter.EntityID, filter.Action, filter.ActorType}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_log`+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count audit entries: %w", err)
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, occurred_at, actor_type, actor_id, action, entity_type, entity_id,
		       request_id, idempotency_key, ip_address, before_state, after_state, metadata
		FROM audit_log`+where+`
		ORDER BY occurred_at DESC, id DESC
		LIMIT $5 OFFSET $6
	`, append(args, limit, offset)...)
	if err != nil {
		return nil, 0, fmt.Errorf("list audit entries: %w", err)
	}
	defer rows.Close()

	entries := []models.AuditRecord{}
	for rows.Next() {
		var e models.AuditRecord
		var actorID, entityID, requestID, idemKey, ip sql.NullString
		var before, after, metadata []byte
		if err := rows.Scan(&e.ID, &e.OccurredAt, &e.ActorType, &actorID, &e.Action, &e.EntityType, &entityID,
			&requestID, &idemKey, &ip, &before, &after, &metadata); err != nil {
			return nil, 0, fmt.Errorf("read audit entry: %w", err)
		}
		e.ActorID, e.EntityID, e.RequestID = nullableString(actorID), nullableString(entityID), nullableString(requestID)
		e.IdempotencyKey, e.IPAddress = nullableString(idemKey), nullableString(ip)
		e.Before, e.After, e.Metadata = before, after, metadata
		entries = append(entries, e)
	}
	return entries, total, rows.Err()
}

func nullableString(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	return &v.String
}
