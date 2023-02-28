package websocket

import (
	"bychat/internal/models"
	"bychat/lib/cache"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

// GetUserList 获取全部用户
func GetUserList(appID uint32) (userList []string) {
	logrus.Info("获取全部用户", appID)

	userList = clientManager.GetUserList(appID)
	return
}

// GetUserClient 获取用户所在的连接
func GetUserClient(appID uint32, userID string) (client *Client) {
	client = clientManager.GetUserClient(appID, userID)
	return
}

// SendUserMessage 给用户发送消息
func SendUserMessage(appID uint32, userID string, msgID, message string) (sendResults bool, err error) {
	// 封装发生数据格式
	data := models.GetTextMsgData(userID, msgID, message)
	// 获取与用户建立的socket client，如果不为空，则是当前机器，否则需要通过redis查找对应的服务，并通过rpc发生消息
	client := GetUserClient(appID, userID)

	if client != nil {
		// 在本机发送
		sendResults, err = SendUserMessageLocal(appID, userID, data)
		if err != nil {
			fmt.Println("给用户发送消息", appID, userID, err)
		}
		return
	}

	// key := GetUserKey(appId, userId)
	// info, err := cache.GetUserOnlineInfo(key)
	// if err != nil {
	// 	fmt.Println("给用户发送消息失败", key, err)

	// 	return false, nil
	// }
	// if !info.IsOnline() {
	// 	fmt.Println("用户不在线", key)
	// 	return false, nil
	// }
	// server := models.NewServer(info.AccIp, info.AccPort)
	// msg, err := grpcclient.SendMsg(server, msgId, appId, userId, models.MessageCmdMsg, models.MessageCmdMsg, message)
	// if err != nil {
	// 	fmt.Println("给用户发送消息失败", key, err)

	// 	return false, err
	// }
	// fmt.Println("给用户发送消息成功-rpc", msg)
	// sendResults = true

	return
}

// SendUserMessageLocal 给本机用户发送消息
func SendUserMessageLocal(appID uint32, userID string, data string) (sendResults bool, err error) {
	client := GetUserClient(appID, userID)
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
func SendUserMessageAll(appID uint32, userID, msgID, cmd, message string) (sendResults bool, err error) {
	sendResults = true

	// currentTime := uint64(time.Now().Unix())

	// servers, err := cache.GetServerAll(currentTime)
	// if err != nil {
	// 	fmt.Println("给全体用户发消息", err)

	// 	return
	// }

	logrus.WithFields(logrus.Fields{
		"appID":  appID,
		"userID": userID,
	}).Info("SendUserMessageAll")

	data := models.GetMsgData(userID, msgID, cmd, message)

	cache.ZSetMessage(appID, data)

	AllSendMessages(appID, userID, data)
	// TODO 发送grpc所有人
	return
}

// AllSendMessages 全员广播
func AllSendMessages(appID uint32, userID string, data string) {
	logrus.WithFields(logrus.Fields{
		"appID":  appID,
		"userID": userID,
		"data":   data,
	}).Info("全员广播")

	// 获取userId对应的client，用于过滤
	ignoreClient := clientManager.GetUserClient(appID, userID)
	// 发生数据给所有人
	clientManager.sendAppIDAll([]byte(data), appID, ignoreClient)
}
