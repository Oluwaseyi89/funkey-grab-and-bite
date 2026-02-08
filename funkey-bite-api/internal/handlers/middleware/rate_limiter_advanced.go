package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginLimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimitConfig defines rate limiting configuration
type RateLimitConfig struct {
	Limit  int
	Window time.Duration
	Prefix string
}

// Default rate limit configurations
var (
	PublicRateLimit = RateLimitConfig{
		Limit:  100, // 100 requests
		Window: time.Hour,
		Prefix: "public_",
	}

	AuthRateLimit = RateLimitConfig{
		Limit:  1000, // 1000 requests
		Window: time.Hour,
		Prefix: "auth_",
	}

	OrderCreationRateLimit = RateLimitConfig{
		Limit:  10, // 10 orders per hour
		Window: time.Hour,
		Prefix: "order_",
	}

	TrackingRateLimit = RateLimitConfig{
		Limit:  30, // 30 tracking requests per hour
		Window: time.Hour,
		Prefix: "track_",
	}
)

// CreateRateLimiter creates a rate limiter middleware
func CreateRateLimiter(config RateLimitConfig) gin.HandlerFunc {
	rate := limiter.Rate{
		Period: config.Window,
		Limit:  int64(config.Limit),
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate, limiter.WithClientIPHeader("X-Forwarded-For"))

	return ginLimiter.NewMiddleware(instance, ginLimiter.WithKeyGetter(func(c *gin.Context) string {
		return config.Prefix + c.ClientIP()
	}))
}

// OrderRateLimitMiddleware specific for order creation
func OrderRateLimitMiddleware() gin.HandlerFunc {
	return CreateRateLimiter(OrderCreationRateLimit)
}

// TrackingRateLimitMiddleware specific for order tracking
func TrackingRateLimitMiddleware() gin.HandlerFunc {
	return CreateRateLimiter(TrackingRateLimit)
}
