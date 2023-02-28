package models

import (
	"errors"
	"fmt"
	"strings"
)

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

// StringToServer 切割转换127.0.0.1:8080
func StringToServer(str string) (server *ServerNode, err error) {
	list := strings.Split(str, ":")
	if len(list) != 2 {
		return nil, errors.New("err")
	}

	server = &ServerNode{
		IP:   list[0],
		Port: list[1],
	}
	return
}
