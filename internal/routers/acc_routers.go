package routers

import "bychat/internal/websocket"

// InitWebsocket 初始化
func InitWebsocket() {
	websocket.Register("heartbeat", websocket.Heartbeat)
	websocket.Register("bind", websocket.BindUser)
}
