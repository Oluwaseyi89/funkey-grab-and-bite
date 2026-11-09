package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

// fakeStore records what the worker decided for each job.
type fakeStore struct {
	mu          sync.Mutex
	queue       []models.Job
	succeeded   []int64
	retried     map[int64]time.Time
	reasons     map[int64]string
	rescheduled map[int64]time.Time
	killed      map[int64]string
}

func newFakeStore(jobs ...models.Job) *fakeStore {
	return &fakeStore{queue: jobs, retried: map[int64]time.Time{}, reasons: map[int64]string{},
		rescheduled: map[int64]time.Time{}, killed: map[int64]string{}}
}

func (s *fakeStore) Claim(ctx context.Context, limit int, lease time.Duration) ([]models.Job, error) {
	claimed := s.queue
	s.queue = nil
	return claimed, nil
}
func (s *fakeStore) Succeed(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.succeeded = append(s.succeeded, id)
	return nil
}
func (s *fakeStore) Retry(ctx context.Context, id int64, runAt time.Time, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.retried[id], s.reasons[id] = runAt, reason
	return nil
}
func (s *fakeStore) Reschedule(ctx context.Context, id int64, runAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rescheduled[id] = runAt
	return nil
}
func (s *fakeStore) Kill(ctx context.Context, id int64, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.killed[id] = reason
	return nil
}

func job(id int64, kind string, attempts int) models.Job {
	return models.Job{ID: id, Kind: kind, Payload: json.RawMessage(`{"orderId":7}`), Attempts: attempts, MaxAttempts: 5}
}

func TestWorkerRecordsEachOutcome(t *testing.T) {
	later := time.Now().Add(10 * time.Minute)
	store := newFakeStore(
		job(1, "ok", 1),
		job(2, "flaky", 2),
		job(3, "flaky", 5),
		job(4, "missing", 1),
		job(5, "too-early", 1),
		job(6, "panics", 1),
		job(7, "unknown-kind", 1),
		job(8, "ok", 6),
	)
	worker := NewWorker(store)
	worker.Register("ok", func(ctx context.Context, payload json.RawMessage) error { return nil })
	worker.Register("flaky", func(ctx context.Context, payload json.RawMessage) error { return errors.New("smtp timeout") })
	worker.Register("missing", func(ctx context.Context, payload json.RawMessage) error {
		return Permanent(errors.New("order not found"))
	})
	worker.Register("too-early", func(ctx context.Context, payload json.RawMessage) error { return RunAgainAt(later) })
	worker.Register("panics", func(ctx context.Context, payload json.RawMessage) error { panic("nil map") })

	if ran := worker.RunDue(context.Background()); ran != 8 {
		t.Fatalf("expected 8 jobs run, got %d", ran)
	}

	if len(store.succeeded) != 1 || store.succeeded[0] != 1 {
		t.Fatalf("only job 1 should succeed, got %v", store.succeeded)
	}
	if at, ok := store.retried[2]; !ok || at.Before(time.Now().Add(7*time.Minute)) || store.reasons[2] != "smtp timeout" {
		t.Fatalf("a failure with attempts left must retry after backoff with its reason, got %v %q", at, store.reasons[2])
	}
	if _, ok := store.killed[3]; !ok {
		t.Fatal("a failure on the last attempt must mark the job dead")
	}
	if store.killed[4] != "order not found" {
		t.Fatalf("a permanent failure must mark the job dead at once, got %q", store.killed[4])
	}
	if !store.rescheduled[5].Equal(later) {
		t.Fatalf("RunAgainAt must reschedule without counting, got %v", store.rescheduled[5])
	}
	if _, ok := store.retried[6]; !ok {
		t.Fatal("a panicking handler must be retried, not crash the worker")
	}
	if _, ok := store.killed[7]; !ok {
		t.Fatal("a job with no handler must be marked dead")
	}
	if _, ok := store.killed[8]; !ok {
		t.Fatal("a job reclaimed past its attempts must be marked dead without running")
	}
}

func TestBackoffDoublesUpToCapWithJitter(t *testing.T) {
	for attempt, want := range map[int]time.Duration{1: 4 * time.Minute, 2: 8 * time.Minute, 3: 16 * time.Minute, 4: 30 * time.Minute, 9: 30 * time.Minute} {
		for range 20 {
			got := Backoff(attempt)
			if got < want-want/10 || got > want+want/10 {
				t.Fatalf("attempt %d: %v is outside %v ±10%%", attempt, got, want)
			}
		}
	}
}

func TestOrderPayload(t *testing.T) {
	if id, err := Order(json.RawMessage(`{"orderId":42}`)); err != nil || id != 42 {
		t.Fatalf("expected 42, got %d %v", id, err)
	}
	var permanent permanentError
	if _, err := Order(json.RawMessage(`{}`)); !errors.As(err, &permanent) {
		t.Fatalf("a payload without an order must be a permanent error, got %v", err)
	}
}

func TestWorkerStopsWhenContextIsCancelled(t *testing.T) {
	worker := NewWorker(newFakeStore())
	worker.pollInterval = 10 * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	worker.Start(ctx)
	worker.Notify()
	cancel()

	done := make(chan struct{})
	go func() { worker.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("the worker must stop when its context is cancelled")
	}
}

func TestWorkerRunsSlowJobsSideBySide(t *testing.T) {
	store := newFakeStore(job(1, "slow", 1), job(2, "slow", 1), job(3, "slow", 1), job(4, "slow", 1))
	worker := NewWorker(store)
	worker.Register("slow", func(ctx context.Context, payload json.RawMessage) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	start := time.Now()
	worker.RunDue(context.Background())
	if elapsed := time.Since(start); elapsed > 600*time.Millisecond {
		t.Fatalf("4 slow jobs took %v; they must not run one after another", elapsed)
	}
	if len(store.succeeded) != 4 {
		t.Fatalf("expected all 4 to succeed, got %d", len(store.succeeded))
	}
}
