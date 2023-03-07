package websocket

import (
	"bychat/internal/grpcclient"
	"bychat/internal/helper"
	"bychat/internal/models"
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
func GetRoomUserList(appID, roomID uint32) (userList []*models.ResponseUserOnline) {
	logrus.WithFields(logrus.Fields{
		"roomID": roomID,
	}).Info("GetRoomUserList")
	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerNodeAll(currentTime)
	if err != nil {
		logrus.Error("给全体用户发消息", err)
		return
	}

	for i, server := range servers {
		if server.IP == serverIP && server.Port == serverPort {
			for _, v := range cache.GetChatRoomUser(roomID) {
				tmp := &models.ResponseUserOnline{
					ID:       v.ID,
					NickName: v.NickName,
					Avatar:   v.Avatar,
				}
				userList = append(userList, tmp)
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
func SendUserMessageAll(appID, roomID, userID uint32, message string) (sendResults bool, err error) {
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
	}).Info("SendUserMessageAll")

	for _, sv := range servers {
		if IsLocal(sv) {
			AllSendMessages(appID, roomID, userID, message)
		} else {
			_, err := grpcclient.SendMsgAll(sv, appID, roomID, userID, message)
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
	if ignoreClient == nil {
		logrus.Error("AllSendMessages")
		return
	}
	// 发送数据给房间所有人
	clientManager.sendAll([]byte(data), roomID, userID, ignoreClient)
}

// EnterChatRoom 进入房间
func EnterChatRoom(appID, roomID, userID uint32) error {
	logrus.WithFields(logrus.Fields{
		"AppId":  appID,
		"UserId": userID,
		"RoomID": roomID,
	}).Info("EnterChatRoom")
	uo, err := cache.GetUserOnlineInfo(userID)
	if err != nil {
		logrus.Error("EnterChatRoom Failed:", err)
		return err
	}
	if uo == nil {
		return errors.New("用户未登录")
	}

	seq := helper.GetOrderIDTime()
	data := models.GetTextMsgDataEnter(uo.NickName, seq, "哈喽~")

	cache.SetChatRoomUser(roomID, uo)
	logrus.Info("EnterChatRoom uo:", uo.ID)

	sendResults, err := SendUserMessageAll(appID, roomID, userID, data)
	if err != nil {
		logrus.Error("SendUserMessageAll Failed:", err)
		return err
	}
	if !sendResults {
		return errors.New("发送失败")
	}
	return nil
}

// ExitChatRoom 离开房间
func ExitChatRoom(appID, roomID, userID uint32) error {
	logrus.WithFields(logrus.Fields{
		"AppId":  appID,
		"UserId": userID,
	}).Info("ExitChatRoom")
	seq := helper.GetOrderIDTime()
	uo, err := cache.GetUserOnlineInfo(userID)
	if err != nil {
		logrus.Error("EnterChatRoom Failed:", err)
		return err
	}
	if uo == nil {
		return errors.New("用户未登录")
	}

	data := models.GetTextMsgDataExit(uo.NickName, seq, "退出~")
	sendResults, err := SendUserMessageAll(appID, roomID, userID, data)
	if err != nil {
		logrus.Error("SendUserMessageAll Failed:", err)
		return err
	}
	if !sendResults {
		return nil
	}

	cache.DelChatRoomUser(roomID, userID)
	return nil
}
