package models

import (
	"time"

	"github.com/sirupsen/logrus"
)

const (
	heartbeatTimeout = 3 * 60 // 用户心跳超时时间
)

// UserOnline 用户在线状态
type UserOnline struct {
	ID            uint32 `json:"id"`            // 用户id
	NickName      string `json:"name"`          // 用户name
	LoginTime     uint64 `json:"loginTime"`     // 用户上次登录时间
	HeartbeatTime uint64 `json:"heartbeatTime"` // 用户上次心跳时间
	LogOutTime    uint64 `json:"logOutTime"`    // 用户退出登录的时间
	DeviceInfo    string `json:"deviceInfo"`    // 设备信息
	IsLogoff      bool   `json:"isLogoff"`      // 是否下线
}

// type RoomInfo struct {
// 	ID     string        `json:"id"`
// 	Name   string        `json:"name"`
// 	People []*UserOnline `json:"people"`
// }

/**********************  数据处理  *********************************/

// UserLogin 用户登录
// func UserLogin(accIP, accPort string, appID uint32, nickName string, addr string, loginTime uint64) (userOnline *UserOnline) {
// 	userOnline = &UserOnline{
// 		NickName:      nickName,
// 		LoginTime:     loginTime,
// 		HeartbeatTime: loginTime,
// 		IsLogoff:      false,
// 	}
// 	return
// }

// Heartbeat 用户心跳
func (u *UserOnline) Heartbeat(currentTime uint64) {
	u.HeartbeatTime = currentTime
	u.IsLogoff = false

	return
}

// LogOut 用户退出登录
func (u *UserOnline) LogOut() {
	currentTime := uint64(time.Now().Unix())
	u.LogOutTime = currentTime
	u.IsLogoff = true

	return
}

/**********************  数据操作  *********************************/

// IsOnline 用户是否在线
func (u *UserOnline) IsOnline() (online bool) {
	if u.IsLogoff {
		return
	}

	currentTime := uint64(time.Now().Unix())

	if u.HeartbeatTime < (currentTime - heartbeatTimeout) {
		logrus.WithFields(logrus.Fields{
			"userID":        u.ID,
			"heartbeatTime": u.HeartbeatTime,
		}).Info("用户是否在线：心跳超时")
		return
	}

	if u.IsLogoff {
		logrus.WithFields(logrus.Fields{
			"userID":        u.ID,
			"heartbeatTime": u.HeartbeatTime,
		}).Info("用户是否在线 用户已经下线")
		return
	}

	return true
}
