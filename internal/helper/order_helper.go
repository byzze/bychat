package helper

import (
	"fmt"
	"time"
)

func GetOrderIdTime() (orderID string) {
	currentTime := time.Now().Nanosecond()
	orderID = fmt.Sprintf("%d", currentTime)
	return
}
