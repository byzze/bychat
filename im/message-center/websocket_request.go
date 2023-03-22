package messagecenter

import (
	"bychat/im/cache"
	"bychat/im/client"
	"bychat/im/models"
	"bychat/pkg/common"
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

// Login websocket处理数据
func Login(c *client.Client, msgSeq string, message []byte) (code uint32, msg string, data interface{}) {
	return
}

// BindUser websocket处理数据
func BindUser(c *client.Client, msgSeq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK

	var request = &models.OpenRequest{}
	err := json.Unmarshal(message, request)
	if err != nil {
		logrus.WithField("msgSeq", msgSeq).WithError(err).Error("BindUser: invalid request")
		return common.ParameterIllegal, "invalid request", nil
	}

	userOnline, err := cache.GetUserOnlineInfo(request.UserID)
	if err != nil {
		logrus.WithField("msgSeq", msgSeq).WithError(err).Error("BindUser: failed to get user online info")
		return common.ParameterIllegal, "failed to get user online info", nil
	}

	c.Login(request.AppID, request.UserID)

	err = cache.SetUserOnlineInfo(request.UserID, userOnline)
	if err != nil {
		logrus.WithField("msgSeq", msgSeq).WithError(err).Error("BindUser: failed to set user online info")
		return common.ServerError, "failed to set user online info", nil
	}

	client.BinUserdChannel(c)
	return
}

/* // Logout websocket处理数据
func Logout(c *client.Client, msgSeq string, message []byte) (code uint32, msg string, data interface{}) {
	return
} */

// MsgProcess websocket处理数据
func MsgProcess(c *client.Client, msgSeq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	msg = ""
	return
}

// Heartbeat 心跳
func Heartbeat(c *client.Client, msgSeq string, message []byte) (code uint32, msg string, data interface{}) {
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
		"userId": request.UserID,
		"appId":  request.AppID,
	}).Info("webSocket_request Heartbeat")

	if !client.IsLogin(c) {
		logrus.WithFields(logrus.Fields{
			"UserId": request.UserID,
		}).Info("webSocket_request Heartbeat 用户未登录")
		code = common.NotLoggedIn
		return
	}

	userOnline, err := cache.GetUserOnlineInfo(request.UserID)
	if err != nil {
		if err == redis.Nil {
			code = common.NotLoggedIn
			logrus.WithFields(logrus.Fields{
				"userId": request.UserID,
				"appId":  request.AppID,
			}).Warn("webSocket_request Heartbeat 用户未登录")
		} else {
			code = common.ServerError
			logrus.WithFields(logrus.Fields{
				"userId": request.UserID,
				"appId":  request.AppID,
				"err":    err,
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
			"seq":     msgSeq,
			"c.AppID": c.AppID,
			"err":     err,
		}).Error("webSocket_request Heartbeat SetUserOnlineInfo")
	}
	return
}
