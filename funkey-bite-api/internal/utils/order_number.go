package utils

import (
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
