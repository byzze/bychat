package user

import (
	"bychat/im/cache"
	messagecenter "bychat/im/message-center"
	"bychat/im/models"
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
	// TODO 数据库，数据校验
	var user = messagecenter.UserLogin(appID, userID, "", "", nickName, "", currentTime)
	return cache.SetUserOnlineInfo(userID, user)
}

// LogOut 退出
func LogOut(appID, userID uint32) {
	logrus.WithFields(logrus.Fields{
		"AppId":  appID,
		"UserId": userID,
	}).Info("webSocket_request 退出")
	// 设置redis缓存
	messagecenter.UserLogout(appID, userID)

}

// GetRoomUserList 获取全部用户
func GetRoomUserList(appID, roomID uint32) (userList []*models.ResponseUserOnline) {
	logrus.WithFields(logrus.Fields{
		"roomID": roomID,
	}).Info("GetRoomUserList")
	userList = messagecenter.GetRoomUserList(appID, roomID)
	return
}

// EnterChatRoom 进入房间
func EnterChatRoom(appID, roomID, userID uint32) error {
	// 记录函数名称和请求参数
	logrus.WithFields(logrus.Fields{
		"appID":  appID,
		"roomID": roomID,
		"userID": userID,
	}).Info("EnterChatRoom called")

	// 获取用户在线信息
	/* 	uo, err := cache.GetUserOnlineInfo(userID)
	   	if err != nil {
	   		logrus.WithError(err).Error("EnterChatRoom: failed to get user online info")
	   		return err
	   	}

	   	if uo == nil {
	   		return errors.New("user is not online")
	   	}

	   	seq := utils.GetOrderIDTime()

	   	cache.SetChatRoomUser(roomID, uo)

	   	data := models.GetTextMsgDataEnter(uo.NickName, "", seq, "哈喽~")

	   	// 记录消息序列号并发送消息
	   	logrus.WithFields(logrus.Fields{
	   		"seq": seq,
	   	}).Info("EnterChatRoom: message sent")
	   	sendResults, err := SendMessageAll(appID, roomID, userID, data)
	   	if err != nil {
	   		logrus.WithError(err).Error("EnterChatRoom: failed to send message")
	   		return err
	   	}
	   	if !sendResults {
	   		return errors.New("failed to send message")
	   	}
	*/
	return nil
}

// ExitChatRoom 离开房间
func ExitChatRoom(appID, roomID, userID uint32) error {
	logrus.WithFields(logrus.Fields{
		"AppId":  appID,
		"UserId": userID,
	}).Info("ExitChatRoom")

	/* uo, err := cache.GetUserOnlineInfo(userID)
	if err != nil {
		logrus.WithError(err).Error("EnterChatRoom Failed")
		return err
	}
	if uo == nil {
		return errors.New("user is not online")
	}

	seq := utils.GetOrderIDTime()
	data := models.GetTextMsgDataExit(uo.NickName, "", seq, "退出~")

	cache.DelChatRoomUser(roomID, userID)

	sendResults, err := SendMessageAll(appID, roomID, userID, data)
	if err != nil {
		logrus.WithError(err).Error("ExitChatRoom: SendUserMessageAll Failed")
		return err
	}
	if !sendResults {
		return errors.New("failed to send message")
	}
	*/
	return nil
}

// SendMessageAll 用户发送消息 群聊
func SendMessageAll(appID, roomID, userID uint32, message string) (sendResults bool, err error) {
	sendResults = true
	logrus.WithFields(logrus.Fields{
		"appID":   appID,
		"roomID":  roomID,
		"userID":  userID,
		"message": message,
	}).Info("SendUserMessageAll")

	messagecenter.SendMsgAllServer(appID, roomID, userID, message)
	return
}
