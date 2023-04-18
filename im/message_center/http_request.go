package messagecenter

import (
	"bychat/im/cache"
	"bychat/im/client"
	"bychat/im/models"
	"bychat/im/rpc/grpcclient"
	"time"

	"github.com/sirupsen/logrus"
)

// UserLogout 退出
func UserLogout(appID, userID uint32) {
	c := client.GetUserClient(appID, userID)
	client.UnregisterChannel(c)
}

// UserLogin 退出
func UserLogin(appID, userID uint32, accIP, accPort string, nickName string, addr string, loginTime uint64) *models.UserOnline {
	return models.UserLogin(appID, userID, "", "", nickName, "", loginTime)
}

// GetRoomUserList 获取用户列表
func GetRoomUserList(appID, roomID uint32) (userList []*models.ResponseUserOnline) {
	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerNodeAll(currentTime)
	if err != nil {
		logrus.WithError(err).Error("GetRoomUserList")
		return
	}

	for i, server := range servers {
		if server.IP == models.ServerNodeInfo.IP && server.Port == models.ServerNodeInfo.Port {
			for _, v := range cache.GetChatRoomUser(roomID) {
				tmp := &models.ResponseUserOnline{
					ID:       v.ID,
					NickName: v.NickName,
					Avatar:   v.Avatar,
				}
				userList = append(userList, tmp)
			}
		} else {
			_, err := grpcclient.GetRoomUserList(servers[i], appID, roomID)
			if err != nil {
				logrus.WithError(err).Error("rpc GetRoomUserList")
				continue
			}
			// TODO 封装
			// userList = append(userList, roomUserList...)
		}
	}
	return
}
