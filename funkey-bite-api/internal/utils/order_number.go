package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

// globalOrderCounter is an atomically-incremented counter seeded from
// the current millisecond timestamp so that it remains distinct across
// process restarts within the same calendar month.
var globalOrderCounter atomic.Int64

func init() {
	globalOrderCounter.Store(time.Now().UnixMilli())
}

func GenerateOrderNumber() string {
	year := time.Now().Year()
	month := time.Now().Month()
	seqNum := getNextOrderSequence()
	return fmt.Sprintf("FG-%d-%02d-%d", year, month, seqNum)
}

// getNextOrderSequence returns a monotonically increasing value that is
// unique within the lifetime of the process and extremely unlikely to
// collide across restarts (restarts land at different millisecond offsets).
func getNextOrderSequence() int64 {
	return globalOrderCounter.Add(1)
}

func GenerateCustomerID(userID int) string {
	return fmt.Sprintf("CUST-%06d", userID)
}

// GeneratePaymentReference builds a Paystack transaction reference we control,
// rather than letting Paystack assign one, so a webhook can be matched straight
// back to the order it belongs to. The random suffix keeps it unguessable even
// though the order number itself isn't secret.
func GeneratePaymentReference(orderNumber string) string {
	suffix := make([]byte, 6)
	if _, err := rand.Read(suffix); err != nil {
		// crypto/rand failing is effectively unrecoverable on any real platform;
		// fall back to a timestamp so this never panics or blocks an order.
		return fmt.Sprintf("PSK-%s-%d", orderNumber, time.Now().UnixNano())
	}
	return fmt.Sprintf("PSK-%s-%s", orderNumber, hex.EncodeToString(suffix))
}
