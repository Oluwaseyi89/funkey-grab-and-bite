package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple sliding window rate limiter
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int           // Maximum requests
	window   time.Duration // Time window
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow checks if a request is allowed based on rate limiting
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Clean old requests
	if timestamps, exists := rl.requests[key]; exists {
		var validTimestamps []time.Time
		for _, ts := range timestamps {
			if ts.After(windowStart) {
				validTimestamps = append(validTimestamps, ts)
			}
		}
		rl.requests[key] = validTimestamps
	}

	// Check if under limit
	if len(rl.requests[key]) >= rl.limit {
		return false
	}

	// Record this request
	rl.requests[key] = append(rl.requests[key], now)
	return true
}

// RateLimitMiddleware creates a Gin middleware for rate limiting
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		// Use IP address as the key (you could also use user ID for authenticated users)
		key := c.ClientIP()

		// For specific paths, you might want different limits
		path := c.Request.URL.Path

		// More strict limits for public tracking endpoints
		if path == "/api/v1/orders/track" || path == "/api/v1/orders/track/:phone/:orderNumber" {
			key = "track_" + key // Separate bucket for tracking
		}

		if !limiter.Allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":      "Rate limit exceeded. Please try again later.",
				"retryAfter": window.Seconds(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// IPBasedRateLimit is a simpler middleware using just IP address
func IPBasedRateLimit(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		if !limiter.Allow(c.ClientIP()) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// UserBasedRateLimit uses user ID for authenticated users, IP for guests
func UserBasedRateLimit(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		var key string

		// Try to get user ID from context (if authenticated)
		if userID, exists := c.Get("user_id"); exists {
			key = "user_" + string(rune(userID.(int)))
		} else {
			key = "ip_" + c.ClientIP()
		}

		if !limiter.Allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
