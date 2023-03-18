package websocket

import (
	"bychat/internal/grpcclient"
	"bychat/internal/models"
	"bychat/pkg/cache"
	"time"

	"github.com/sirupsen/logrus"
)

// GetUserClient 获取用户client
func GetUserClient(appID, userID uint32) *Client {
	return clientManager.getUserClient(appID, userID)
}

// SendMsgAllClient 全员广播
func SendMsgAllClient(appID, roomID, userID uint32, message string) {
	clientManager.sendAll([]byte(message), appID, roomID, userID)
}

// SendMsgAllServer 全员广播
func SendMsgAllServer(appID, roomID, userID uint32, message string) (err error) {
	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerNodeAll(currentTime)
	if err != nil {
		logrus.Error("给全体用户发消息", err)
		return
	}
	for _, sv := range servers {
		if IsLocal(sv) {
			SendMsgAllClient(appID, roomID, userID, message)
		} else {
			err := RPCSendMsgAll(sv, appID, roomID, userID, message)
			if err != nil {
				logrus.WithError(err).Error("grpcclient SendMsgAll 给全体用户发消息")
				continue
			}
		}
	}
	return
}

/************GRPC*****************/

// RPCSendMsgAll 发送数据
func RPCSendMsgAll(sv *models.ServerNode, appID, roomID, userID uint32, message string) (err error) {
	err = grpcclient.SendMsgAll(sv, appID, roomID, userID, message)
	if err != nil {
		logrus.Error("grpcclient SendMsgAll 给全体用户发消息", err)
	}
	return
}
