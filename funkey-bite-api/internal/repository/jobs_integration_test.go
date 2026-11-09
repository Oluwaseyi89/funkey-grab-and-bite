package repository

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

// Runs only with TEST_DATABASE_URL; see openIntegrationDB.
func TestJobRepositoryQueueAndClaim(t *testing.T) {
	db := openIntegrationDB(t)
	ctx := context.Background()
	repo := NewJobRepository(db)
	kind := "test." + fmt.Sprint(time.Now().UnixNano())
	t.Cleanup(func() { db.Exec(`DELETE FROM jobs WHERE kind = $1`, kind) })

	count := func(where string, args ...any) int {
		t.Helper()
		var n int
		if err := db.QueryRow(`SELECT COUNT(*) FROM jobs WHERE kind = $1 AND `+where, append([]any{kind}, args...)...).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}

	// A duplicate dedupe key is skipped; a job queued in a rolled-back
	// transaction never exists.
	for range 2 {
		if err := repo.Enqueue(ctx, nil, models.NewJob{Kind: kind, Payload: models.OrderJobPayload{OrderID: 1}, DedupeKey: kind + ":1"}); err != nil {
			t.Fatal(err)
		}
	}
	tx, _ := db.Begin()
	_ = repo.Enqueue(ctx, tx, models.NewJob{Kind: kind, Payload: models.OrderJobPayload{OrderID: 2}})
	_ = tx.Rollback()
	if n := count("true"); n != 1 {
		t.Fatalf("expected 1 job after a duplicate and a rollback, got %d", n)
	}

	// Two workers claiming at once never get the same job; a future job waits.
	for i := 3; i <= 6; i++ {
		_ = repo.Enqueue(ctx, nil, models.NewJob{Kind: kind, Payload: models.OrderJobPayload{OrderID: i}})
	}
	_ = repo.Enqueue(ctx, nil, models.NewJob{Kind: kind, Payload: models.OrderJobPayload{OrderID: 99}, RunAt: time.Now().Add(time.Hour)})

	var mu sync.Mutex
	seen := map[int64]int{}
	var wg sync.WaitGroup
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			claimed, err := repo.Claim(ctx, 5, time.Minute)
			if err != nil {
				t.Error(err)
				return
			}
			mu.Lock()
			defer mu.Unlock()
			for _, job := range claimed {
				if job.Kind == kind {
					seen[job.ID]++
				}
			}
		}()
	}
	wg.Wait()
	if len(seen) != 5 {
		t.Fatalf("expected the 5 due jobs claimed, got %d", len(seen))
	}
	for id, n := range seen {
		if n != 1 {
			t.Fatalf("job %d was claimed %d times", id, n)
		}
	}
	if n := count("status = 'pending'"); n != 1 {
		t.Fatalf("the future job must still be waiting, got %d pending", n)
	}

	// Outcomes.
	ids := make([]int64, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	_ = repo.Succeed(ctx, ids[0])
	_ = repo.Retry(ctx, ids[1], time.Now().Add(-time.Second), "smtp timeout")
	_ = repo.Reschedule(ctx, ids[2], time.Now().Add(time.Hour))
	_ = repo.Kill(ctx, ids[3], "order not found")
	if count("status = 'succeeded' AND completed_at IS NOT NULL") != 1 ||
		count("status = 'dead' AND last_error = 'order not found'") != 1 ||
		count("id = $2 AND status = 'pending' AND attempts = 0", ids[2]) != 1 {
		t.Fatal("expected succeeded, dead and an uncounted reschedule to be recorded")
	}

	// The retried job is due again, with its attempt counted.
	claimed, _ := repo.Claim(ctx, 10, -time.Second) // lease already over
	var retried *models.Job
	for i := range claimed {
		if claimed[i].ID == ids[1] {
			retried = &claimed[i]
		}
	}
	if retried == nil || retried.Attempts != 2 {
		t.Fatalf("expected the retried job claimed on its second attempt, got %+v", retried)
	}

	// A running job whose lease ran out (its worker stopped) is claimed again.
	again, _ := repo.Claim(ctx, 10, time.Minute)
	found := false
	for _, job := range again {
		if job.ID == ids[1] && job.Attempts == 3 {
			found = true
		}
	}
	if !found {
		t.Fatal("an abandoned running job must be claimable again")
	}

	removed, err := repo.DeleteSucceededBefore(ctx, time.Now().Add(time.Minute))
	if err != nil || removed < 1 {
		t.Fatalf("expected the finished job removed, got %d %v", removed, err)
	}
}
