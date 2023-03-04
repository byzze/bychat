package routers

import "bychat/internal/servers/websocket"

// InitWebsocket 初始化
func InitWebsocket() {
	websocket.Register("heartbeat", websocket.Heartbeat)
	websocket.Register("bind", websocket.BindUser)
}
