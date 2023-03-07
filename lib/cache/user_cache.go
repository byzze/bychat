/**
 * Created by GoLand.
 * User: link1st
 * Date: 2019-07-25
 * Time: 17:28
 */

package cache

import (
	"bychat/internal/models"
	"bychat/lib/redislib"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const (
	userOnlinePrefix    = "acc:user:online:" // 用户在线状态
	userOnlineCacheTime = 24 * 60 * 60
)

/*********************  查询用户是否在线  ************************/
func getUserOnlineKey(userID uint32) (key string) {
	key = fmt.Sprintf("%s%d", userOnlinePrefix, userID)
	return
}

// GetUserOnlineInfo 用户在线信息
func GetUserOnlineInfo(userID uint32) (userOnline *models.UserOnline, err error) {
	redisClient := redislib.GetClient()

	key := getUserOnlineKey(userID)
	data, err := redisClient.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			logrus.WithFields(logrus.Fields{
				"userKey": userID,
				"err":     err,
			}).Info("GetUserOnlineInfo")
			return
		}
		logrus.WithFields(logrus.Fields{
			"userKey": userID,
			"err":     err,
		}).Error("GetUserOnlineInfo")
		return
	}

	userOnline = &models.UserOnline{}
	err = json.Unmarshal(data, userOnline)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userKey": userID,
			"err":     err,
		}).Error("GetUserOnlineInfo json Unmarshal")
		return
	}
	return
}

// SetUserOnlineInfo 设置用户在线数据
func SetUserOnlineInfo(userKey uint32, userOnline *models.UserOnline) (err error) {
	redisClient := redislib.GetClient()
	key := getUserOnlineKey(userKey)

	valueByte, err := json.Marshal(userOnline)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"key": key,
			"err": err,
		}).Error("SetUserOnlineInfo json Marshal Err")
		return
	}

	_, err = redisClient.Do("setEx", key, userOnlineCacheTime, string(valueByte)).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"key": key,
			"err": err,
		}).Error("SetUserOnlineInfo Err")
	}
	return
}
