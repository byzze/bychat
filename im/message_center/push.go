package messagecenter

import (
	"bychat/im/cache"
	"bychat/im/client"
	"bychat/im/models"
	"bychat/pkg/common"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// ReponseMsg 发送数据
func ReponseMsg(c *client.Client, code uint32, msgSeq, codeMsg string, message interface{}, msgCmd models.MessageCmd) {
	responseHead := models.NewResponse(msgSeq, code, codeMsg, message, msgCmd)
	headByte, err := json.Marshal(responseHead)
	if err != nil {
		logrus.Error("处理数据 json Marshal", err)
		return
	}
	codeMsg = common.GetErrorMessage(code, codeMsg)
	client.SendMsg(c, headByte)
}

// SendMsgAllClient 全员广播
func SendMsgAllClient(appID, roomID, userID uint32, message string) {
	client.SendMsgAllClient(appID, roomID, userID, message)
}

// GetChatRoomUser 获取房间用户
func GetChatRoomUser(roomID uint32) []*models.UserOnline {
	userResList := cache.GetChatRoomUser(roomID)
	return userResList
}

// SendMsgAllServer 全员广播RPC
func SendMsgAllServer(appID, roomID, userID uint32, message string) (err error) {
	return client.SendMsgAllServer(appID, roomID, userID, message)
}
