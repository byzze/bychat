package websocket

import (
	"bychat/internal/helper"
	"bychat/internal/models"
	"bychat/internal/servers/grpcclient"
	"bychat/lib/cache"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

// Login 登录
func Login(appID, userID uint32, nickName string) error {
	logrus.WithFields(logrus.Fields{
		"AppId":  appID,
		"UserId": userID,
	}).Info("webSocket_request 登录")

	currentTime := uint64(time.Now().Unix())

	var user = models.UserOnline{
		ID:            userID,
		NickName:      nickName,
		LoginTime:     currentTime,
		HeartbeatTime: 0,
		LogOutTime:    0,
		DeviceInfo:    "",
		IsLogoff:      false,
	}
	return cache.SetUserOnlineInfo(userID, &user)
}

// LogOut 退出
func LogOut(appID, userID uint32) {
	logrus.WithFields(logrus.Fields{
		"AppId":  appID,
		"UserId": userID,
	}).Info("webSocket_request 退出")
	// 设置redis缓存
	client := GetUserClient(appID, userID)
	if client == nil {
		return
	}
	unregisterChannel(client)
}

// GetRoomUserList 获取全部用户
func GetRoomUserList(appID, roomID uint32) (userList []string) {
	logrus.WithFields(logrus.Fields{
		"roomID": roomID,
	}).Info("获取全部用户")
	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerNodeAll(currentTime)
	if err != nil {
		logrus.Error("给全体用户发消息", err)
		return
	}
	for i, server := range servers {
		if server.IP == serverIP && server.Port == serverPort {
			for _, v := range cache.GetRoomUser(roomID) {
				userList = append(userList, v.NickName)
			}
		} else {
			roomUserList, err := grpcclient.GetRoomUserList(servers[i], appID, roomID)
			if err != nil {
				logrus.Error("grpcclient GetRoomUserList", err)
				continue
			}
			userList = append(userList, roomUserList...)
		}
	}

	return
}

// GetUserClient 获取用户所在的连接
func GetUserClient(appID, userID uint32) (client *Client) {
	client = clientManager.GetUserClient(appID, userID)
	return
}

// SendUserMessageAll 发送消息 群聊
func SendUserMessageAll(appID, roomID, userID uint32, msgID, cmd, message string) (sendResults bool, err error) {
	sendResults = true
	currentTime := uint64(time.Now().Unix())

	servers, err := cache.GetServerNodeAll(currentTime)
	if err != nil {
		sendResults = false
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

	uo, err := cache.GetUserOnlineInfo(userID)
	if err != nil {
		logrus.Error("给全体用户发消息", err)
		sendResults = false
		return
	}
	var data string

	switch cmd {
	case models.MessageCmdEnter:
		data = models.GetTextMsgDataEnter(uo.NickName, msgID, message)
	case models.MessageCmdExit:
		data = models.GetTextMsgDataExit(uo.NickName, msgID, message)
	default:
		data = models.GetMsgData(uo.NickName, msgID, cmd, message)
		cache.ZSetMessage(roomID, data)
	}

	for _, sv := range servers {
		if IsLocal(sv) {
			AllSendMessages(appID, roomID, userID, data)
		} else {
			_, err := grpcclient.SendMsgAll(sv, appID, roomID, userID, msgID, cmd, message)
			if err != nil {
				logrus.Error("grpcclient SendMsgAll 给全体用户发消息", err)
				sendResults = false
				continue
			}
		}
	}
	return
}

// AllSendMessages 全员广播
func AllSendMessages(appID, roomID, userID uint32, data string) {
	logrus.WithFields(logrus.Fields{
		"appID":  appID,
		"roomID": roomID,
		"userID": userID,
		"data":   data,
	}).Info("AllSendMessages")

	// 获取userId对应的client，用于过滤
	ignoreClient := clientManager.GetUserClient(appID, userID)
	// 发送数据给房间所有人
	clientManager.sendAll([]byte(data), roomID, ignoreClient)
}

// EnterRoom 进入房间
func EnterRoom(appID, roomID, userID uint32) error {
	logrus.WithFields(logrus.Fields{
		"AppId":  appID,
		"UserId": userID,
		"RoomID": roomID,
	}).Info("EnterRoom")
	uo, err := cache.GetUserOnlineInfo(userID)
	if err != nil {
		logrus.Error("EnterRoom Failed:", err)
		return err
	}
	if uo == nil {
		return errors.New("用户未登录")
	}
	seq := helper.GetOrderIDTime()
	sendResults, err := SendUserMessageAll(appID, roomID, userID, seq, models.MessageCmdEnter, "哈喽~")
	if err != nil {
		logrus.Error("SendUserMessageAll Failed:", err)
		return err
	}
	if !sendResults {
		return nil
	}
	logrus.Info("EnterRoom uo:", uo.ID)
	cache.SetRoomUser(roomID, uo)
	return nil
}

// ExitRoom 离开房间
func ExitRoom(appID, roomID, userID uint32) error {
	logrus.WithFields(logrus.Fields{
		"AppId":  appID,
		"UserId": userID,
	}).Info("ExitRoom")
	seq := helper.GetOrderIDTime()
	sendResults, err := SendUserMessageAll(appID, roomID, userID, seq, models.MessageCmdExit, "退出~")
	if err != nil {
		logrus.Error("SendUserMessageAll Failed:", err)
		return err
	}
	if !sendResults {
		return nil
	}

	cache.DelRoomUser(roomID, userID)
	return nil
}
