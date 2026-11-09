package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/database"
)

// Purposes of a one-time account code.
const (
	CodePurposeClaimAccount  = "claim_account"
	CodePurposePasswordReset = "password_reset"
)

// AccountCodeRepository stores one-time codes emailed to customers, as
// hashes only.
type AccountCodeRepository struct {
	db *sql.DB
}

func NewAccountCodeRepository(db *sql.DB) *AccountCodeRepository {
	return &AccountCodeRepository{db: db}
}

// Create stores a new code for user and purpose. Earlier unused codes for the
// same purpose stop working, so only the latest email counts.
func (r *AccountCodeRepository) Create(ctx context.Context, userID int, purpose, email, codeHash string, ttl time.Duration) error {
	return database.WithTx(ctx, r.db, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `
			UPDATE account_codes SET consumed_at = NOW()
			WHERE user_id = $1 AND purpose = $2 AND consumed_at IS NULL
		`, userID, purpose); err != nil {
			return fmt.Errorf("retire earlier codes: %w", err)
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO account_codes (user_id, purpose, email, code_hash, expires_at)
			VALUES ($1, $2, $3, $4, NOW() + make_interval(secs => $5))
		`, userID, purpose, email, codeHash, ttl.Seconds()); err != nil {
			return fmt.Errorf("store code: %w", err)
		}
		return nil
	})
}

// Use checks codeHash against the user's current code for purpose and, if it
// matches, uses it up. A wrong code counts as a try; after maxAttempts the
// code stops working. It returns false for a wrong, used, expired or
// exhausted code.
func (r *AccountCodeRepository) Use(ctx context.Context, userID int, purpose, codeHash string, maxAttempts int) (bool, error) {
	ok := false
	err := database.WithTx(ctx, r.db, func(tx *sql.Tx) error {
		ok = false
		var id int64
		var stored string
		var attempts int
		err := tx.QueryRowContext(ctx, `
			SELECT id, code_hash, attempts FROM account_codes
			WHERE user_id = $1 AND purpose = $2 AND consumed_at IS NULL AND expires_at > NOW()
			ORDER BY created_at DESC LIMIT 1
			FOR UPDATE
		`, userID, purpose).Scan(&id, &stored, &attempts)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read code: %w", err)
		}
		if attempts >= maxAttempts {
			return nil
		}
		if stored != codeHash {
			_, err := tx.ExecContext(ctx, `UPDATE account_codes SET attempts = attempts + 1 WHERE id = $1`, id)
			return err
		}
		if _, err := tx.ExecContext(ctx, `UPDATE account_codes SET consumed_at = NOW() WHERE id = $1`, id); err != nil {
			return fmt.Errorf("use code: %w", err)
		}
		ok = true
		return nil
	})
	return ok, err
}
