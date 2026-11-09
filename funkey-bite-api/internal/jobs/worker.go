// Package jobs runs background work with bounded retries. Jobs live in
// Postgres (see repository.JobRepository); each API instance runs a Worker
// that claims due jobs, so no external queue is needed.
package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

// Job kinds.
const (
	KindOrderConfirmation = "order.confirmation"
	KindOrderInvoice      = "order.invoice"
	KindOrderPaymentCheck = "order.payment_check"
	// The payment-confirmed message, one job per channel so a failed email
	// is retried without sending the SMS again.
	KindOrderPaidEmail = "order.paid_email"
	KindOrderPaidSMS   = "order.paid_sms"
	// KindEmail sends an email rendered when it was queued, e.g. a catering
	// quote. Its payload can hold a private link, so it is queued as
	// sensitive and wiped once sent.
	KindEmail = "email.send"
)

// Store is where jobs are kept.
type Store interface {
	Claim(ctx context.Context, limit int, lease time.Duration) ([]models.Job, error)
	Succeed(ctx context.Context, id int64) error
	Retry(ctx context.Context, id int64, runAt time.Time, reason string) error
	Reschedule(ctx context.Context, id int64, runAt time.Time) error
	Kill(ctx context.Context, id int64, reason string) error
}

// Handler does one job. Returning an error retries it with backoff until its
// attempts run out; see Permanent and RunAgainAt for other outcomes. Jobs may
// run more than once (a retry after a timeout, or after an instance stopped
// mid-job), so handlers must be safe to repeat.
type Handler func(ctx context.Context, payload json.RawMessage) error

type permanentError struct{ err error }

func (e permanentError) Error() string { return e.err.Error() }
func (e permanentError) Unwrap() error { return e.err }

// Permanent marks an error that retrying won't fix, such as a missing order,
// so the job is given up at once.
func Permanent(err error) error { return permanentError{err} }

type runAgainError struct{ at time.Time }

func (e runAgainError) Error() string { return "run again at " + e.at.Format(time.RFC3339) }

// RunAgainAt asks for the job to run again at t without counting this
// attempt, for a job that found it was too early to act.
func RunAgainAt(t time.Time) error { return runAgainError{t} }

const (
	defaultPollInterval = 3 * time.Second
	defaultBatchSize    = 10
	// A slow provider (e.g. an email service retrying) shouldn't hold up the
	// rest of a batch, so up to this many jobs run at once.
	maxConcurrentJobs = 4
	jobTimeout        = time.Minute
	// The lease outlasts the job timeout, so a job is only reclaimed when its
	// worker stopped.
	jobLease = 2 * jobTimeout

	backoffBase = 4 * time.Minute
	backoffCap  = 30 * time.Minute
)

// Worker claims and runs due jobs.
type Worker struct {
	store        Store
	handlers     map[string]Handler
	pollInterval time.Duration
	batchSize    int
	now          func() time.Time
	wake         chan struct{}
	wg           sync.WaitGroup
}

func NewWorker(store Store) *Worker {
	return &Worker{
		store:        store,
		handlers:     map[string]Handler{},
		pollInterval: defaultPollInterval,
		batchSize:    defaultBatchSize,
		now:          time.Now,
		wake:         make(chan struct{}, 1),
	}
}

// Register sets the handler for kind. Call it before Start.
func (w *Worker) Register(kind string, handler Handler) {
	w.handlers[kind] = handler
}

// Notify asks the worker to look for jobs now instead of at its next poll,
// e.g. after queueing an email.
func (w *Worker) Notify() {
	select {
	case w.wake <- struct{}{}:
	default:
	}
}

// Start runs the worker until ctx is cancelled; Wait blocks until it stops.
func (w *Worker) Start(ctx context.Context) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		ticker := time.NewTicker(w.pollInterval)
		defer ticker.Stop()
		for {
			w.RunDue(ctx)
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			case <-w.wake:
			}
		}
	}()
}

// Wait blocks until a started worker has stopped.
func (w *Worker) Wait() { w.wg.Wait() }

// RunDue claims and runs one batch of due jobs, returning how many ran.
func (w *Worker) RunDue(ctx context.Context) int {
	claimed, err := w.store.Claim(ctx, w.batchSize, jobLease)
	if err != nil {
		if ctx.Err() == nil {
			log.Printf("[jobs] %v", err)
		}
		return 0
	}
	slots := make(chan struct{}, maxConcurrentJobs)
	var batch sync.WaitGroup
	for _, job := range claimed {
		slots <- struct{}{}
		batch.Add(1)
		go func() {
			defer batch.Done()
			defer func() { <-slots }()
			w.run(ctx, job)
		}()
	}
	batch.Wait()
	return len(claimed)
}

func (w *Worker) run(ctx context.Context, job models.Job) {
	// Finishing must not be skipped because the app is shutting down.
	done := context.WithoutCancel(ctx)

	handler, ok := w.handlers[job.Kind]
	if !ok {
		w.record(w.store.Kill(done, job.ID, "no handler for job kind "+job.Kind))
		return
	}
	if job.Attempts > job.MaxAttempts {
		// Reclaimed after its worker stopped on the last attempt.
		w.record(w.store.Kill(done, job.ID, "attempts exhausted"))
		return
	}

	err := w.call(ctx, handler, job)

	var runAgain runAgainError
	var permanent permanentError
	switch {
	case err == nil:
		w.record(w.store.Succeed(done, job.ID))
	case errors.As(err, &runAgain):
		w.record(w.store.Reschedule(done, job.ID, runAgain.at))
	case errors.As(err, &permanent) || job.Attempts >= job.MaxAttempts:
		log.Printf("[jobs] %s job %d failed for good after %d attempt(s): %v", job.Kind, job.ID, job.Attempts, err)
		w.record(w.store.Kill(done, job.ID, err.Error()))
	default:
		next := w.now().Add(Backoff(job.Attempts))
		log.Printf("[jobs] %s job %d attempt %d failed, retrying at %s: %v", job.Kind, job.ID, job.Attempts, next.Format(time.RFC3339), err)
		w.record(w.store.Retry(done, job.ID, next, err.Error()))
	}
}

func (w *Worker) call(ctx context.Context, handler Handler, job models.Job) (err error) {
	ctx, cancel := context.WithTimeout(ctx, jobTimeout)
	defer cancel()
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("panic: %v", p)
		}
	}()
	return handler(ctx, job.Payload)
}

func (w *Worker) record(err error) {
	if err != nil {
		log.Printf("[jobs] %v", err)
	}
}

// Backoff is the wait before retrying after attempt failed: it doubles from
// four minutes up to half an hour, with jitter so jobs that failed together
// don't retry together. With five attempts a job is retried for about an hour.
func Backoff(attempt int) time.Duration {
	wait := backoffBase
	for i := 1; i < attempt && wait < backoffCap; i++ {
		wait *= 2
	}
	wait = min(wait, backoffCap)
	jitter := time.Duration(rand.Int63n(int64(wait / 5)))
	return wait - wait/10 + jitter
}

// Order returns the order ID in an order job's payload.
func Order(payload json.RawMessage) (int, error) {
	var p models.OrderJobPayload
	if err := json.Unmarshal(payload, &p); err != nil || p.OrderID <= 0 {
		return 0, Permanent(fmt.Errorf("invalid order job payload %s", payload))
	}
	return p.OrderID, nil
}
