package api

import (
	"bychat/internal/common"
	"bychat/internal/models"
	"bychat/internal/websocket"
	"bychat/pkg/cache"
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

/**
	解析websocket操作指令执行的方法
**/

// BindUser 绑定用户信息
func BindUser(client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK

	var request = &models.OpenRequest{}
	err := json.Unmarshal(message, request)
	if err != nil {
		logrus.WithField("seq", seq).WithError(err).Error("BindUser: invalid request")
		return common.ParameterIllegal, "invalid request", nil
	}
	userOnline, err := cache.GetUserOnlineInfo(request.UserID)
	if err != nil {
		logrus.WithField("seq", seq).WithError(err).Error("BindUser: failed to get user online info")
		return common.ParameterIllegal, "failed to get user online info", nil
	}
	client.Login(request.AppID, userOnline)
	err = cache.SetUserOnlineInfo(request.UserID, userOnline)
	if err != nil {
		logrus.WithField("seq", seq).WithError(err).Error("BindUser: failed to set user online info")
		return common.ServerError, "failed to set user online info", nil
	}

	websocket.BinUserdChannel(client)
	return
}

// Heartbeat 心跳
func Heartbeat(client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
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

	if !client.IsLogin() {
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
				"c.AppID": client.AppID,
			}).Warn("webSocket_request Heartbeat 用户未登录")
		} else {
			code = common.ServerError
			logrus.WithFields(logrus.Fields{
				"seq":     seq,
				"c.AppID": client.AppID,
				"err":     err,
			}).Error("webSocket_request Heartbeat GetUserOnlineInfo")
		}
		return
	}

	client.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)

	err = cache.SetUserOnlineInfo(request.UserID, userOnline)
	if err != nil {
		code = common.ServerError
		logrus.WithFields(logrus.Fields{
			"seq":     seq,
			"c.AppID": client.AppID,
			"err":     err,
		}).Error("webSocket_request Heartbeat SetUserOnlineInfo")
	}
	return
}
