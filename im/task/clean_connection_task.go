package task

import (
	"bychat/im/client"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

/**
1. 获取所有client
2. 计算client心跳时长是否过期，进行相关剔除操作（关闭通道，关闭链接）
3. 定期执行
**/

// CleanConnctionInit 清楚链接
func CleanConnctionInit() {
	Timer(3*time.Second, 60*time.Second, cleanConnection, "", nil, nil)
}

// 清理超时连接
func cleanConnection(param interface{}) (result bool) {
	logrus.WithFields(logrus.Fields{
		"param": param,
	}).Info("定时任务，清理超时连接 启动")

	result = true
	// 需要捕获，否则子协程出现panic时会导致程序奔溃
	defer func() {
		if r := recover(); r != nil {
			logrus.WithFields(logrus.Fields{
				"r":             r,
				"debug.Stack()": string(debug.Stack()),
			}).Info("ClearTimeoutConnections stop")
		}
	}()

	client.ClearTimeoutConnections()

	logrus.WithFields(logrus.Fields{
		"param": param,
	}).Info("定时任务，清理超时连接 完成")
	return
}
