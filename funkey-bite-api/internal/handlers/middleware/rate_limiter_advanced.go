package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginLimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

type RateLimitConfig struct {
	Limit  int
	Window time.Duration
	Prefix string
}

var (
	PublicRateLimit = RateLimitConfig{
		Limit:  100,
		Window: time.Hour,
		Prefix: "public_",
	}

	AuthRateLimit = RateLimitConfig{
		Limit:  1000,
		Window: time.Hour,
		Prefix: "auth_",
	}

	OrderCreationRateLimit = RateLimitConfig{
		Limit:  10,
		Window: time.Hour,
		Prefix: "order_",
	}

	TrackingRateLimit = RateLimitConfig{
		Limit:  30,
		Window: time.Hour,
		Prefix: "track_",
	}

	// IdentityRateLimit caps requests per phone number or signed-in customer
	// on each route group, however many IP addresses they come from.
	IdentityRateLimit = RateLimitConfig{
		Limit:  10,
		Window: time.Hour,
		Prefix: "identity_",
	}
)

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

// OrderRateLimitMiddleware caps order creation per IP address.
func OrderRateLimitMiddleware(counter RateLimitCounter) gin.HandlerFunc {
	return SharedRateLimit(counter, OrderCreationRateLimit, ByClientIP)
}

// IdentityRateLimitMiddleware caps requests per phone number or signed-in
// customer within group, e.g. "orders" or "auth", so each kind of action has
// its own allowance.
func IdentityRateLimitMiddleware(counter RateLimitCounter, group string) gin.HandlerFunc {
	config := IdentityRateLimit
	config.Prefix = IdentityRateLimit.Prefix + group + "_"
	return SharedRateLimit(counter, config, ByPhoneOrUser)
}

func TrackingRateLimitMiddleware() gin.HandlerFunc {
	return CreateRateLimiter(TrackingRateLimit)
}
