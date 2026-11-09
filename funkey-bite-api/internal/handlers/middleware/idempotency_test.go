package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/requestctx"
)

// memoryIdempotencyStore follows the same rules as IdempotencyRepository.
type memoryIdempotencyStore struct {
	mu      sync.Mutex
	records map[string]*models.IdempotencyRecord
}

func (s *memoryIdempotencyStore) Begin(ctx context.Context, scope, key, hash string, lease, ttl time.Duration) (models.IdempotencyOutcome, *models.IdempotencyRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.records == nil {
		s.records = map[string]*models.IdempotencyRecord{}
	}
	record, ok := s.records[scope+key]
	switch {
	case !ok:
		s.records[scope+key] = &models.IdempotencyRecord{Scope: scope, Key: key, RequestHash: hash, Status: "in_progress"}
		return models.IdempotencyNew, nil, nil
	case record.RequestHash != hash:
		return models.IdempotencyMismatch, record, nil
	case record.Status == "completed":
		return models.IdempotencyReplay, record, nil
	default:
		return models.IdempotencyInProgress, record, nil
	}
}

func (s *memoryIdempotencyStore) SaveResponse(ctx context.Context, scope, key string, status int, body []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	record := s.records[scope+key]
	record.Status, record.ResponseStatus, record.ResponseBody = "completed", status, append([]byte(nil), body...)
	return nil
}

func (s *memoryIdempotencyStore) Release(ctx context.Context, scope, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if record, ok := s.records[scope+key]; ok && record.Status == "in_progress" {
		delete(s.records, scope+key)
	}
	return nil
}

func idempotentRouter(store IdempotencyStore, opts IdempotencyOptions, handler gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/orders", Idempotency(store, opts), handler)
	return router
}

func send(router http.Handler, key, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if key != "" {
		req.Header.Set(IdempotencyKeyHeader, key)
	}
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}

func TestIdempotencyReplaysCompletedRequestWithoutRunningItAgain(t *testing.T) {
	calls := 0
	router := idempotentRouter(&memoryIdempotencyStore{}, IdempotencyOptions{}, func(c *gin.Context) {
		calls++
		c.JSON(http.StatusCreated, gin.H{"orderNumber": "FG-1", "call": calls})
	})

	first := send(router, "checkout-abc-123", `{"items":[1]}`)
	retry := send(router, "checkout-abc-123", `{"items":[1]}`)

	if calls != 1 {
		t.Fatalf("the handler must run once, ran %d times", calls)
	}
	if retry.Code != http.StatusCreated || retry.Body.String() != first.Body.String() {
		t.Fatalf("retry must get the first response: %d %s vs %s", retry.Code, retry.Body.String(), first.Body.String())
	}
	if retry.Header().Get("Idempotent-Replayed") != "true" {
		t.Fatal("a replayed response must say so")
	}
}

func TestIdempotencyRejectsKeyReusedForDifferentRequest(t *testing.T) {
	router := idempotentRouter(&memoryIdempotencyStore{}, IdempotencyOptions{}, func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{})
	})

	send(router, "checkout-abc-123", `{"items":[1]}`)
	res := send(router, "checkout-abc-123", `{"items":[2]}`)
	if res.Code != http.StatusUnprocessableEntity || !strings.Contains(res.Body.String(), "IDEMPOTENCY_KEY_REUSED") {
		t.Fatalf("expected 422 IDEMPOTENCY_KEY_REUSED, got %d %s", res.Code, res.Body.String())
	}
}

func TestIdempotencyAnswersInProgressWhileFirstRequestRuns(t *testing.T) {
	store := &memoryIdempotencyStore{}
	_, _, _ = store.Begin(context.Background(), "POST /orders", "checkout-abc-123", hashRequest("/orders", []byte(`{}`)), time.Minute, time.Hour)
	router := idempotentRouter(store, IdempotencyOptions{}, func(c *gin.Context) {
		t.Fatal("the handler must not run while the first attempt is in progress")
	})

	res := send(router, "checkout-abc-123", `{}`)
	if res.Code != http.StatusConflict || res.Header().Get("Retry-After") == "" {
		t.Fatalf("expected 409 with Retry-After, got %d %v", res.Code, res.Header())
	}
}

func TestIdempotencyFreesKeyWhenRequestFails(t *testing.T) {
	calls := 0
	router := idempotentRouter(&memoryIdempotencyStore{}, IdempotencyOptions{}, func(c *gin.Context) {
		calls++
		if calls == 1 {
			c.JSON(http.StatusConflict, gin.H{"error": "Burger is out of stock"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{})
	})

	send(router, "checkout-abc-123", `{}`)
	res := send(router, "checkout-abc-123", `{}`)
	if calls != 2 || res.Code != http.StatusCreated {
		t.Fatalf("a failed request must free its key for a retry, got %d calls and %d", calls, res.Code)
	}
}

func TestIdempotencyPassesKeyToHandlerContext(t *testing.T) {
	var got requestctx.Idempotency
	router := idempotentRouter(&memoryIdempotencyStore{}, IdempotencyOptions{}, func(c *gin.Context) {
		got, _ = requestctx.IdempotencyFrom(RequestContext(c))
		c.JSON(http.StatusCreated, gin.H{})
	})

	send(router, "checkout-abc-123", `{}`)
	if got.Key != "checkout-abc-123" || got.Scope != "POST /orders" {
		t.Fatalf("expected key and scope in the context, got %+v", got)
	}
}

func TestIdempotencyKeyRules(t *testing.T) {
	store := &memoryIdempotencyStore{}
	handler := func(c *gin.Context) { c.JSON(http.StatusCreated, gin.H{}) }

	if res := send(idempotentRouter(store, IdempotencyOptions{}, handler), "", `{}`); res.Code != http.StatusCreated {
		t.Fatalf("without Required, a request without a key runs as normal, got %d", res.Code)
	}
	if res := send(idempotentRouter(store, IdempotencyOptions{Required: true}, handler), "", `{}`); res.Code != http.StatusBadRequest {
		t.Fatalf("with Required, a missing key is rejected, got %d", res.Code)
	}
	if res := send(idempotentRouter(store, IdempotencyOptions{}, handler), "short", `{}`); res.Code != http.StatusBadRequest {
		t.Fatalf("a key under 8 characters is rejected, got %d", res.Code)
	}
}

func TestIdempotencyReplaysFromEntityWhenResponseWasNotSaved(t *testing.T) {
	store := &memoryIdempotencyStore{records: map[string]*models.IdempotencyRecord{
		"POST /orderscheckout-abc-123": {
			RequestHash: hashRequest("/orders", []byte(`{}`)), Status: "completed",
			EntityType: "order", EntityID: "17",
		},
	}}
	router := idempotentRouter(store, IdempotencyOptions{Replay: func(c *gin.Context, record *models.IdempotencyRecord) {
		c.JSON(http.StatusCreated, gin.H{"id": record.EntityID})
	}}, func(c *gin.Context) { t.Fatal("a completed request must not run again") })

	res := send(router, "checkout-abc-123", `{}`)
	if res.Code != http.StatusCreated || !strings.Contains(res.Body.String(), `"17"`) {
		t.Fatalf("expected the order re-rendered from its ID, got %d %s", res.Code, res.Body.String())
	}
}
