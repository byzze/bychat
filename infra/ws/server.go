package ws

import (
	"bychat/infra/models"
	"bychat/internal/cache"
	"errors"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
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

// GetServerNodeAll 获取全部节点
func GetServerNodeAll(currentTime uint64) (servers []*models.ServerNode, err error) {
	servers = make([]*models.ServerNode, 0)
	serverMap, err := cache.GetServerNodeAll(currentTime)
	if err != nil {
		return
	}

	for key, value := range serverMap {
		valueUint64, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"key": key,
				"err": err,
			}).Error("GetServerNodeAll")
			return nil, err
		}

		// 超时
		if valueUint64+cache.ServerNodesHashTimeout <= currentTime {
			continue
		}

		server, err := StringToServer(key)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"key": key,
				"err": err,
			}).Error("GetServerNodeAll")
			return nil, err
		}
		servers = append(servers, server)
	}
	return
}
