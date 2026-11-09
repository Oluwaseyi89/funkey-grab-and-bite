package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/requestctx"
)

// IdempotencyRepository stores the result of each client request made with an
// Idempotency-Key, so a retry returns the first result instead of repeating it.
type IdempotencyRepository struct {
	db *sql.DB
}

func NewIdempotencyRepository(db *sql.DB) *IdempotencyRepository {
	return &IdempotencyRepository{db: db}
}

// Begin claims key for a request with the given body hash. It returns
// IdempotencyNew when the caller should run the request: the key is new, its
// earlier record expired, or an earlier attempt with the same body crashed
// before committing (its lease ran out while still in progress). Otherwise it
// returns the existing record and whether to replay it, wait, or reject a
// different body.
func (r *IdempotencyRepository) Begin(ctx context.Context, scope, key, requestHash string, lease, ttl time.Duration) (models.IdempotencyOutcome, *models.IdempotencyRecord, error) {
	var claimed bool
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO idempotency_keys (scope, idempotency_key, request_hash, locked_until, expires_at)
		VALUES ($1, $2, $3, NOW() + make_interval(secs => $4), NOW() + make_interval(secs => $5))
		ON CONFLICT (scope, idempotency_key) DO UPDATE SET
			request_hash = EXCLUDED.request_hash,
			status = 'in_progress',
			entity_type = NULL,
			entity_id = NULL,
			response_status = NULL,
			response_body = NULL,
			locked_until = EXCLUDED.locked_until,
			created_at = NOW(),
			completed_at = NULL,
			expires_at = EXCLUDED.expires_at
		WHERE idempotency_keys.expires_at < NOW()
		   OR (idempotency_keys.status = 'in_progress'
		       AND idempotency_keys.locked_until < NOW()
		       AND idempotency_keys.request_hash = EXCLUDED.request_hash)
		RETURNING true
	`, scope, key, requestHash, lease.Seconds(), ttl.Seconds()).Scan(&claimed)
	if err == nil {
		return models.IdempotencyNew, nil, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, nil, fmt.Errorf("claim idempotency key: %w", err)
	}

	record, err := r.get(ctx, scope, key)
	if err != nil {
		return 0, nil, err
	}
	if record == nil {
		// Expired and removed between the two statements; the caller can retry.
		return models.IdempotencyInProgress, nil, nil
	}
	switch {
	case record.RequestHash != requestHash:
		return models.IdempotencyMismatch, record, nil
	case record.Status == "completed":
		return models.IdempotencyReplay, record, nil
	default:
		return models.IdempotencyInProgress, record, nil
	}
}

func (r *IdempotencyRepository) get(ctx context.Context, scope, key string) (*models.IdempotencyRecord, error) {
	var record models.IdempotencyRecord
	var entityType, entityID sql.NullString
	var responseStatus sql.NullInt64
	err := r.db.QueryRowContext(ctx, `
		SELECT scope, idempotency_key, request_hash, status, entity_type, entity_id,
		       response_status, response_body
		FROM idempotency_keys
		WHERE scope = $1 AND idempotency_key = $2
	`, scope, key).Scan(
		&record.Scope, &record.Key, &record.RequestHash, &record.Status,
		&entityType, &entityID, &responseStatus, &record.ResponseBody,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read idempotency key: %w", err)
	}
	record.EntityType = entityType.String
	record.EntityID = entityID.String
	record.ResponseStatus = int(responseStatus.Int64)
	return &record, nil
}

// MarkCompleted records, inside the transaction making the change, that the
// request in ctx created or changed the given entity. Because it commits or
// rolls back with the change, a completed key always means the change was
// saved. It does nothing when the request has no idempotency key.
func (r *IdempotencyRepository) MarkCompleted(ctx context.Context, exec Execer, entityType, entityID string) error {
	idem, ok := requestctx.IdempotencyFrom(ctx)
	if !ok {
		return nil
	}
	if exec == nil {
		exec = r.db
	}
	_, err := exec.ExecContext(ctx, `
		UPDATE idempotency_keys
		SET status = 'completed', entity_type = $3, entity_id = $4, completed_at = NOW()
		WHERE scope = $1 AND idempotency_key = $2
	`, idem.Scope, idem.Key, entityType, entityID)
	if err != nil {
		return fmt.Errorf("complete idempotency key: %w", err)
	}
	return nil
}

// SaveResponse stores the response sent for a successful request so retries
// get exactly the same answer, and marks the key completed.
func (r *IdempotencyRepository) SaveResponse(ctx context.Context, scope, key string, status int, body []byte) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE idempotency_keys
		SET status = 'completed', response_status = $3, response_body = $4,
		    completed_at = COALESCE(completed_at, NOW())
		WHERE scope = $1 AND idempotency_key = $2
	`, scope, key, status, body)
	if err != nil {
		return fmt.Errorf("save idempotent response: %w", err)
	}
	return nil
}

// Release frees a key whose request failed without changing anything, so the
// client can retry with it. A key the change already completed is kept.
func (r *IdempotencyRepository) Release(ctx context.Context, scope, key string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM idempotency_keys
		WHERE scope = $1 AND idempotency_key = $2 AND status = 'in_progress'
	`, scope, key)
	if err != nil {
		return fmt.Errorf("release idempotency key: %w", err)
	}
	return nil
}

// DeleteExpired removes keys past their expiry and returns how many.
func (r *IdempotencyRepository) DeleteExpired(ctx context.Context) (int64, error) {
	result, err := r.db.ExecContext(ctx, `DELETE FROM idempotency_keys WHERE expires_at < NOW()`)
	if err != nil {
		return 0, fmt.Errorf("delete expired idempotency keys: %w", err)
	}
	return result.RowsAffected()
}
