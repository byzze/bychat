package websocket

import (
	"bychat/internal/helper"
	"bychat/internal/models"
	"bychat/lib/cache"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ClientManager 连接管理
type ClientManager struct {
	Clients     map[*Client]bool   // 全部的连接
	ClientsLock sync.RWMutex       // 读写锁
	Users       map[string]*Client // 登录的用户 // appID+uuid
	UserLock    sync.RWMutex       // 读写锁
	Register    chan *Client       // 连接连接处理
	BindUser    chan *Client       // 绑定用户信息
	Unregister  chan *Client       // 断开连接处理程序
	// Broadcast   chan []byte        // 广播 向全部成员发送数据
}

// NewClientManager 管理client
func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]*Client),
		Register:   make(chan *Client, 1000),
		BindUser:   make(chan *Client, 1000),
		Unregister: make(chan *Client, 1000),
		// Broadcast:  make(chan []byte, 1000),
	}

	return
}

/**************************  manager  ***************************************/

// InClient 校验
func (manager *ClientManager) InClient(client *Client) (ok bool) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	// 连接存在，在添加
	_, ok = manager.Clients[client]

	return
}

// GetClients 获取clients
func (manager *ClientManager) GetClients() (clients map[*Client]bool) {
	clients = make(map[*Client]bool)

	manager.ClientsRange(func(client *Client, value bool) (result bool) {
		clients[client] = value
		return true
	})
	return
}

// ClientsRange 遍历
func (manager *ClientManager) ClientsRange(f func(client *Client, value bool) (result bool)) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	for key, value := range manager.Clients {
		result := f(key, value)
		if result == false {
			return
		}
	}

	return
}

// GetClientsLen 获取client长度
func (manager *ClientManager) GetClientsLen() (clientsLen int) {
	clientsLen = len(manager.Clients)
	return
}

// AddClients 添加客户端
func (manager *ClientManager) AddClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	manager.Clients[client] = true
}

// DelClients 删除客户端
func (manager *ClientManager) DelClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	if _, ok := manager.Clients[client]; ok {
		delete(manager.Clients, client)
	}
}

// GetUserClient 获取用户的连接
func (manager *ClientManager) GetUserClient(appID, userID uint32) (client *Client) {
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	logrus.WithFields(logrus.Fields{
		"appID":  appID,
		"userID": userID,
	}).Info("ClientManager GetUserClient")
	key := fmt.Sprintf("%d_%d", appID, userID)
	if value, ok := manager.Users[key]; ok {
		client = value
	}
	return
}

// GetUsersLen 获取用户
func (manager *ClientManager) GetUsersLen() (userLen int) {
	userLen = len(manager.Users)
	return
}

// AddUsers 添加用户
func (manager *ClientManager) AddUsers(client *Client) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	key := client.GetKey()
	manager.Users[key] = client
}

// DelUsers 删除用户
func (manager *ClientManager) DelUsers(client *Client) (result bool) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	key := client.GetKey()
	if value, ok := manager.Users[key]; ok {
		// 判断是否为相同的用户
		if value.Addr != client.Addr {
			return
		}
		delete(manager.Users, key)
		result = true
	}

	return
}

