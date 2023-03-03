package websocket

import (
	"bychat/internal/models"
	"bychat/lib/cache"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// GetRoomUserList 获取全部用户
func GetRoomUserList(appID, roomID uint32) (userList []string) {
	logrus.WithFields(logrus.Fields{
		"roomID": roomID,
	}).Info("获取全部用户")

	// key := roomID
	// for _, v := range cache.Rooms[key] {
	// 	userList = append(userList, v.ID)
	// }
	return
}

// GetUserClient 获取用户所在的连接
func GetUserClient(appID uint32, userID string) (client *Client) {
	client = clientManager.GetUserClient(appID, userID)
	return
}

// SendUserMessage 给用户发送消息
func SendUserMessage(roomID uint32, userID string, msgID, message string) (sendResults bool, err error) {
	// 封装发生数据格式
	data := models.GetTextMsgData(userID, msgID, message)
	// 获取与用户建立的socket client，如果不为空，则是当前机器，否则需要通过redis查找对应的服务，并通过rpc发生消息
	client := GetUserClient(roomID, userID)

	if client != nil {
		// 在本机发送
		sendResults, err = SendUserMessageLocal(roomID, userID, data)
		if err != nil {
			fmt.Println("给用户发送消息", roomID, userID, err)
		}
		return
	}

	key := userID
	info, err := cache.GetUserOnlineInfo(key)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"key": key,
			"err": err,
		}).Error("给用户发送消息失败")
		return false, nil
	}
	if !info.IsOnline() {
		fmt.Println("用户不在线", key)
		return false, nil
	}
	// server := models.NewServer(info.AccIP, info.AccPort)
	// msg, err := grpcclient.SendMsg(server, msgID, userID, models.MessageCmdMsg, models.MessageCmdMsg, message)
	// if err != nil {
	// 	fmt.Println("给用户发送消息失败", key, err)
	// 	return false, err
	// }
	// fmt.Println("给用户发送消息成功-rpc", msg)
	sendResults = true

	return
}

// SendUserMessageLocal 给本机用户发送消息
func SendUserMessageLocal(roomID uint32, userID string, data string) (sendResults bool, err error) {
	client := GetUserClient(roomID, userID)
	if client == nil {
		err = errors.New("用户不在线")
		return
	}

	// 发送消息
	client.SendMsg([]byte(data))
	sendResults = true
	return
}

// SendUserMessageAll 发送消息
func SendUserMessageAll(appID, roomID uint32, userID, msgID, cmd, message string) (sendResults bool, err error) {
	sendResults = true
	currentTime := uint64(time.Now().Unix())

	servers, err := cache.GetServerNodeAll(currentTime)
	if err != nil {
		logrus.Error("给全体用户发消息", err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"servers": servers,
		"appID":   appID,
		"roomID":  roomID,
		"userID":  userID,
		"message": message,
		"cmd":     cmd,
		"msgID":   msgID,
	}).Info("SendUserMessageAll")

	data := models.GetMsgData(userID, msgID, cmd, message)
	cache.ZSetMessage(roomID, data)

	for _, sv := range servers {
		if sv.IP == serverIP && sv.Port == sv.Port {
			AllSendMessages(appID, roomID, userID, data)
		}
		// TODO 发送grpc所有人
	}
	return
}

// AllSendMessages 全员广播
func AllSendMessages(appID, roomID uint32, userID string, data string) {
	logrus.WithFields(logrus.Fields{
		"appID":  appID,
		"roomID": roomID,
		"userID": userID,
		"data":   data,
	}).Info("全员广播")

	// 获取userId对应的client，用于过滤
	ignoreClient := clientManager.Users[userID]
	// 发送数据给房间所有人
	clientManager.sendRoomIDAll([]byte(data), roomID, ignoreClient)
}

// EnterRoom 进入房间
func EnterRoom(appID, roomID uint32, userID string) {
	logrus.WithFields(logrus.Fields{
		"AppId":  appID,
		"UserId": userID,
		"RoomID": roomID,
	}).Info("webSocket_request 进入房间接口")

}

// ExitRoom 进入房间
func ExitRoom(appID uint32, userID string) {
	logrus.WithFields(logrus.Fields{
		"AppId":  appID,
		"UserId": userID,
	}).Info("webSocket_request 离开房间接口")

}

// LogOut 退出
func LogOut(appID uint32, userID string) {
	logrus.WithFields(logrus.Fields{
		"AppId":  appID,
		"UserId": userID,
	}).Info("webSocket_request 退出")
	delete(cache.UserMap, userID)

	c := GetUserClient(appID, userID)
	clientManager.Unregister <- c
}

// Login 登录
func Login(appID uint32, userID, userName string) {
	logrus.WithFields(logrus.Fields{
		"AppId":  appID,
		"UserId": userID,
	}).Info("webSocket_request 登录")

	var user = models.UserOnline{
		ID:            userID,
		LoginTime:     0,
		HeartbeatTime: 0,
		LogOutTime:    0,
		DeviceInfo:    "",
		IsLogoff:      false,
	}
	cache.UserMap[userID] = &user
}
