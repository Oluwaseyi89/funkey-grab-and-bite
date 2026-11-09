package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
	"funkey-grab-and-bite/funkey-bite-api/internal/handlers"
	"funkey-grab-and-bite/funkey-bite-api/internal/requestctx"
)

// IdempotencyKeyHeader is the header clients send with a unique value per
// action (e.g. one checkout), reusing it when they retry that action.
const IdempotencyKeyHeader = "Idempotency-Key"

const (
	defaultIdempotencyLease = 30 * time.Second
	defaultIdempotencyTTL   = 24 * time.Hour
)

var idempotencyKeyPattern = regexp.MustCompile(`^[A-Za-z0-9_.:-]{8,255}$`)

// IdempotencyStore saves the result of each request made with a key.
type IdempotencyStore interface {
	Begin(ctx context.Context, scope, key, requestHash string, lease, ttl time.Duration) (models.IdempotencyOutcome, *models.IdempotencyRecord, error)
	SaveResponse(ctx context.Context, scope, key string, status int, body []byte) error
	Release(ctx context.Context, scope, key string) error
}

// IdempotencyReplayer answers a retry of a request that completed but whose
// response was never saved (the API stopped right after committing), from
// the entity the request created.
type IdempotencyReplayer func(c *gin.Context, record *models.IdempotencyRecord)

type IdempotencyOptions struct {
	// Required rejects requests without a key; otherwise they run as normal.
	Required bool
	// Lease is how long a request may run before a retry with the same key
	// may assume it crashed and take over. Defaults to 30 seconds.
	Lease time.Duration
	// TTL is how long results are kept for retries. Defaults to 24 hours.
	TTL    time.Duration
	Replay IdempotencyReplayer
}

// Idempotency makes a write safe to retry: the first request with a given
// Idempotency-Key runs and its successful response is stored; any retry with
// that key gets the stored response instead of repeating the change. A retry
// while the first is still running gets 409, and reusing a key for a
// different request gets 422. Failed requests free their key so the client
// can try again with it.
//
// Handlers that change data should record completion inside their
// transaction (IdempotencyRepository.MarkCompleted with RequestContext), so a
// key is only ever completed together with the change it made.
func Idempotency(store IdempotencyStore, opts IdempotencyOptions) gin.HandlerFunc {
	if opts.Lease <= 0 {
		opts.Lease = defaultIdempotencyLease
	}
	if opts.TTL <= 0 {
		opts.TTL = defaultIdempotencyTTL
	}

	return func(c *gin.Context) {
		key := c.GetHeader(IdempotencyKeyHeader)
		if key == "" {
			if opts.Required {
				handlers.Error(c, http.StatusBadRequest, "IDEMPOTENCY_KEY_REQUIRED", "This request needs an Idempotency-Key header.")
				c.Abort()
				return
			}
			c.Next()
			return
		}
		if !idempotencyKeyPattern.MatchString(key) {
			handlers.Error(c, http.StatusBadRequest, "IDEMPOTENCY_KEY_INVALID", "The Idempotency-Key must be 8 to 255 letters, digits or - _ . :")
			c.Abort()
			return
		}

		body, err := peekBody(c)
		if err != nil {
			status := http.StatusBadRequest
			if errors.Is(err, errBodyTooLarge) {
				status = http.StatusRequestEntityTooLarge
			}
			handlers.Error(c, status, "INVALID_BODY", "Failed to read the request body.")
			c.Abort()
			return
		}

		scope := c.Request.Method + " " + routeOf(c)
		requestHash := hashRequest(c.Request.URL.RequestURI(), body)

		outcome, record, err := store.Begin(c.Request.Context(), scope, key, requestHash, opts.Lease, opts.TTL)
		if err != nil {
			log.Printf("[idempotency] %s: %v", scope, err)
			handlers.Error(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Please try again in a moment.")
			c.Abort()
			return
		}

		switch outcome {
		case models.IdempotencyMismatch:
			handlers.Error(c, http.StatusUnprocessableEntity, "IDEMPOTENCY_KEY_REUSED", "This Idempotency-Key was already used for a different request.")
			c.Abort()
			return
		case models.IdempotencyInProgress:
			c.Header("Retry-After", "1")
			handlers.Error(c, http.StatusConflict, "REQUEST_IN_PROGRESS", "This request is still being processed. Please wait a moment.")
			c.Abort()
			return
		case models.IdempotencyReplay:
			replay(c, record, opts.Replay)
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(requestctx.WithIdempotency(
			c.Request.Context(), requestctx.Idempotency{Scope: scope, Key: key},
		))
		capture := &idempotencyCaptureWriter{ResponseWriter: c.Writer}
		c.Writer = capture

		c.Next()

		// The request may already be cancelled; storing the result must not be.
		ctx := context.WithoutCancel(c.Request.Context())
		status := capture.Status()
		if status >= 200 && status < 300 {
			if err := store.SaveResponse(ctx, scope, key, status, capture.body.Bytes()); err != nil {
				log.Printf("[idempotency] save response for %s: %v", scope, err)
			}
			return
		}
		if err := store.Release(ctx, scope, key); err != nil {
			log.Printf("[idempotency] release key for %s: %v", scope, err)
		}
	}
}

func replay(c *gin.Context, record *models.IdempotencyRecord, replayer IdempotencyReplayer) {
	c.Header("Idempotent-Replayed", "true")
	switch {
	case len(record.ResponseBody) > 0:
		c.Data(record.ResponseStatus, "application/json; charset=utf-8", record.ResponseBody)
	case replayer != nil && record.EntityID != "":
		replayer(c, record)
	default:
		handlers.Error(c, http.StatusConflict, "REQUEST_ALREADY_PROCESSED", "This request was already processed.")
	}
}

func routeOf(c *gin.Context) string {
	if route := c.FullPath(); route != "" {
		return route
	}
	return c.Request.URL.Path
}

// hashRequest fingerprints the exact request: its path and query (which may
// carry IDs) and its body.
func hashRequest(uri string, body []byte) string {
	h := sha256.New()
	h.Write([]byte(uri))
	h.Write([]byte{'\n'})
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}

// idempotencyCaptureWriter keeps a copy of the handler's response body while
// passing it through.
type idempotencyCaptureWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *idempotencyCaptureWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w *idempotencyCaptureWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
