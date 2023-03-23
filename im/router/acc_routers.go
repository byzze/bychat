package router

import (
	messagecenter "bychat/im/message_center"
	"bychat/im/models"
)

// InitWebsocket 初始化
func InitWebsocket() {
	messagecenter.Register(models.MessageCmdHeartbeat, messagecenter.Heartbeat)
	messagecenter.Register(models.MessageCmdBindUser, messagecenter.Login)
	// messagecenter.Register("logout", messagecenter.Logout)
	// messagecenter.Register("msg", messagecenter.MsgProcess)
}
