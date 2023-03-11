package user

import (
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	dst := "sda"
	tokenFileName := "sda"
	dst = fmt.Sprintf(dst+"%s", tokenFileName)
	fmt.Println(dst)
	// logrus.SetReportCaller(true)
	// logrus.WithFields(logrus.Fields{
	// 	"param": "param",
	// }).Info("定时任务，清理超时连接")
}
