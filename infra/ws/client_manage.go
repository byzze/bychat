package ws

import (
	"bychat/infra/models"
	"bychat/internal/cache"
	"bychat/pkg/utils"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	clientManager = NewClientManager() // 管理者
)

// ClientManager 连接管理
type ClientManager struct {
	Clients     map[*models.Client]bool   // 全部的连接
	ClientsLock sync.RWMutex              // 读写锁
	Users       map[string]*models.Client // 登录的用户 // appID+uuid
	UserLock    sync.RWMutex              // 读写锁
	Register    chan *models.Client       // 连接连接处理
	BindUser    chan *models.Client       // 绑定用户信息
	Unregister  chan *models.Client       // 断开连接处理程序
	// Broadcast   chan []byte        // 广播 向全部成员发送数据
}

var sOnce sync.Once

// NewClientManager 管理client
func NewClientManager() (clientManager *ClientManager) {
	sOnce.Do(func() {
		clientManager = &ClientManager{
			Clients:    make(map[*models.Client]bool),
			Users:      make(map[string]*models.Client),
			Register:   make(chan *models.Client, 1000),
			BindUser:   make(chan *models.Client, 1000),
			Unregister: make(chan *models.Client, 1000),
			// Broadcast:  make(chan []byte, 1000),
		}
	})
	return
}

// GetClientManager 获取clients
func GetClientManager() *ClientManager {
	return clientManager
}

/**************************  manager  ***************************************/

// InClient 校验
func (manager *ClientManager) InClient(client *models.Client) (ok bool) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	// 连接存在，在添加
	_, ok = manager.Clients[client]

	return
}

// GetClients 获取clients
func (manager *ClientManager) GetClients() (clients map[*models.Client]bool) {
	clients = make(map[*models.Client]bool)

	manager.ClientsRange(func(client *models.Client, value bool) (result bool) {
		clients[client] = value
		return true
	})
	return
}

// ClientsRange 遍历
func (manager *ClientManager) ClientsRange(f func(client *models.Client, value bool) (result bool)) {
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
func (manager *ClientManager) AddClients(client *models.Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	manager.Clients[client] = true
}

// DelClients 删除客户端
func (manager *ClientManager) DelClients(client *models.Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	if _, ok := manager.Clients[client]; ok {
		delete(manager.Clients, client)
	}
}

// GetUserClient 获取用户的连接
func (manager *ClientManager) GetUserClient(appID, userID uint32) (client *models.Client) {
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
func (manager *ClientManager) AddUsers(client *models.Client) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	key := client.GetKey()
	manager.Users[key] = client
}

// DelUsers 删除用户
func (manager *ClientManager) DelUsers(client *models.Client) (result bool) {
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
func (manager *ClientManager) GetUserClientList() (clients []*models.Client) {
	clients = make([]*models.Client, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()

	for _, v := range manager.Users {
		clients = append(clients, v)
	}
	return
}

// 向全部成员(除了自己)发送数据
func (manager *ClientManager) sendAll(message []byte, appID, roomID, userID uint32) {
	logrus.WithFields(logrus.Fields{
		"appID":  appID,
		"roomID": roomID,
		"userID": userID,
	}).Info("sendAll 发送消息")

	roomUserList := cache.GetChatRoomUser(roomID)
	for _, user := range roomUserList {
		conn := manager.GetUserClient(appID, user.ID)
		if conn != nil && user.ID != userID {
			conn.SendMsg(message)
		}
	}
}

// EventRegister 用户建立连接事件
func (manager *ClientManager) EventRegister(client *models.Client) {
	manager.AddClients(client)
	logrus.Info("EventRegister 用户建立连接:", client.Addr)
}

// EventBindUser 用户信息绑定
func (manager *ClientManager) EventBindUser(client *models.Client) {
	manager.AddUsers(client)
	logrus.Info("EvenBindUser 用户信息绑定:", client.Addr, client)
}

// EventUnregister 用户断开连接
func (manager *ClientManager) EventUnregister(client *models.Client) {
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
		orderID := utils.GetOrderIDTime()
		// 根据用户ID查询房间id TODO
		data := models.GetTextMsgDataExit(userOnline.NickName, "", orderID, "用户已经离开~")
		SendMsgAllServer(client.AppID, roomID, client.UserID, data)
	}
	client.Socket.Close()
}

// BinUserdChannel 注销client
func BinUserdChannel(client *models.Client) {
	GetClientManager().BindUser <- client
}

// UnregisterChannel 注销client
func UnregisterChannel(client *models.Client) {
	GetClientManager().Unregister <- client
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

	managerInfo["clientsLen"] = GetClientManager().GetClientsLen()        // 客户端连接数
	managerInfo["usersLen"] = GetClientManager().GetUsersLen()            // 登录用户数
	managerInfo["chanRegisterLen"] = len(GetClientManager().Register)     // 未处理连接事件数
	managerInfo["chanUnregisterLen"] = len(GetClientManager().Unregister) // 未处理退出登录事件数
	// managerInfo["chanBroadcastLen"] = len(clientManager.Broadcast)   // 未处理广播事件数

	if isDebug == "true" {
		addrList := make([]string, 0)
		GetClientManager().ClientsRange(func(client *models.Client, value bool) (result bool) {
			addrList = append(addrList, client.Addr)

			return true
		})

		users := GetClientManager().GetUserKeys()

		managerInfo["clients"] = addrList // 客户端列表
		managerInfo["users"] = users      // 登录用户列表
	}

	return
}

// IsLogin 是否登录了
func (manager *ClientManager) IsLogin(client *models.Client) (isLogin bool) {
	c := manager.GetUserClient(client.AppID, client.UserID)
	if c != nil {
		isLogin = true
	}
	return
}

// ClearTimeoutConnections 定时清理超时连接
func ClearTimeoutConnections() {
	currentTime := uint64(time.Now().Unix())

	clients := GetClientManager().GetClients()
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
