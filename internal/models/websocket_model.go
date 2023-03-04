package models

import (
	"github.com/gorilla/websocket"
)

const (
	//HeartbeatExpirationTime 用户连接超时时间
	HeartbeatExpirationTime = 6 * 60
)

// Client 用户连接
type Client struct {
	AppID         uint32          `json:"appID"`                   // 登录的平台ID app/web/ios
	UserID        uint32          `json:"userID"`                  // userID
	AccIP         string          `json:"accIp"`                   // acc Ip
	AccPort       string          `json:"accPort"`                 // acc 端口
	ClientIP      string          `json:"clientIp"`                // 客户端Ip
	ClientPort    string          `json:"clientPort"`              // 客户端端口
	Addr          string          `json:"addr,omitempty"`          // 客户端地址
	Socket        *websocket.Conn `json:"socket,omitempty"`        // 用户连接
	Send          chan []byte     `json:"send,omitempty"`          // 待发送的数据
	FirstTime     uint64          `json:"firstTime,omitempty"`     // 首次连接事件
	HeartbeatTime uint64          `json:"heartbeatTime,omitempty"` // 用户上次心跳时间
	LoginTime     uint64          `json:"loginTime,omitempty"`     // 登录时间 登录以后才有
}
