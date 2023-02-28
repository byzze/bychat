package task

import (
	"bychat/internal/servers/websocket"
	"bychat/lib/cache"
	"fmt"
	"runtime/debug"
	"time"
)

// ServerNodeInit 服务器初始化
func ServerNodeInit() {
	Timer(2*time.Second, 60*time.Second, serverNode, "", serverDefer, "")
}

// 服务注册
func serverNode(param interface{}) (result bool) {
	result = true

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("服务注册 stop", r, string(debug.Stack()))
		}
	}()

	server := websocket.GetServerNode()
	currentTime := uint64(time.Now().Unix())
	fmt.Println("定时任务，服务注册", param, server, currentTime)

	cache.SetServerNodeInfo(server, currentTime)

	return
}

// 服务下线
func serverDefer(param interface{}) (result bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("服务下线 stop", r, string(debug.Stack()))
		}
	}()

	fmt.Println("服务下线", param)

	server := websocket.GetServerNode()
	cache.DelServerNodeInfo(server)

	return
}
