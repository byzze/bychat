package im

import (
	"bychat/im/client"
	messagecenter "bychat/im/message_center"
	"bychat/im/models"
	"bychat/im/router"
	"bychat/im/task"
	"bychat/pkg/utils"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultAppID = iota // 默认平台ID web
)
const (
	defaultRoomID = 101 // 默认房间ID
)

var (
	serverIP   string
	serverPort string
)

// StartWebSocket 启动程序
func StartWebSocket() {
	serverIP = utils.GetServerNodeIP()

	webSocketPort := viper.GetString("app.webSocketPort")
	rpcPort := viper.GetString("app.rpcPort")

	serverPort = rpcPort
	// 初始化路由``
	router.InitWebsocket()

	http.HandleFunc("/acc", wsPage)

	models.NewServerNode(serverIP, serverPort)

	task.ServerNodeInit()
	task.CleanConnctionInit()
	// 添加处理程序
	go client.ManagerStart()
	logrus.Infof("WebSocket 启动程序成功:%s:%s", serverIP, serverPort)

	http.ListenAndServe(":"+webSocketPort, nil)
}

func wsPage(w http.ResponseWriter, req *http.Request) {
	// 升级协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		logrus.Info("升级协议", "ua:", r.Header["User-Agent"], "referer:", r.Header["Referer"])
		return true
	}}).Upgrade(w, req, nil)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	logrus.Info("webSocket 建立连接:", conn.RemoteAddr().String())

	currentTime := uint64(time.Now().Unix())
	c := client.NewClient(0, serverIP, serverPort, conn.RemoteAddr().String(), conn, currentTime)

	go c.Read(messagecenter.ProcessData)
	go c.Write()

	// 用户连接事件
	client.GetManager().Register <- c
}
