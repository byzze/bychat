package ws

import (
	"bychat/infra/models"
	"bychat/infra/rpc/grpcclient"
	"bychat/internal/cache"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

// GetUserClient 获取用户client
func GetUserClient(appID, userID uint32) *models.Client {
	return GetClientManager().GetUserClient(appID, userID)
}

// SendMessageLocalClient 给本机用户发送消息
func SendMessageLocalClient(appID, userID uint32, data string) (err error) {
	client := GetUserClient(appID, userID)
	if client == nil {
		err = errors.New("用户不在线")
		return
	}

	// 发送消息
	client.SendMsg([]byte(data))
	return
}

// SendMsgAllClient 全员广播
func SendMsgAllClient(appID, roomID, userID uint32, message string) {
	GetClientManager().sendAll([]byte(message), appID, roomID, userID)
}

// GetChatRoomUser 获取房间用户
func GetChatRoomUser(roomID uint32) []*models.UserOnline {
	userResList := cache.GetChatRoomUser(roomID)
	return userResList
}

// SendMsgAllServer 全员广播RPC
func SendMsgAllServer(appID, roomID, userID uint32, message string) (err error) {
	currentTime := uint64(time.Now().Unix())
	servers, err := GetServerNodeAll(currentTime)
	if err != nil {
		logrus.Error("给全体用户发消息", err)
		return
	}
	for _, sv := range servers {
		if IsLocal(sv) {
			SendMsgAllClient(appID, roomID, userID, message)
		} else {
			err = grpcclient.SendMsgAll(sv, appID, roomID, userID, message)
			if err != nil {
				logrus.WithError(err).Error("rpc SendMsgAll 给用户发消息")
				continue
			}
		}
	}
	return
}
