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

	for _, v := range cache.GetRoomUser(roomID) {
		userList = append(userList, v.NickName)
	}

	return
}

// GetUserClient 获取用户所在的连接
func GetUserClient(appID, userID uint32) (client *Client) {
	client = clientManager.GetUserClient(appID, userID)
	return
}

// SendUserMessage 给用户发送消息
func SendUserMessage(roomID, userID uint32, msgID, message string) (sendResults bool, err error) {
	user, err := cache.GetUserOnlineInfo(userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"err":    err,
		}).Error("redis 获取失败")
		return false, nil
	}
	nickname := user.NickName
	// 封装发生数据格式
	data := models.GetTextMsgData(nickname, msgID, message)
	// 获取与用户建立的socket client，如果不为空，则是当前机器，否则需要通过redis查找对应的服务，并通过rpc发生消息
	client := GetUserClient(roomID, userID)

	if client != nil {
		// 在本机发送
		sendResults, err = SendUserMessageLocal(roomID, userID, data)
		if err != nil {
			fmt.Println("给用户发送消息", roomID, nickname, err)
		}
		return
	}

	info, err := cache.GetUserOnlineInfo(userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"err":    err,
		}).Error("给用户发送消息失败")
		return false, nil
	}
	if !info.IsOnline() {
		fmt.Println("用户不在线", userID)
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
func SendUserMessageLocal(roomID, userID uint32, data string) (sendResults bool, err error) {
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
func SendUserMessageAll(appID, roomID, userID uint32, msgID, cmd, message string) (sendResults bool, err error) {
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

	uo, err := cache.GetUserOnlineInfo(userID)
	if err != nil {
		logrus.Error("给全体用户发消息", err)
		return
	}
	data := models.GetMsgData(uo.NickName, msgID, cmd, message)
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
func AllSendMessages(appID, roomID, userID uint32, data string) {
	logrus.WithFields(logrus.Fields{
		"appID":  appID,
		"roomID": roomID,
		"userID": userID,
		"data":   data,
	}).Info("全员广播")

	// 获取userId对应的client，用于过滤
	ignoreClient := clientManager.GetUserClient(appID, userID)
	// 发送数据给房间所有人
	clientManager.sendAll([]byte(data), ignoreClient)
}

// EnterRoom 进入房间
func EnterRoom(appID, roomID, userID uint32) error {
	logrus.WithFields(logrus.Fields{
		"AppId":  appID,
		"UserId": userID,
		"RoomID": roomID,
	}).Info("webSocket_request 进入房间接口")
	uo, err := cache.GetUserOnlineInfo(userID)
	if err != nil {
		logrus.Error("EnterRoom Failed:", err)
		return err
	}
	logrus.Info("EnterRoom uo:", uo.ID)
	cache.SetRoomUser(roomID, uo)
	return nil
}

// ExitRoom 离开房间
func ExitRoom(appID, roomID, userID uint32) {
	logrus.WithFields(logrus.Fields{
		"AppId":  appID,
		"UserId": userID,
	}).Info("webSocket_request 离开房间接口")
	uo, err := cache.GetUserOnlineInfo(userID)
	if err != nil {
		logrus.Error("EnterRoom Failed:", err)
	}
	cache.DelRoomUser(roomID, uo)
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
	clientManager.Unregister <- client
}

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
