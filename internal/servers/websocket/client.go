package websocket

import (
	"bychat/internal/models"
	"runtime/debug"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Client models.Client

// NewClient 初始化
func NewClient(appID uint32, accIP, accPort, ClientIP, ClientPort, addr string, socket *websocket.Conn, firstTime uint64) (client *Client) {
	client = &Client{
		AppID:         appID,
		AccIP:         accIP,
		AccPort:       accPort,
		ClientIP:      ClientIP,
		ClientPort:    ClientPort,
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, 100),
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
	}
	return
}

// 读取客户端数据
func (client *Client) read() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("write stop", string(debug.Stack()), r)
		}
	}()

	defer func() {
		logrus.Info("读取客户端数据 关闭send", client)
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
		ProcessData(client, message)
	}
}

// 向客户端写数据
func (client *Client) write() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("write stop", string(debug.Stack()), r)
		}
	}()

	defer func() {
		clientManager.Unregister <- client
		client.Socket.Close()
		logrus.Info("Client发送数据 defer", client)
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
func (client *Client) Login(appID uint32, userOnline *models.UserOnline) {
	client.LoginTime = userOnline.LoginTime
	client.AppID = appID
	// 登录成功=心跳一次
	client.Heartbeat(client.LoginTime)
}

// Heartbeat 用户心跳
func (client *Client) Heartbeat(currentTime uint64) {
	client.HeartbeatTime = currentTime
	return
}

// IsHeartbeatTimeout 心跳超时
func (client *Client) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	if client.HeartbeatTime+models.HeartbeatExpirationTime <= currentTime {
		timeout = true
	}
	return
}

// IsLogin 是否登录了
func (client *Client) IsLogin() (isLogin bool) {
	// 用户登录了
	return
}

// UserIsLocal 用户是否在本台机器上
func (client *Client) UserIsLocal(localIP, localPort string) (result bool) {
	if client.AccIP == localIP && client.AccPort == localPort {
		result = true

		return
	}
	return
}
