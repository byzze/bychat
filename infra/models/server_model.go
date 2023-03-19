package models

import "fmt"

// ServerNode 服务器节点
type ServerNode struct {
	IP   string `json:"ip"`   // ip
	Port string `json:"port"` // 端口
}

// NewServerNode 新建
func NewServerNode(ip string, port string) *ServerNode {
	return &ServerNode{IP: ip, Port: port}
}

func (s *ServerNode) String() (str string) {
	if s == nil {
		return
	}

	str = fmt.Sprintf("%s:%s", s.IP, s.Port)
	return
}
