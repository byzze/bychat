package client

import (
	"bychat/im/cache"
	"bychat/im/models"
	"bychat/im/rpc/grpcclient"
	"time"

	"github.com/sirupsen/logrus"
)

// ManagerStart 管理启动
func ManagerStart() {
	getManager().start()
}

// GetUserClient 对外获取管理者
func GetUserClient(appID, userID uint32) *Client {
	return getManager().getUserClient(appID, userID)
}

// GetManager 对外获取管理者
func GetManager() *Manager {
	return getManager()
}

// SendMsg 针对client 发送消息
func SendMsg(c *Client, headByte []byte) {
	c.sendMsg(headByte)
}

// SendMsgAllClient 发消息至本机的所有client 向全部成员(除了自己)发送数据
func SendMsgAllClient(appID, roomID, userID uint32, message string) {
	logrus.WithFields(logrus.Fields{
		"appID":  appID,
		"roomID": roomID,
		"userID": userID,
	}).Info("sendAll 发送消息")

	roomUserList := cache.GetChatRoomUser(roomID)
	for _, v := range roomUserList {
		c := getManager().getUserClient(appID, userID)
		if c != nil && v.ID != userID {
			c.sendMsg([]byte(message))
		}
	}
}

// SendMsgAllServer 全员广播RPC
func SendMsgAllServer(appID, roomID, userID uint32, message string) (err error) {
	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerNodeAll(currentTime)
	if err != nil {
		logrus.Error("给全体用户发消息", err)
		return
	}
	for _, sv := range servers {
		if models.IsLocal(sv) {
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

// ClearTimeoutConnections 定时清理超时连接
func ClearTimeoutConnections() {
	currentTime := uint64(time.Now().Unix())

	clients := getManager().getClients()
	for client := range clients {
		if client.IsHeartbeatTimeout(currentTime) {
			logrus.WithFields(logrus.Fields{
				"client.Addr":          client.Addr,
				"client.UserID":        client.UserID,
				"client.LoginTime":     client.LoginTime,
				"client.HeartbeatTime": client.HeartbeatTime,
			}).Info("心跳时间超时 关闭连接")
			UnregisterChannel(client)
			// client.Socket.Close()
		}
	}
}

// IsLogin 是否登录了
func IsLogin(client *Client) (isLogin bool) {
	c := GetManager().getUserClient(client.AppID, client.UserID)
	if c != nil {
		isLogin = true
	}
	return
}

// GetManagerInfo 获取管理者信息
func GetManagerInfo(isDebug string) (managerInfo map[string]interface{}) {
	managerInfo = make(map[string]interface{})

	managerInfo["clientsLen"] = getManager().getClientsLen()        // 客户端连接数
	managerInfo["usersLen"] = getManager().getUsersLen()            // 登录用户数
	managerInfo["chanRegisterLen"] = len(getManager().Register)     // 未处理连接事件数
	managerInfo["chanUnregisterLen"] = len(getManager().Unregister) // 未处理退出登录事件数
	// managerInfo["chanBroadcastLen"] = len(clientManager.Broadcast)   // 未处理广播事件数

	if isDebug == "true" {
		addrList := make([]string, 0)
		getManager().clientsRange(func(client *Client, value bool) (result bool) {
			addrList = append(addrList, client.Addr)

			return true
		})

		users := getManager().getUserKeys()

		managerInfo["clients"] = addrList // 客户端列表
		managerInfo["users"] = users      // 登录用户列表
	}

	return
}
