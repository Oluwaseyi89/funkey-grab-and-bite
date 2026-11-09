package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// RateLimitRepository counts requests per bucket per fixed window in
// Postgres, so limits hold across restarts and every API instance.
type RateLimitRepository struct {
	db *sql.DB
}

func NewRateLimitRepository(db *sql.DB) *RateLimitRepository {
	return &RateLimitRepository{db: db}
}

// Hit counts one request in bucket's window and returns the count so far,
// including this one.
func (r *RateLimitRepository) Hit(ctx context.Context, bucket string, windowStart time.Time) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO rate_limit_counters (bucket, window_start, request_count)
		VALUES ($1, $2, 1)
		ON CONFLICT (bucket, window_start)
		DO UPDATE SET request_count = rate_limit_counters.request_count + 1
		RETURNING request_count
	`, bucket, windowStart).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count rate limited request: %w", err)
	}
	return count, nil
}

// DeleteBefore removes windows that started before cutoff and returns how many.
func (r *RateLimitRepository) DeleteBefore(ctx context.Context, cutoff time.Time) (int64, error) {
	result, err := r.db.ExecContext(ctx, `DELETE FROM rate_limit_counters WHERE window_start < $1`, cutoff)
	if err != nil {
		return 0, fmt.Errorf("delete old rate limit windows: %w", err)
	}
	return result.RowsAffected()
}
