package routers

import (
	"bychat/internal/websocket"
	"bychat/pkg/api"
)

// InitWebsocket 初始化
func InitWebsocket() {
	websocket.Register("heartbeat", api.Heartbeat)
	websocket.Register("bind", api.BindUser)
}
