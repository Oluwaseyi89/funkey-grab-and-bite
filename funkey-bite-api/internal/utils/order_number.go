package utils

import (
	"fmt"
	"time"
)

func GenerateOrderNumber() string {
	year := time.Now().Year()
	month := time.Now().Month()

	seqNum := getNextOrderSequence()

	return fmt.Sprintf("FG-%d-%02d-%04d", year, month, seqNum)
}

func getNextOrderSequence() int {
	return int(time.Now().UnixNano() % 10000)
}

func GenerateCustomerID(userID int) string {
	return fmt.Sprintf("CUST-%06d", userID)
}
