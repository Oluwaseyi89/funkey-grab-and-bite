package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/requestctx"
)

// These tests run against a real Postgres with the migrations applied, since
// the idempotency and rate limit SQL depends on Postgres behaviour a mock
// can't reproduce. Set TEST_DATABASE_URL to run them; they clean up after
// themselves.
func openIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	db, err := sql.Open("postgres", url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestIdempotencyRepositoryLifecycle(t *testing.T) {
	db := openIntegrationDB(t)
	ctx := context.Background()
	repo := NewIdempotencyRepository(db)
	scope := "POST /test/" + fmt.Sprint(time.Now().UnixNano())
	t.Cleanup(func() { db.Exec(`DELETE FROM idempotency_keys WHERE scope = $1`, scope) })

	expect := func(want models.IdempotencyOutcome, key, hash string, lease time.Duration) *models.IdempotencyRecord {
		t.Helper()
		got, record, err := repo.Begin(ctx, scope, key, hash, lease, time.Hour)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("key %s: expected outcome %d, got %d", key, want, got)
		}
		return record
	}

	expect(models.IdempotencyNew, "key-one-0001", "hash-a", time.Minute)
	expect(models.IdempotencyInProgress, "key-one-0001", "hash-a", time.Minute)
	expect(models.IdempotencyMismatch, "key-one-0001", "hash-b", time.Minute)

	// Completed inside a transaction, then the response saved.
	keyCtx := requestctx.WithIdempotency(ctx, requestctx.Idempotency{Scope: scope, Key: "key-one-0001"})
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.MarkCompleted(keyCtx, tx, "order", "17"); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	record := expect(models.IdempotencyReplay, "key-one-0001", "hash-a", time.Minute)
	if record.EntityType != "order" || record.EntityID != "17" || record.ResponseBody != nil {
		t.Fatalf("expected the committed entity and no response yet, got %+v", record)
	}
	if err := repo.SaveResponse(ctx, scope, "key-one-0001", 201, []byte(`{"id":17}`)); err != nil {
		t.Fatal(err)
	}
	record = expect(models.IdempotencyReplay, "key-one-0001", "hash-a", time.Minute)
	if record.ResponseStatus != 201 || string(record.ResponseBody) != `{"id":17}` {
		t.Fatalf("expected the saved response, got %+v", record)
	}
	if err := repo.Release(ctx, scope, "key-one-0001"); err != nil {
		t.Fatal(err)
	}
	expect(models.IdempotencyReplay, "key-one-0001", "hash-a", time.Minute)

	// A rolled-back completion leaves the key in progress.
	expect(models.IdempotencyNew, "key-two-0002", "hash-a", time.Minute)
	keyCtx = requestctx.WithIdempotency(ctx, requestctx.Idempotency{Scope: scope, Key: "key-two-0002"})
	tx, _ = db.Begin()
	_ = repo.MarkCompleted(keyCtx, tx, "order", "18")
	_ = tx.Rollback()
	expect(models.IdempotencyInProgress, "key-two-0002", "hash-a", time.Minute)

	// Once its lease runs out, a crashed attempt is taken over by a retry with
	// the same body, but never by a different request.
	expect(models.IdempotencyNew, "key-three-03", "hash-a", -time.Second)
	expect(models.IdempotencyMismatch, "key-three-03", "hash-b", time.Minute)
	expect(models.IdempotencyNew, "key-three-03", "hash-a", time.Minute)

	// A failed request frees its key.
	expect(models.IdempotencyNew, "key-four-004", "hash-a", time.Minute)
	if err := repo.Release(ctx, scope, "key-four-004"); err != nil {
		t.Fatal(err)
	}
	expect(models.IdempotencyNew, "key-four-004", "hash-b", time.Minute)
}

func TestRateLimitRepositoryCountsPerWindow(t *testing.T) {
	db := openIntegrationDB(t)
	ctx := context.Background()
	repo := NewRateLimitRepository(db)
	bucket := "test_" + fmt.Sprint(time.Now().UnixNano())
	t.Cleanup(func() { db.Exec(`DELETE FROM rate_limit_counters WHERE bucket = $1`, bucket) })

	window := time.Now().UTC().Truncate(time.Hour)
	for want := 1; want <= 3; want++ {
		got, err := repo.Hit(ctx, bucket, window)
		if err != nil || got != want {
			t.Fatalf("hit %d: got %d, %v", want, got, err)
		}
	}
	if got, _ := repo.Hit(ctx, bucket, window.Add(time.Hour)); got != 1 {
		t.Fatalf("a new window starts from 1, got %d", got)
	}

	removed, err := repo.DeleteBefore(ctx, window.Add(time.Minute))
	if err != nil || removed < 1 {
		t.Fatalf("expected the old window removed, got %d, %v", removed, err)
	}
}

func TestAuditRecordJoinsCallerTransaction(t *testing.T) {
	db := openIntegrationDB(t)
	ctx := context.Background()
	entityID := fmt.Sprint(time.Now().UnixNano())
	t.Cleanup(func() { db.Exec(`DELETE FROM audit_log WHERE entity_id = $1`, entityID) })

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if err := NewAuditRepository(db).Record(ctx, tx, models.AuditEntry{
		Action: "test.rolled_back", EntityType: "test", EntityID: entityID,
	}); err != nil {
		t.Fatal(err)
	}
	_ = tx.Rollback()

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM audit_log WHERE entity_id = $1`, entityID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("an entry written in a rolled-back transaction must not remain")
	}
}
