package helper

import (
	"fmt"
	"time"
)

// GetOrderIDTime 获取时间戳
func GetOrderIDTime() (orderID string) {
	currentTime := time.Now().Nanosecond()
	orderID = fmt.Sprintf("%d", currentTime)
	return
}