// GetUserKeys 获取用户keys
func (manager *ClientManager) GetUserKeys() (userKeys []string) {
	userKeys = make([]string, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()

	for key := range manager.Users {
		userKeys = append(userKeys, key)
	}
	return
}

// GetUserClientList 获取用户客户端信息
func (manager *ClientManager) GetUserClientList() (clients []*Client) {
	clients = make([]*Client, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()

	for _, v := range manager.Users {
		clients = append(clients, v)
	}
	return
}

// 向全部成员(除了自己)发送数据
func (manager *ClientManager) sendAll(message []byte, roomID, userID uint32, ignoreClient *Client) {
	clients := manager.GetUserClientList()

	tmpRoomID := cache.GetChatRoomID(userID)
	logrus.WithFields(logrus.Fields{
		"roomID":        roomID,
		"userID":        userID,
		"tmpRoomID":     tmpRoomID,
		"ignoreClient":  ignoreClient.UserID,
		"ignoreClient1": ignoreClient.AppID,
	}).Info("sendAll 发送消息")
	for _, conn := range clients {
		if conn != ignoreClient && roomID == tmpRoomID {
			conn.SendMsg(message)
		}
	}
}

// EventRegister 用户建立连接事件
func (manager *ClientManager) EventRegister(client *Client) {
	manager.AddClients(client)
	logrus.Info("EventRegister 用户建立连接:", client.Addr)
}

// EventBindUser 用户信息绑定
func (manager *ClientManager) EventBindUser(client *Client) {

	manager.AddUsers(client)
	logrus.Info("EvenBindUser 用户信息绑定:", client.Addr, client)
}

// EventUnregister 用户断开连接
func (manager *ClientManager) EventUnregister(client *Client) {
	manager.DelClients(client)

	// 删除用户连接
	deleteResult := manager.DelUsers(client)
	if deleteResult == false {
		// 不是当前连接的客户端
		return
	}

	roomID := cache.GetChatRoomID(client.UserID)
	// 清除redis登录数据
	userOnline, err := cache.GetUserOnlineInfo(client.UserID)
	if err == nil && client.Addr == userOnline.Addr {
		userOnline.LogOut()
		cache.SetUserOnlineInfo(client.UserID, userOnline)
		cache.DelChatRoomUser(roomID, client.UserID)
	}

	// 关闭 chan
	close(client.Send)

	logrus.Info("EventUnregister 用户断开连接", client.Addr)

	if client.UserID != 0 {
		orderID := helper.GetOrderIDTime()
		// 更加用户ID查询房间id TODO
		SendUserMessageAll(client.AppID, roomID, client.UserID, orderID, models.MessageCmdExit, "用户已经离开~")
	}
	client.Socket.Close()
}

// bindChannel 注销client
func binUserdChannel(client *Client) {
	clientManager.BindUser <- client
}

// unregisterChannel 注销client
func unregisterChannel(client *Client) {
	clientManager.Unregister <- client
}

// 管道处理程序
func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.Register:
			// 建立连接事件
			manager.EventRegister(conn)

		case conn := <-manager.BindUser:
			// 绑定用户信息
			manager.EventBindUser(conn)

		case conn := <-manager.Unregister:
			// 断开连接事件
			manager.EventUnregister(conn)

			// case message := <-manager.Broadcast:
			// 	// 广播事件
			// 	clients := manager.GetClients()
			// 	for conn := range clients {
			// 		select {
			// 		case conn.Send <- message:
			// 		default:
			// 			close(conn.Send)
			// 		}
			// 	}
		}
	}
}

/**************************  manager info  ***************************************/

// GetManagerInfo 获取管理者信息
func GetManagerInfo(isDebug string) (managerInfo map[string]interface{}) {
	managerInfo = make(map[string]interface{})

	managerInfo["clientsLen"] = clientManager.GetClientsLen()        // 客户端连接数
	managerInfo["usersLen"] = clientManager.GetUsersLen()            // 登录用户数
	managerInfo["chanRegisterLen"] = len(clientManager.Register)     // 未处理连接事件数
	managerInfo["chanUnregisterLen"] = len(clientManager.Unregister) // 未处理退出登录事件数
	// managerInfo["chanBroadcastLen"] = len(clientManager.Broadcast)   // 未处理广播事件数

	if isDebug == "true" {
		addrList := make([]string, 0)
		clientManager.ClientsRange(func(client *Client, value bool) (result bool) {
			addrList = append(addrList, client.Addr)

			return true
		})

		users := clientManager.GetUserKeys()

		managerInfo["clients"] = addrList // 客户端列表
		managerInfo["users"] = users      // 登录用户列表
	}

	return
}

// ClearTimeoutConnections 定时清理超时连接
func ClearTimeoutConnections() {
	currentTime := uint64(time.Now().Unix())

	clients := clientManager.GetClients()
	for client := range clients {
		if client.IsHeartbeatTimeout(currentTime) {
			logrus.WithFields(logrus.Fields{
				"client.Addr":          client.Addr,
				"client.UserID":        client.UserID,
				"client.LoginTime":     client.LoginTime,
				"client.HeartbeatTime": client.HeartbeatTime,
			}).Info("心跳时间超时 关闭连接")
			unregisterChannel(client)
			// client.Socket.Close()
		}
	}
}
