package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"funkey-grab-and-bite/funkey-bite-api/internal/handlers"
	"funkey-grab-and-bite/funkey-bite-api/internal/utils"
)

// RateLimitCounter counts requests per bucket per window. It is backed by
// Postgres so limits hold across restarts and every API instance.
type RateLimitCounter interface {
	Hit(ctx context.Context, bucket string, windowStart time.Time) (int, error)
}

// RateLimitKeyFunc returns who a request counts against, or false when the
// limit doesn't apply to it.
type RateLimitKeyFunc func(c *gin.Context) (string, bool)

// SharedRateLimit allows config.Limit requests per identity per window,
// counted with counter. Identities are hashed before they are stored, so the
// counters hold no phone numbers or IP addresses. If counting fails the request
// is let through: a database hiccup should not lock customers out.
func SharedRateLimit(counter RateLimitCounter, config RateLimitConfig, key RateLimitKeyFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		identity, ok := key(c)
		if !ok {
			c.Next()
			return
		}

		now := time.Now().UTC()
		windowStart := now.Truncate(config.Window)
		resetAt := windowStart.Add(config.Window)
		bucket := config.Prefix + hashIdentity(identity)

		count, err := counter.Hit(c.Request.Context(), bucket, windowStart)
		if err != nil {
			log.Printf("[rate limit] %s: %v; letting the request through", config.Prefix, err)
			c.Next()
			return
		}

		remaining := max(config.Limit-count, 0)
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.Limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetAt.Unix(), 10))

		if count > config.Limit {
			wait := resetAt.Sub(now)
			c.Header("Retry-After", strconv.Itoa(int(math.Ceil(wait.Seconds()))))
			handlers.Error(c, http.StatusTooManyRequests, "RATE_LIMITED", fmt.Sprintf(
				"Too many requests. Please try again in %d minutes.", int(math.Ceil(wait.Minutes())),
			))
			c.Abort()
			return
		}

		c.Next()
	}
}

func hashIdentity(identity string) string {
	sum := sha256.Sum256([]byte(identity))
	return hex.EncodeToString(sum[:16])
}

// ByClientIP counts requests against the client's IP address.
func ByClientIP(c *gin.Context) (string, bool) {
	ip := c.ClientIP()
	return "ip:" + ip, ip != ""
}

// phoneBodyFields are the JSON fields that carry a customer's phone number.
var phoneBodyFields = []string{"customerPhone", "contactPhone", "phone"}

// ByPhoneOrUser counts requests against the signed-in customer, or else the
// phone number the request is for (from a :phone path parameter or the JSON
// body), or else its email. Requests with neither, such as an emailed quote link, aren't counted.
func ByPhoneOrUser(c *gin.Context) (string, bool) {
	if userID, ok := c.Get("user_id"); ok {
		if isAdmin, _ := c.Get("is_admin"); isAdmin != true {
			return fmt.Sprintf("user:%v", userID), true
		}
	}

	if phone := utils.CanonicalPhone(c.Param("phone")); phone != "" {
		return "phone:" + phone, true
	}

	body, err := peekBody(c)
	if err != nil || len(body) == 0 {
		return "", false
	}
	var fields map[string]any
	if json.Unmarshal(body, &fields) != nil {
		return "", false
	}
	for _, name := range phoneBodyFields {
		if raw, ok := fields[name].(string); ok {
			if phone := utils.CanonicalPhone(raw); phone != "" {
				return "phone:" + phone, true
			}
		}
	}
	// Requests identified by email only, e.g. a password reset.
	if raw, ok := fields["email"].(string); ok && strings.Contains(raw, "@") {
		return "email:" + strings.ToLower(strings.TrimSpace(raw)), true
	}
	return "", false
}
