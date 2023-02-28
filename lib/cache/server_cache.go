package cache

import (
	"bychat/internal/models"
	"bychat/lib/redislib"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	serverNodesHashKey       = "acc:hash:server:nodes" // 全部的服务器
	serverNodesHashCacheTime = 2 * 60 * 60             // key过期时间
	serverNodesHashTimeout   = 3 * 60                  // 超时时间
)

func getServerNodesHashKey() (key string) {
	key = fmt.Sprintf("%s", serverNodesHashKey)
	return
}

// SetServerNodeInfo 设置服务器信息
func SetServerNodeInfo(serverNode *models.ServerNode, currentTime uint64) (err error) {
	key := getServerNodesHashKey()

	value := fmt.Sprintf("%d", currentTime)

	redisClient := redislib.GetClient()
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
	redisClient := redislib.GetClient()
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
func GetServerNodeAll(currentTime uint64) (servers []*models.ServerNode, err error) {
	servers = make([]*models.ServerNode, 0)
	key := getServerNodesHashKey()

	redisClient := redislib.GetClient()

	val, err := redisClient.Do("hGetAll", key).Result()

	valByte, _ := json.Marshal(val)

	logrus.WithFields(logrus.Fields{
		"key":     key,
		"valByte": string(valByte),
	}).Info("GetServerNodeAll")

	serverMap, err := redisClient.HGetAll(key).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"key": key,
			"err": err,
		}).Error("SetServerNodeInfo")
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
		if valueUint64+serverNodesHashTimeout <= currentTime {
			continue
		}

		server, err := models.StringToServer(key)
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
