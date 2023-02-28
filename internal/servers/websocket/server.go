package websocket

import (
	"bychat/internal/common"
	"bychat/internal/models"
	"bychat/lib/cache"
	"encoding/json"
	"fmt"
	"time"

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
			"client.AppId":  c.AppID,
			"client.UserId": c.UserID,
			"seq":           seq,
		}).Info("用户登录 用户已经登录")
		code = common.OperationFailure
		return
	}

	c.Login(request.AppID, request.UserID, currentTime)

	// 存储redis数据
	userOnline := models.UserLogin(serverIP, serverPort, request.AppID, request.UserID, c.Addr, currentTime)
	err = cache.SetUserOnlineInfo(c.GetKey(), userOnline)
	if err != nil {
		code = common.ServerError
		fmt.Println("用户登录 SetUserOnlineInfo", seq, err)
		return
	}

	// 用户登录
	login := &login{
		AppID:  request.AppID,
		UserID: request.UserID,
		Client: c,
	}
	clientManager.Login <- login
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
		"AppId":  c.AppID,
		"UserId": c.UserID,
		"RoomID": c.RoomID,
	}).Info("webSocket_request 心跳接口")

	if !c.IsLogin() {
		logrus.WithFields(logrus.Fields{
			"AppId":  c.AppID,
			"UserId": c.UserID,
			"RoomID": c.RoomID,
			"seq":    seq,
		}).Info("心跳接口 用户未登录")
		code = common.NotLoggedIn

		return
	}

	c.Heartbeat(currentTime)
	// todo cache
	return
}
