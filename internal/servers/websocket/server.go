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

// Login 登陆
func Login(c *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	currentTime := uint64(time.Now().Unix())

	var request = &models.LoginRequest{}
	err := json.Unmarshal(message, request)
	if err != nil {
		code = common.ParameterIllegal
		logrus.WithField("err", err.Error()).Error("Login")
		return
	}
	if c.IsLogin() {
		logrus.WithFields(logrus.Fields{
			"client.AppId":  request.AppID,
			"client.UserId": request.UserID,
			"seq":           seq,
		}).Info("用户登录 用户已经登录")
		code = common.OperationFailure
		return
	}
	userOnline := models.UserLogin(serverIP, serverPort, request.AppID, request.UserID, c.Addr, currentTime)

	c.UserOnline = userOnline

	clientManager.Login <- c

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
		logrus.WithField("err", err.Error()).Error("Heartbeat")
		return
	}

	logrus.WithFields(logrus.Fields{
		"UserId": request.UserID,
	}).Info("webSocket_request 心跳接口")

	if !c.IsLogin() {
		logrus.WithFields(logrus.Fields{
			"UserId": request.UserID,
			"seq":    seq,
		}).Info("心跳接口 用户未登录")
		code = common.NotLoggedIn

		return
	}

	userOnline, err := cache.GetUserOnlineInfo(c.GetKey())
	if err != nil {
		if err == redis.Nil {
			code = common.NotLoggedIn
			logrus.WithFields(logrus.Fields{
				"seq":          seq,
				"c.AppID":      c.AppID,
				"c.UserOnline": c.UserOnline,
			}).Warn("心跳接口 用户未登录")
		} else {
			code = common.ServerError
			logrus.WithFields(logrus.Fields{
				"seq":          seq,
				"c.AppID":      c.AppID,
				"c.UserOnline": c.UserOnline,
				"err":          err,
			}).Error("心跳接口 GetUserOnlineInfo")
		}
		return
	}

	c.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)
	err = cache.SetUserOnlineInfo(c.GetKey(), userOnline)
	if err != nil {
		code = common.ServerError
		logrus.WithFields(logrus.Fields{
			"seq":          seq,
			"c.AppID":      c.AppID,
			"c.UserOnline": c.UserOnline,
			"err":          err,
		}).Error("心跳接口 SetUserOnlineInfo")
	}
	return
}
