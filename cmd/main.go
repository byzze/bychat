package main

import (
	"bychat/config"
	"bychat/internal/routers"
	"bychat/internal/servers/websocket"
	"bychat/lib/redislib"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// common.InitLogger() 初始化日志服务
	logrus.SetReportCaller(true)

	config.InitConfig()
	redislib.ExampleNewClient()

	r := gin.Default()
	// 初始化路由
	routers.InitWeb(r)
	routers.InitWebsocket()

	go websocket.StartWebSocket()

	r.Run()
}
