package user

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestXxx(t *testing.T) {
	logrus.SetReportCaller(true)
	logrus.WithFields(logrus.Fields{
		"param": "param",
	}).Info("定时任务，清理超时连接")
}
