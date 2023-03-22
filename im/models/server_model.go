package models

import (
	"fmt"
	"sync"
)

// ServerNode 服务器节点
type ServerNode struct {
	IP   string `json:"ip"`   // ip
	Port string `json:"port"` // 端口
}

// ServerNodeInfo 服务节点信息
var ServerNodeInfo *ServerNode

var once sync.Once

// NewServerNode 新建
func NewServerNode(ip string, port string) *ServerNode {
	once.Do(func() {
		ServerNodeInfo = &ServerNode{IP: ip, Port: port}
	})
	return ServerNodeInfo
}

// IsLocal 校验本地
func IsLocal(server *ServerNode) (isLocal bool) {
	if server.IP == ServerNodeInfo.IP && server.Port == ServerNodeInfo.Port {
		isLocal = true
	}
	return
}

func (s *ServerNode) String() (str string) {
	if s == nil {
		return
	}

	str = fmt.Sprintf("%s:%s", s.IP, s.Port)
	return
}
