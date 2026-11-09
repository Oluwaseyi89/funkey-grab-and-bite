package middleware

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type fakeCounter struct {
	mu     sync.Mutex
	counts map[string]int
	err    error
}

func (f *fakeCounter) Hit(ctx context.Context, bucket string, windowStart time.Time) (int, error) {
	if f.err != nil {
		return 0, f.err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.counts == nil {
		f.counts = map[string]int{}
	}
	key := bucket + windowStart.String()
	f.counts[key]++
	return f.counts[key], nil
}

func limitedRouter(handler ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/orders", handler...)
	router.GET("/track/:phone", handler...)
	return router
}

func postJSON(router http.Handler, body string, ip string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = ip + ":1234"
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}

func ok(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) }

func TestIdentityLimitCountsSamePhoneAcrossFormatsAndIPs(t *testing.T) {
	router := limitedRouter(IdentityRateLimitMiddleware(&fakeCounter{}, "orders"), ok)

	for i := range IdentityRateLimit.Limit {
		// Alternate phone formats and IPs: they must all count as one customer.
		phone := "08031112222"
		if i%2 == 1 {
			phone = "+234 803 111 2222"
		}
		res := postJSON(router, `{"customerPhone":"`+phone+`"}`, "198.51.100."+string(rune('1'+i%5)))
		if res.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i+1, res.Code)
		}
	}

	res := postJSON(router, `{"customerPhone":"0803 111 2222"}`, "203.0.113.9")
	if res.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 once the phone's limit is used, got %d", res.Code)
	}
	if res.Header().Get("Retry-After") == "" || !strings.Contains(res.Body.String(), "RATE_LIMITED") {
		t.Fatalf("expected Retry-After and RATE_LIMITED, got headers %v body %s", res.Header(), res.Body.String())
	}

	if res := postJSON(router, `{"customerPhone":"08035556666"}`, "203.0.113.9"); res.Code != http.StatusOK {
		t.Fatalf("another phone must have its own allowance, got %d", res.Code)
	}
}

func TestIdentityLimitReadsPhoneFromPathAndLeavesBodyForHandler(t *testing.T) {
	counter := &fakeCounter{}
	var seen string
	router := limitedRouter(IdentityRateLimitMiddleware(counter, "lookup"), func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		seen = string(body)
		c.JSON(http.StatusOK, gin.H{})
	})

	postJSON(router, `{"contactPhone":"08031112222","eventName":"Wedding"}`, "198.51.100.1")
	if seen != `{"contactPhone":"08031112222","eventName":"Wedding"}` {
		t.Fatalf("handler must still receive the full body, got %q", seen)
	}

	req := httptest.NewRequest(http.MethodGet, "/track/08031112222", nil)
	router.ServeHTTP(httptest.NewRecorder(), req)
	if len(counter.counts) != 1 {
		t.Fatalf("body and path phone must share one bucket, got %d buckets", len(counter.counts))
	}
}

func TestIdentityLimitUsesSignedInCustomerAndSkipsRequestsWithoutIdentity(t *testing.T) {
	counter := &fakeCounter{}
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/orders", func(c *gin.Context) { c.Set("user_id", 42); c.Next() },
		IdentityRateLimitMiddleware(counter, "orders"), ok)
	router.POST("/anon", IdentityRateLimitMiddleware(counter, "orders"), ok)

	for range IdentityRateLimit.Limit {
		postJSON(router, `{"customerPhone":"08031112222"}`, "198.51.100.1")
	}
	if res := postJSON(router, `{"customerPhone":"08099999999"}`, "198.51.100.1"); res.Code != http.StatusTooManyRequests {
		t.Fatalf("a signed-in customer is limited however the phone changes, got %d", res.Code)
	}

	req := httptest.NewRequest(http.MethodPost, "/anon", strings.NewReader(`{"token":"abc"}`))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("a request with no phone or user is not counted, got %d", res.Code)
	}
}

func TestOrderRateLimitCapsPerIP(t *testing.T) {
	router := limitedRouter(OrderRateLimitMiddleware(&fakeCounter{}), ok)

	for i := range OrderCreationRateLimit.Limit {
		if res := postJSON(router, `{}`, "198.51.100.7"); res.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i+1, res.Code)
		}
	}
	if res := postJSON(router, `{}`, "198.51.100.7"); res.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 after %d orders from one IP, got %d", OrderCreationRateLimit.Limit, res.Code)
	}
	if res := postJSON(router, `{}`, "198.51.100.8"); res.Code != http.StatusOK {
		t.Fatalf("another IP must not be affected, got %d", res.Code)
	}
}

func TestSharedRateLimitLetsRequestsThroughWhenCountingFails(t *testing.T) {
	router := limitedRouter(OrderRateLimitMiddleware(&fakeCounter{err: errors.New("db down")}), ok)
	if res := postJSON(router, `{}`, "198.51.100.7"); res.Code != http.StatusOK {
		t.Fatalf("a counting failure must not block customers, got %d", res.Code)
	}
}
