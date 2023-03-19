package routers

import (
	"bychat/infra/models"
	"bychat/internal/api/websocket"
)

// InitWebsocket 初始化
func InitWebsocket() {
	models.Register("heartbeat", websocket.Heartbeat)
	models.Register("bind", websocket.BindUser)
}
