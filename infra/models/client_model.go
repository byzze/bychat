package models

import (
	"bychat/pkg/common"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Client client 管理
const (
	//HeartbeatExpirationTime 用户连接超时时间
	HeartbeatExpirationTime = 60
)

// Client 用户连接
type Client struct {
	AppID   uint32 `json:"appID"`   // 登录的平台ID app/web/ios
	UserID  uint32 `json:"userID"`  // userID
	AccIP   string `json:"accIp"`   // acc Ip
	AccPort string `json:"accPort"` // acc 端口
	// ClientIP      string          `json:"clientIp"`                // 客户端Ip
	// ClientPort    string          `json:"clientPort"`              // 客户端端口
	Addr          string          `json:"addr,omitempty"`          // 客户端地址
	Socket        *websocket.Conn `json:"socket,omitempty"`        // 用户连接
	Send          chan []byte     `json:"send,omitempty"`          // 待发送的数据
	FirstTime     uint64          `json:"firstTime,omitempty"`     // 首次连接事件
	HeartbeatTime uint64          `json:"heartbeatTime,omitempty"` // 用户上次心跳时间
	LoginTime     uint64          `json:"loginTime,omitempty"`     // 登录时间 登录以后才有
}

// NewClient 初始化
func NewClient(appID uint32, accIP, accPort, addr string, socket *websocket.Conn, firstTime uint64) (client *Client) {
	client = &Client{
		AppID:   appID,
		AccIP:   accIP,
		AccPort: accPort,
		// ClientIP:      ClientIP,
		// ClientPort:    ClientPort,
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, 100),
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
	}
	return
}

// GetKey 获取client key
func (client *Client) GetKey() string {
	key := fmt.Sprintf("%d_%d", client.AppID, client.UserID)
	return key
}

// 读取客户端数据
func (client *Client) Read() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("write stop", string(debug.Stack()), r)
		}
	}()

	defer func() {
		logrus.WithFields(logrus.Fields{
			"client.Addr":   client.Addr,
			"client.UserID": client.UserID,
		}).Info("读取客户端数据 关闭send")
		// close(c.Send)
	}()

	for {
		_, message, err := client.Socket.ReadMessage()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Addr": client.Addr,
				"err":  err,
			}).Error("读取客户端数据 错误")
			return
		}

		// 处理程序
		logrus.Info("读取客户端数据 处理:", string(message))
		client.ProcessData(message)
	}
}

// 向客户端写数据
func (client *Client) Write() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("write stop", string(debug.Stack()), r)
		}
	}()

	defer func() {
		// clientManager.Unregister <- client
		// client.Socket.Close()
		logrus.WithFields(logrus.Fields{
			"client.Addr":   client.Addr,
			"client.UserID": client.UserID,
		}).Info("Client发送数据 defer")
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				// 发送数据错误 关闭连接
				logrus.Info("Client发送数据 关闭连接:", client.Addr, "ok:", ok)
				return
			}
			client.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// SendMsg 发送数据
func (client *Client) SendMsg(msg []byte) {
	if client == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			logrus.Error("SendMsg stop:", r, string(debug.Stack()))
		}
	}()

	client.Send <- msg
}

// close 关闭客户端连接
func (client *Client) close() {
	close(client.Send)
}

// Login 用户登录
func (client *Client) Login(appID uint32, userOnline *UserOnline) {
	userOnline.Addr = client.Addr

	client.LoginTime = userOnline.LoginTime
	client.AppID = appID
	client.UserID = userOnline.ID
	// 登录成功=心跳一次
	client.Heartbeat(client.LoginTime)
}

// Heartbeat 用户心跳
func (client *Client) Heartbeat(currentTime uint64) {
	client.HeartbeatTime = currentTime
}

// IsHeartbeatTimeout 心跳超时
func (client *Client) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	if client.HeartbeatTime+HeartbeatExpirationTime <= currentTime {
		timeout = true
	}
	return
}

// UserIsLocal 用户是否在本台机器上
func (client *Client) UserIsLocal(localIP, localPort string) (result bool) {
	if client.AccIP == localIP && client.AccPort == localPort {
		result = true
	}
	return
}

// DisposeFunc 处理函数
type DisposeFunc func(client *Client, seq string, message []byte) (code uint32, msg string, data interface{})

var (
	handlers        = make(map[MessageCmd]DisposeFunc)
	handlersRWMutex sync.RWMutex
)

// Register 注册
func Register(key MessageCmd, value DisposeFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value

	return
}

func getHandlers(key MessageCmd) (value DisposeFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()

	value, ok = handlers[key]
	return
}

// ProcessData websocket处理数据
func (client *Client) ProcessData(message []byte) {
	logrus.WithFields(logrus.Fields{
		"addr": client.Addr,
		"data": string(message),
	}).Info("ProcessData Request")

	var req = &Request{}
	err := json.Unmarshal(message, req)
	if err != nil {
		logrus.Error(err)
		return
	}
	requestData, err := json.Marshal(req.Data)
	if err != nil {
		logrus.Error("处理数据 json Marshal", err)
		client.SendMsg([]byte("处理数据失败"))
		return
	}

	seq := req.MsgSeq
	cmd := MessageCmd(req.Cmd)

	var (
		code uint32
		msg  string
		data interface{}
	)

	if v, ok := getHandlers(cmd); ok {
		code, msg, data = v(client, seq, requestData)
	} else {
		code = common.RoutingNotExist
		logrus.WithFields(logrus.Fields{
			"client.Addr": client.Addr,
			"cmd":         cmd,
		}).Error("处理数据 路由不存在")
	}

	msg = common.GetErrorMessage(code, msg)

	responseHead := NewResponse(seq, code, msg, data, cmd)

	headByte, err := json.Marshal(responseHead)
	if err != nil {
		logrus.Error("处理数据 json Marshal", err)
		return
	}

	client.SendMsg(headByte)

	logrus.WithFields(logrus.Fields{
		"cmd":      cmd,
		"code":     code,
		"headByte": string(headByte),
	}).Info("acc_response send")
}
