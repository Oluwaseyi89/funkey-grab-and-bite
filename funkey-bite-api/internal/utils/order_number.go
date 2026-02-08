// internal/utils/order_number.go
package utils

import (
	"fmt"
	// "strconv"
	"time"
)

func GenerateOrderNumber() string {
	year := time.Now().Year()
	month := time.Now().Month()

	// Get sequence number (from DB sequence)
	seqNum := getNextOrderSequence()

	return fmt.Sprintf("FG-%d-%02d-%04d", year, month, seqNum)
}

func getNextOrderSequence() int {
	// This would query the PostgreSQL sequence
	// For now, we'll use a timestamp-based approach
	return int(time.Now().UnixNano() % 10000)
}

// Alternative: Pattern-based like your example
func GenerateCustomerID(userID int) string {
	return fmt.Sprintf("CUST-%06d", userID)
}
