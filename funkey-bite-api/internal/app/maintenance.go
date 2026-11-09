package app

import (
	"context"
	"log"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/repository"
)

const (
	maintenanceInterval = 15 * time.Minute
	// Rate limit windows are an hour at most; a day leaves room for any limit
	// added later.
	rateLimitRetention = 24 * time.Hour
)

// startMaintenance clears expired idempotency keys and old rate limit
// windows every maintenanceInterval until ctx is cancelled.
func startMaintenance(ctx context.Context, idempotency *repository.IdempotencyRepository, rateLimits *repository.RateLimitRepository) {
	go func() {
		ticker := time.NewTicker(maintenanceInterval)
		defer ticker.Stop()
		for {
			if removed, err := idempotency.DeleteExpired(ctx); err != nil {
				log.Printf("[maintenance] %v", err)
			} else if removed > 0 {
				log.Printf("[maintenance] removed %d expired idempotency keys", removed)
			}
			if _, err := rateLimits.DeleteBefore(ctx, time.Now().Add(-rateLimitRetention)); err != nil {
				log.Printf("[maintenance] %v", err)
			}

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
}
