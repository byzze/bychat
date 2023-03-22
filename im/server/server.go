package server

import (
	"bychat/im/models"
	"errors"
	"strings"
)

// StringToServer 切割转换127.0.0.1:8080
func StringToServer(str string) (server *models.ServerNode, err error) {
	list := strings.Split(str, ":")
	if len(list) != 2 {
		return nil, errors.New("err")
	}

	server = &models.ServerNode{
		IP:   list[0],
		Port: list[1],
	}
	return
}
