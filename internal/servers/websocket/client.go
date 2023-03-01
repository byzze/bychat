package websocket

import (
	"bychat/internal/models"
	"runtime/debug"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	// 用户连接超时时间
	heartbeatExpirationTime = 6 * 60
)

// Client 用户连接
type Client struct {
	AppID         uint32          // 登录的平台ID app/web/ios
	Addr          string          // 客户端地址
	Socket        *websocket.Conn // 用户连接
	Send          chan []byte     // 待发送的数据
	FirstTime     uint64          // 首次连接事件
	HeartbeatTime uint64          // 用户上次心跳时间
	LoginTime     uint64          // 登录时间 登录以后才有
	UserOnline    *models.UserOnline
}

// NewClient 初始化
func NewClient(addr string, socket *websocket.Conn, firstTime uint64) (client *Client) {
	client = &Client{
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, 100),
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
	}

	return
}

// GetKey 获取 key
func (c *Client) GetKey() (key string) {
	key = GetUserKey(c.AppID, c.UserOnline.UserID)
	return
}

// 读取客户端数据
func (c *Client) read() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("write stop", string(debug.Stack()), r)
		}
	}()

	defer func() {
		logrus.Info("读取客户端数据 关闭send", c)
		// close(c.Send)
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Addr": c.Addr,
				"err":  err,
			}).Error("读取客户端数据 错误")
			return
		}

		// 处理程序
		logrus.Info("读取客户端数据 处理:", string(message))
		ProcessData(c, message)
	}
}

// 向客户端写数据
func (c *Client) write() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error("write stop", string(debug.Stack()), r)
		}
	}()

	defer func() {
		clientManager.Unregister <- c
		c.Socket.Close()
		logrus.Info("Client发送数据 defer", c)
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 发送数据错误 关闭连接
				logrus.Info("Client发送数据 关闭连接:", c.Addr, "ok:", ok)
				return
			}
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// SendMsg 发送数据
func (c *Client) SendMsg(msg []byte) {
	if c == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			logrus.Error("SendMsg stop:", r, string(debug.Stack()))
		}
	}()

	c.Send <- msg
}

// close 关闭客户端连接
func (c *Client) close() {
	close(c.Send)
}

// Login 用户登录
func (c *Client) Login(appID uint32, userOnline *models.UserOnline) {
	c.UserOnline = userOnline
	c.LoginTime = userOnline.LoginTime
	c.AppID = appID
	// 登录成功=心跳一次
	c.Heartbeat(c.LoginTime)
}

// Heartbeat 用户心跳
func (c *Client) Heartbeat(currentTime uint64) {
	c.HeartbeatTime = currentTime
	return
}

// IsHeartbeatTimeout 心跳超时
func (c *Client) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	if c.HeartbeatTime+heartbeatExpirationTime <= currentTime {
		timeout = true
	}
	return
}

// IsLogin 是否登录了
func (c *Client) IsLogin() (isLogin bool) {
	// 用户登录了
	if c.UserOnline != nil {
		isLogin = true
		return
	}
	return
}

// DeleteRoomUser 移出房间，进入新房间
func (c *Client) DeleteRoomUser(userOnline *models.UserOnline) {
	// 是否在房间
	if c.UserOnline != nil {
		if list, ok := clientManager.Rooms[c.UserOnline.RoomID]; ok {
			tmpList := DeleteSlice(list, userOnline)
			clientManager.Rooms[c.UserOnline.RoomID] = tmpList
		}
	}
}

// DeleteSlice 移位法删除元素
func DeleteSlice(arr []*models.UserOnline, elem *models.UserOnline) []*models.UserOnline {
	j := 0
	for _, v := range arr {
		if v != elem {
			arr[j] = v
			j++
		}
	}
	return arr[:j]
}
