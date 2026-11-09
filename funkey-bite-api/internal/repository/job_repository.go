package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

const defaultJobMaxAttempts = 5

// JobRepository stores background jobs in Postgres.
type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

// Enqueue queues job using exec: pass the transaction making the change that
// needs the job, so the job is queued only if the change is saved (nil queues
// it on its own). A job whose DedupeKey is already queued is skipped.
func (r *JobRepository) Enqueue(ctx context.Context, exec Execer, job models.NewJob) error {
	if exec == nil {
		exec = r.db
	}
	payload, err := json.Marshal(job.Payload)
	if err != nil {
		return fmt.Errorf("encode %s job payload: %w", job.Kind, err)
	}
	maxAttempts := job.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = defaultJobMaxAttempts
	}
	var runAt any
	if !job.RunAt.IsZero() {
		runAt = job.RunAt
	}

	_, err = exec.ExecContext(ctx, `
		INSERT INTO jobs (kind, payload, max_attempts, run_at, dedupe_key, sensitive)
		VALUES ($1, $2, $3, COALESCE($4, NOW()), $5, $6)
		ON CONFLICT (dedupe_key) DO NOTHING
	`, job.Kind, string(payload), maxAttempts, runAt, nullIfEmpty(job.DedupeKey), job.Sensitive)
	if err != nil {
		return fmt.Errorf("queue %s job: %w", job.Kind, err)
	}
	return nil
}

// Claim takes up to limit due jobs for this worker for lease, counting an
// attempt for each. Jobs another worker holds are skipped, and a running job
// whose lease ran out (its worker stopped) is claimed again.
func (r *JobRepository) Claim(ctx context.Context, limit int, lease time.Duration) ([]models.Job, error) {
	rows, err := r.db.QueryContext(ctx, `
		UPDATE jobs
		SET status = 'running', attempts = attempts + 1,
		    locked_until = NOW() + make_interval(secs => $2), updated_at = NOW()
		WHERE id IN (
			SELECT id FROM jobs
			WHERE (status = 'pending' AND run_at <= NOW())
			   OR (status = 'running' AND locked_until < NOW())
			ORDER BY run_at
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, kind, payload, attempts, max_attempts
	`, limit, lease.Seconds())
	if err != nil {
		return nil, fmt.Errorf("claim jobs: %w", err)
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		var payload []byte
		if err := rows.Scan(&job.ID, &job.Kind, &payload, &job.Attempts, &job.MaxAttempts); err != nil {
			return nil, fmt.Errorf("read claimed job: %w", err)
		}
		job.Payload = payload
		jobs = append(jobs, job)
	}
	return jobs, rows.Err()
}

// Succeed marks a job done, wiping a sensitive job's payload.
func (r *JobRepository) Succeed(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE jobs
		SET status = 'succeeded', locked_until = NULL, last_error = NULL,
		    payload = CASE WHEN sensitive THEN '{}'::jsonb ELSE payload END,
		    completed_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("complete job %d: %w", id, err)
	}
	return nil
}

// Retry records a failed attempt and runs the job again at runAt.
func (r *JobRepository) Retry(ctx context.Context, id int64, runAt time.Time, reason string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE jobs
		SET status = 'pending', run_at = $2, locked_until = NULL, last_error = $3, updated_at = NOW()
		WHERE id = $1
	`, id, runAt, reason)
	if err != nil {
		return fmt.Errorf("retry job %d: %w", id, err)
	}
	return nil
}

// Reschedule runs the job again at runAt without counting this attempt, for
// a job that found it was too early to act.
func (r *JobRepository) Reschedule(ctx context.Context, id int64, runAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE jobs
		SET status = 'pending', run_at = $2, locked_until = NULL,
		    attempts = GREATEST(attempts - 1, 0), updated_at = NOW()
		WHERE id = $1
	`, id, runAt)
	if err != nil {
		return fmt.Errorf("reschedule job %d: %w", id, err)
	}
	return nil
}

// Kill gives up on a job; it stays visible as dead until retried by hand.
func (r *JobRepository) Kill(ctx context.Context, id int64, reason string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE jobs
		SET status = 'dead', locked_until = NULL, last_error = $2, updated_at = NOW()
		WHERE id = $1
	`, id, reason)
	if err != nil {
		return fmt.Errorf("mark job %d dead: %w", id, err)
	}
	return nil
}

// DeleteSucceededBefore removes finished jobs completed before cutoff and
// returns how many.
func (r *JobRepository) DeleteSucceededBefore(ctx context.Context, cutoff time.Time) (int64, error) {
	result, err := r.db.ExecContext(ctx, `DELETE FROM jobs WHERE status = 'succeeded' AND completed_at < $1`, cutoff)
	if err != nil {
		return 0, fmt.Errorf("delete finished jobs: %w", err)
	}
	return result.RowsAffected()
}

// List returns jobs with status (all when empty), newest first, and how many
// there are. Payloads are left out.
func (r *JobRepository) List(ctx context.Context, status string, limit, offset int) ([]models.JobSummary, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM jobs WHERE ($1 = '' OR status = $1)`, status).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count jobs: %w", err)
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, kind, status, attempts, max_attempts, run_at, last_error, created_at, updated_at, completed_at
		FROM jobs
		WHERE ($1 = '' OR status = $1)
		ORDER BY updated_at DESC, id DESC
		LIMIT $2 OFFSET $3
	`, status, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list jobs: %w", err)
	}
	defer rows.Close()

	jobs := []models.JobSummary{}
	for rows.Next() {
		var job models.JobSummary
		var lastError sql.NullString
		var completedAt sql.NullTime
		if err := rows.Scan(&job.ID, &job.Kind, &job.Status, &job.Attempts, &job.MaxAttempts,
			&job.RunAt, &lastError, &job.CreatedAt, &job.UpdatedAt, &completedAt); err != nil {
			return nil, 0, fmt.Errorf("read job: %w", err)
		}
		if lastError.Valid {
			job.LastError = &lastError.String
		}
		if completedAt.Valid {
			job.CompletedAt = &completedAt.Time
		}
		jobs = append(jobs, job)
	}
	return jobs, total, rows.Err()
}

// RetryDead gives a dead job a fresh set of attempts, starting now. It
// returns false when the job doesn't exist or isn't dead.
func (r *JobRepository) RetryDead(ctx context.Context, id int64) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE jobs
		SET status = 'pending', attempts = 0, run_at = NOW(), locked_until = NULL, updated_at = NOW()
		WHERE id = $1 AND status = 'dead'
	`, id)
	if err != nil {
		return false, fmt.Errorf("retry job %d: %w", id, err)
	}
	n, err := result.RowsAffected()
	return n > 0, err
}
