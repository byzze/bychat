package cache

import (
	"bychat/lib/redislib"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	messageZSortKey = "acc:zset:message:roomid:" // 房间聊天数据
)

func getmessageZSortKey(roomID uint32) (key string) {
	key = fmt.Sprintf("%s%d", messageZSortKey, roomID)
	return
}

// ZSetMessage 设置数据
func ZSetMessage(roomID uint32, message string) (err error) {
	key := getmessageZSortKey(roomID)

	currentTime := float64(time.Now().Unix())

	redisClient := redislib.GetClient()
	_, err = redisClient.Do("zAdd", key, currentTime, message).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"key": key,
			"err": err,
		}).Error("ZSetMessage")
		return
	}
	return
}

// ZGetMessageAll 获取数据
func ZGetMessageAll(roomID uint32) (res []string, err error) {
	key := getmessageZSortKey(roomID)
	redisClient := redislib.GetClient()

	res, err = redisClient.ZRange(key, 0, -1).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"key": key,
			"err": err,
		}).Error("ZGetMessageAll")
	}
	return
}

// ZGetMessageByOffset 获取数据
func ZGetMessageByOffset(roomID uint32, start, limit int64) (res []string, err error) {
	key := getmessageZSortKey(roomID)
	redisClient := redislib.GetClient()

	res, err = redisClient.ZRange(key, start, limit).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"key": key,
			"err": err,
		}).Error("ZGetMessageAll")
	}
	return
}
