package task

import (
	"bychat/im/cache"
	"bychat/im/models"
	"log"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

// ServerNodeInit 服务器初始化
func ServerNodeInit() {
	Timer(2*time.Second, 60*time.Second, serverRegister, "", serverDefer, "")
}

// 服务注册
func serverRegister(param interface{}) (result bool) {
	result = true

	defer func() {
		if r := recover(); r != nil {
			log.Println(string(debug.Stack()))
			logrus.WithFields(logrus.Fields{
				"r": r,
			}).Info("服务注册 stop")
		}
	}()

	serverNode := models.ServerNodeInfo

	currentTime := uint64(time.Now().Unix())

	logrus.WithFields(logrus.Fields{
		"param":       param,
		"server":      serverNode,
		"currentTime": currentTime,
	}).Info("定时任务，服务注册")

	cache.SetServerNodeInfo(serverNode, currentTime)
	return
}

// 服务下线
func serverDefer(param interface{}) (result bool) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithFields(logrus.Fields{
				"r":           r,
				"debug.Stack": string(debug.Stack()),
			}).Info("服务下线 stop")
		}
	}()

	logrus.WithFields(logrus.Fields{
		"param": param,
	}).Info("服务下线")

	server := models.ServerNodeInfo

	cache.DelServerNodeInfo(server)
	return
}
