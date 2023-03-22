package client

import (
	"bychat/im/cache"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	manager = NewManager() // 管理者
)

// Manager 连接管理
type Manager struct {
	Clients     map[*Client]bool   // 全部的连接
	ClientsLock sync.RWMutex       // 读写锁
	Users       map[string]*Client // 登录的用户 // appID+uuid
	UserLock    sync.RWMutex       // 读写锁
	Register    chan *Client       // 连接连接处理
	BindUser    chan *Client       // 绑定用户信息
	Unregister  chan *Client       // 断开连接处理程序
	// Broadcast   chan []byte        // 广播 向全部成员发送数据
}

var sOnce sync.Once

// NewManager 管理client
func NewManager() (manager *Manager) {
	sOnce.Do(func() {
		manager = &Manager{
			Clients:    make(map[*Client]bool),
			Users:      make(map[string]*Client),
			Register:   make(chan *Client, 1000),
			BindUser:   make(chan *Client, 1000),
			Unregister: make(chan *Client, 1000),
			// Broadcast:  make(chan []byte, 1000),
		}
	})
	return
}

// getManager 获取clients
func getManager() *Manager {
	return manager
}

/**************************  manager  ***************************************/

// inClient 校验
func (manager *Manager) inClient(client *Client) (ok bool) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	// 连接存在，在添加
	_, ok = manager.Clients[client]

	return
}

// getClients 获取clients
func (manager *Manager) getClients() (clients map[*Client]bool) {
	clients = make(map[*Client]bool)

	manager.clientsRange(func(client *Client, value bool) (result bool) {
		clients[client] = value
		return true
	})
	return
}

// clientsRange 遍历
func (manager *Manager) clientsRange(f func(client *Client, value bool) (result bool)) {
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

// getClientsLen 获取client长度
func (manager *Manager) getClientsLen() (clientsLen int) {
	clientsLen = len(manager.Clients)
	return
}

// addClients 添加客户端
func (manager *Manager) addClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	manager.Clients[client] = true
}

// delClients 删除客户端
func (manager *Manager) delClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	if _, ok := manager.Clients[client]; ok {
		delete(manager.Clients, client)
	}
}

// getUserClient 获取用户的连接
func (manager *Manager) getUserClient(appID, userID uint32) (client *Client) {
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	logrus.WithFields(logrus.Fields{
		"appID":  appID,
		"userID": userID,
	}).Info("Manager GetUserClient")
	key := fmt.Sprintf("%d_%d", appID, userID)
	if value, ok := manager.Users[key]; ok {
		client = value
	}
	return
}

// getUsersLen 获取用户
func (manager *Manager) getUsersLen() (userLen int) {
	userLen = len(manager.Users)
	return
}

// addUsers 添加用户
func (manager *Manager) addUsers(client *Client) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	key := client.GetKey()
	manager.Users[key] = client
}

// delUsers 删除用户
func (manager *Manager) delUsers(client *Client) (result bool) {
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

// getUserKeys 获取用户keys
func (manager *Manager) getUserKeys() (userKeys []string) {
	userKeys = make([]string, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()

	for key := range manager.Users {
		userKeys = append(userKeys, key)
	}
	return
}

// getUserClientList 获取用户客户端信息
func (manager *Manager) getUserClientList() (clients []*Client) {
	clients = make([]*Client, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()

	for _, v := range manager.Users {
		clients = append(clients, v)
	}
	return
}

// eventRegister 用户建立连接事件
func (manager *Manager) eventRegister(client *Client) {
	manager.addClients(client)
	logrus.Info("EventRegister 用户建立连接:", client.Addr)
}

// eventBindUser 用户信息绑定
func (manager *Manager) eventBindUser(client *Client) {
	manager.addUsers(client)
	logrus.Info("EvenBindUser 用户信息绑定:", client.Addr, client)
}

// eventUnregister 用户断开连接
func (manager *Manager) eventUnregister(client *Client) {
	manager.delClients(client)
	// 删除用户连接
	deleteResult := manager.delUsers(client)
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
		/* 		orderID := utils.GetOrderIDTime()
		   		// 根据用户ID查询房间id TODO
		   		data := GetTextMsgDataExit(userOnline.NickName, "", orderID, "用户已经离开~")
		   		SendMsgAllServer(client.AppID, roomID, client.UserID, data) */
	}
	client.Socket.Close()
}

// 管道处理程序
func (manager *Manager) start() {
	for {
		select {
		case conn := <-manager.Register:
			// 建立连接事件
			manager.eventRegister(conn)

		case conn := <-manager.BindUser:
			// 绑定用户信息
			manager.eventBindUser(conn)

		case conn := <-manager.Unregister:
			// 断开连接事件
			manager.eventUnregister(conn)

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

// BinUserdChannel 注销client
func BinUserdChannel(client *Client) {
	getManager().BindUser <- client
}

// UnregisterChannel 注销client
func UnregisterChannel(client *Client) {
	getManager().Unregister <- client
}
