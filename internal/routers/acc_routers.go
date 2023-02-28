package routers

import "bychat/internal/servers/websocket"

// InitWebsocket 初始化
func InitWebsocket() {
	websocket.Register("login", websocket.Login)
	websocket.Register("heartbeat", websocket.Heartbeat)
}
