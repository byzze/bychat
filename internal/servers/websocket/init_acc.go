package websocket

import (
	"bychat/internal/helper"
	"bychat/internal/models"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultAppId = 101 // 默认平台Id
)

var (
	clientManager = NewClientManager()                    // 管理者
	appIds        = []uint32{defaultAppId, 102, 103, 104} // 全部的平台

	serverIP   string
	serverPort string
)

func GetAppIds() []uint32 {
	return appIds
}

func GetServer() (server *models.Server) {
	server = models.NewServer(serverIP, serverPort)
	return
}

func IsLocal(server *models.Server) (isLocal bool) {
	if server.IP == serverIP && server.Port == serverPort {
		isLocal = true
	}
	return
}

func InAppIds(appID uint32) (inAppID bool) {

	for _, value := range appIds {
		if value == appID {
			inAppID = true
			return
		}
	}
	return
}

func GetDefaultAppID() (appID uint32) {
	appID = defaultAppId

	return
}

// 启动程序
func StartWebSocket() {

	serverIP = helper.GetServerIp()

	webSocketPort := viper.GetString("app.webSocketPort")
	rpcPort := viper.GetString("app.rpcPort")

	serverPort = rpcPort

	http.HandleFunc("/acc", wsPage)

	// 添加处理程序
	go clientManager.start()
	logrus.Info("WebSocket 启动程序成功", serverIP, serverPort)

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
	client := NewClient(conn.RemoteAddr().String(), conn, currentTime)

	go client.read()
	go client.write()

	// 用户连接事件
	clientManager.Register <- client
}
