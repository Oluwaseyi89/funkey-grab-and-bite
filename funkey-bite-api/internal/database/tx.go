package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/lib/pq"
)

// maxTxAttempts bounds how often WithTx retries a transaction that Postgres
// aborted because of a serialization failure or deadlock.
const maxTxAttempts = 3

// WithTx runs fn in a transaction: it commits when fn returns nil and rolls
// back when fn returns an error or panics, so every write in fn is saved or
// none is.
//
// Postgres asks clients to retry transactions it aborts for a serialization
// failure or deadlock, so those are retried up to maxTxAttempts times. fn must
// therefore only change the database through tx and have no other side
// effects, such as calling Paystack or sending email.
func WithTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error {
	var err error
	for attempt := 1; attempt <= maxTxAttempts; attempt++ {
		err = runTx(ctx, db, fn)
		if err == nil || !isRetryableTxError(err) || attempt == maxTxAttempts {
			return err
		}
		select {
		case <-ctx.Done():
			return err
		case <-time.After(retryDelay(attempt)):
		}
	}
	return err
}

func runTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

// isRetryableTxError reports a serialization failure (40001) or deadlock (40P01).
func isRetryableTxError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "40001" || pqErr.Code == "40P01"
	}
	return false
}

// retryDelay grows with each attempt, with jitter so competing transactions
// don't retry in lockstep.
func retryDelay(attempt int) time.Duration {
	base := time.Duration(attempt) * 25 * time.Millisecond
	return base + time.Duration(rand.Int63n(int64(base)))
}
