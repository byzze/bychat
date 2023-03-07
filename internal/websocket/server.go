package websocket

import (
	"bychat/internal/common"
	"bychat/internal/models"
	"bychat/lib/cache"
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

/**
	解析websocket操作指令执行的方法
**/

// BindUser 绑定用户信息
func BindUser(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK

	var request = &models.OpenRequest{}
	err := json.Unmarshal(message, request)
	if err != nil {
		code = common.ParameterIllegal
		logrus.WithField("err", err.Error()).Error("Open")
		return
	}
	userOnline, err := cache.GetUserOnlineInfo(request.UserID)
	if err != nil {
		code = common.ParameterIllegal
		logrus.WithField("err", err.Error()).Error("Open")
		return
	}
	client.Login(request.AppID, userOnline)
	err = cache.SetUserOnlineInfo(request.UserID, userOnline)
	if err != nil {
		code = common.ServerError
		logrus.WithFields(logrus.Fields{
			"seq": seq,
			"err": err,
		}).Error("webSocket_request SetUserOnlineInfo")
	}

	binUserdChannel(client)
	return
}

// Heartbeat 心跳
func Heartbeat(c *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	currentTime := uint64(time.Now().Unix())

	var request = &models.HeartBeatRequest{}
	err := json.Unmarshal(message, request)
	if err != nil {
		code = common.ParameterIllegal
		logrus.WithField("err", err.Error()).Error("webSocket_request Heartbeat")
		return
	}

	logrus.WithFields(logrus.Fields{
		"UserId": request.UserID,
	}).Info("webSocket_request Heartbeat")

	if !c.IsLogin() {
		logrus.WithFields(logrus.Fields{
			"UserId": request.UserID,
			"seq":    seq,
		}).Info("webSocket_request Heartbeat 用户未登录")
		code = common.NotLoggedIn

		return
	}

	userOnline, err := cache.GetUserOnlineInfo(request.UserID)
	if err != nil {
		if err == redis.Nil {
			code = common.NotLoggedIn
			logrus.WithFields(logrus.Fields{
				"seq":     seq,
				"c.AppID": c.AppID,
			}).Warn("webSocket_request Heartbeat 用户未登录")
		} else {
			code = common.ServerError
			logrus.WithFields(logrus.Fields{
				"seq":     seq,
				"c.AppID": c.AppID,
				"err":     err,
			}).Error("webSocket_request Heartbeat GetUserOnlineInfo")
		}
		return
	}

	c.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)

	err = cache.SetUserOnlineInfo(request.UserID, userOnline)
	if err != nil {
		code = common.ServerError
		logrus.WithFields(logrus.Fields{
			"seq":     seq,
			"c.AppID": c.AppID,
			"err":     err,
		}).Error("webSocket_request Heartbeat SetUserOnlineInfo")
	}
	return
}
