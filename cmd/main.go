package main

import (
	config "bychat/configs"
	"bychat/internal/common"
	"bychat/internal/grpcserver"
	"bychat/internal/task"
	"bychat/internal/websocket"
	"bychat/pkg/redislib"
	"bychat/pkg/routers"
	"bytes"
	"flag"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var cname string

func main() {
	flag.StringVar(&cname, "cname", "app", "your name")
	flag.Parse()

	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(LoggerToFile())
	// 初始化路由
	routers.InitWeb(r)
	routers.InitWebsocket()

	config.InitConfig(cname)
	common.SetOutPutFile(logrus.TraceLevel)

	redislib.InitRedlisClient()

	go websocket.StartWebSocket()

	task.ServerNodeInit()
	task.CleanConnctionInit()
	go grpcserver.Init()

	// http
	httpPort := viper.GetString("app.httpPort")
	http.ListenAndServe(":"+httpPort, r)
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
