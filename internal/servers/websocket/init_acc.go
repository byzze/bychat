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
	defaultAppID = iota // 默认平台ID web
)
const (
	defaultRoomID = 101 // 默认房间ID
)

var (
	clientManager = NewClientManager()                     // 管理者
	roomIDs       = []uint32{defaultRoomID, 102, 103, 104} // 全部的平台

	appIDs = []uint32{defaultRoomID, 1, 2, 3} // 全部的平台

	serverIP   string
	serverPort string
)

// GetRoomIDs 获取id
func GetRoomIDs() []uint32 {
	return roomIDs
}

// GetAppIds 获取id
func GetAppIds() []uint32 {
	return appIDs
}

// GetServerNode 获取id
func GetServerNode() (server *models.ServerNode) {
	server = models.NewServerNode(serverIP, serverPort)
	return
}

// IsLocal 校验本地
func IsLocal(server *models.ServerNode) (isLocal bool) {
	if server.IP == serverIP && server.Port == serverPort {
		isLocal = true
	}
	return
}

// InRoomIDs 校验是否在房间id
func InRoomIDs(roomID uint32) (inRoomID bool) {
	for _, value := range roomIDs {
		if value == roomID {
			inRoomID = true
			return
		}
	}
	return
}

// GetDefaultRoomID 获取df id
func GetDefaultRoomID() (roomID uint32) {
	roomID = defaultRoomID
	return
}

// GetDefaultAppID 获取df id
func GetDefaultAppID() (appID uint32) {
	appID = defaultAppID
	return
}

// StartWebSocket 启动程序
func StartWebSocket() {
	serverIP = helper.GetServerNodeIP()

	webSocketPort := viper.GetString("app.webSocketPort")
	rpcPort := viper.GetString("app.rpcPort")

	serverPort = rpcPort

	http.HandleFunc("/acc", wsPage)

	// 添加处理程序
	go clientManager.start()
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
	client := NewClient(0, conn.RemoteAddr().String(), "", "", "", "", conn, currentTime)

	go client.read()
	go client.write()

	// 用户连接事件
	clientManager.Register <- client
}
