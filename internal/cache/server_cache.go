package cache

import (
	"bychat/infra/models"
	"bychat/infra/redislib"
	"fmt"

	"github.com/sirupsen/logrus"
)

// 服务节点信息
const (
	serverNodesHashKey       = "acc:hash:server:nodes" // 全部的服务器
	serverNodesHashCacheTime = 2 * 60 * 60             // key过期时间
	ServerNodesHashTimeout   = 3 * 60                  // 超时时间
)

var redisClient = redislib.GetClient()

func getServerNodesHashKey() (key string) {
	key = fmt.Sprintf("%s", serverNodesHashKey)
	return
}

// SetServerNodeInfo 设置服务器信息
func SetServerNodeInfo(serverNode *models.ServerNode, currentTime uint64) (err error) {
	key := getServerNodesHashKey()
	value := fmt.Sprintf("%d", currentTime)

	number, err := redisClient.Do("hSet", key, serverNode.String(), value).Int()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"key":    key,
			"number": number,
			"err":    err,
		}).Info("SetServerNodeInfo")
		return
	}

	if number != 1 {
		return
	}

	redisClient.Do("Expire", key, serverNodesHashCacheTime)

	return
}

// DelServerNodeInfo 下线服务器信息
func DelServerNodeInfo(serverNode *models.ServerNode) (err error) {
	key := getServerNodesHashKey()
	number, err := redisClient.Do("hDel", key, serverNode.String()).Int()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"key":    key,
			"number": number,
			"err":    err,
		}).Info("DelServerNodeInfo")
		return
	}

	if number != 1 {
		return
	}
	// 下线服务器，重新设置过期时间
	redisClient.Do("Expire", key, serverNodesHashCacheTime)
	return
}

// GetServerNodeAll 获取所有服务器
func GetServerNodeAll(currentTime uint64) (serverMap map[string]string, err error) {
	key := getServerNodesHashKey()
	serverMap, err = redisClient.HGetAll(key).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"key": key,
			"err": err,
		}).Error("SetServerNodeInfo")
	}
	return
}
