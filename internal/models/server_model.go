package models

import (
	"errors"
	"fmt"
	"strings"
)

type Server struct {
	IP   string `json:"ip"`   // ip
	Port string `json:"port"` // 端口
}

// NewServer 新建
func NewServer(ip string, port string) *Server {
	return &Server{IP: ip, Port: port}
}

func (s *Server) String() (str string) {
	if s == nil {
		return
	}

	str = fmt.Sprintf("%s:%s", s.IP, s.Port)
	return
}

// StringToServer 切割转换127.0.0.1:8080
func StringToServer(str string) (server *Server, err error) {
	list := strings.Split(str, ":")
	if len(list) != 2 {
		return nil, errors.New("err")
	}

	server = &Server{
		IP:   list[0],
		Port: list[1],
	}
	return
}
