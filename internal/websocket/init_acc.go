package websocket

import (
	"bychat/internal/common"
	"bychat/internal/models"
	"bychat/internal/utils"
	"encoding/json"
	"net/http"
	"sync"
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
	clientManager = NewClientManager() // 管理者
	serverIP      string
	serverPort    string
)

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

// StartWebSocket 启动程序
func StartWebSocket() {
	serverIP = utils.GetServerNodeIP()

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
	client := NewClient(0, serverIP, serverPort, conn.RemoteAddr().String(), conn, currentTime)

	go client.read()
	go client.write()

	// 用户连接事件
	clientManager.Register <- client
}

// DisposeFunc 处理函数
type DisposeFunc func(client *Client, seq string, message []byte) (code uint32, msg string, data interface{})

var (
	handlers        = make(map[models.MessageCmd]DisposeFunc)
	handlersRWMutex sync.RWMutex
)

// Register 注册
func Register(key models.MessageCmd, value DisposeFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value

	return
}

func getHandlers(key models.MessageCmd) (value DisposeFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()

	value, ok = handlers[key]
	return
}

// ProcessData websocket处理数据
func ProcessData(c *Client, message []byte) {
	logrus.WithFields(logrus.Fields{
		"addr": c.Addr,
		"data": string(message),
	}).Info("ProcessData Request")

	var req = &models.Request{}
	err := json.Unmarshal(message, req)
	if err != nil {
		logrus.Error(err)
		return
	}
	requestData, err := json.Marshal(req.Data)
	if err != nil {
		logrus.Error("处理数据 json Marshal", err)
		c.SendMsg([]byte("处理数据失败"))
		return
	}

	seq := req.MsgSeq
	cmd := models.MessageCmd(req.Cmd)

	var (
		code uint32
		msg  string
		data interface{}
	)

	if v, ok := getHandlers(cmd); ok {
		code, msg, data = v(c, seq, requestData)
	} else {
		code = common.RoutingNotExist
		logrus.WithFields(logrus.Fields{
			"client.Addr": c.Addr,
			"cmd":         cmd,
		}).Error("处理数据 路由不存在")
	}

	msg = common.GetErrorMessage(code, msg)

	responseHead := models.NewResponse(seq, code, msg, data, cmd)

	headByte, err := json.Marshal(responseHead)
	if err != nil {
		logrus.Error("处理数据 json Marshal", err)
		return
	}

	c.SendMsg(headByte)

	logrus.WithFields(logrus.Fields{
		"cmd":      cmd,
		"code":     code,
		"headByte": string(headByte),
	}).Info("acc_response send")
}
