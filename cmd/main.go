package main

import (
	"bychat/config"
	"bychat/internal/common"
	"bychat/internal/routers"
	"bychat/internal/servers/websocket"
	"bytes"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(LoggerToFile())
	// 初始化路由
	routers.InitWeb(r)
	routers.InitWebsocket()

	config.InitConfig()
	common.SetOutPutFile(logrus.TraceLevel)

	// redislib.InitRedlisClient()

	go websocket.StartWebSocket()

	// task.ServerNodeInit()
	// task.CleanConnctionInit()

	r.Run()
}

// LoggerToFile 日志中间件
func LoggerToFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 开始时间
		start := time.Now()
		// 请求报文
		var requestBody []byte
		if ctx.Request.Body != nil {
			var err error
			requestBody, err = ctx.GetRawData()
			if err != nil {
				logrus.Warn(map[string]interface{}{"err": err.Error()}, "get http request body error")
			}
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		}
		// 处理请求
		ctx.Next()
		// 结束时间
		end := time.Now()
		logrus.Info(map[string]interface{}{
			"statusCode": ctx.Writer.Status(),
			"cost":       float64(end.Sub(start).Nanoseconds()/1e4) / 100.0,
			"clientIp":   ctx.ClientIP(),
			"method":     ctx.Request.Method,
			"uri":        ctx.Request.RequestURI,
		})
	}
}
