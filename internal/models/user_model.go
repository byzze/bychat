package models

import (
	"fmt"
	"time"
)

const (
	heartbeatTimeout = 3 * 60 // 用户心跳超时时间
)

// 用户在线状态
type UserOnline struct {
	AccIP         string `json:"accIp"`         // acc Ip
	AccPort       string `json:"accPort"`       // acc 端口
	AppID         uint32 `json:"appID"`         // appID
	RoomID        uint32 `json:"roomID"`        // roomID
	UserID        string `json:"userID"`        // 用户Id
	ClientIP      string `json:"clientIp"`      // 客户端Ip
	ClientPort    string `json:"clientPort"`    // 客户端端口
	LoginTime     uint64 `json:"loginTime"`     // 用户上次登录时间
	HeartbeatTime uint64 `json:"heartbeatTime"` // 用户上次心跳时间
	LogOutTime    uint64 `json:"logOutTime"`    // 用户退出登录的时间
	// Qua           string `json:"qua"`           // qua
	DeviceInfo string `json:"deviceInfo"` // 设备信息
	IsLogoff   bool   `json:"isLogoff"`   // 是否下线
}

/**********************  数据处理  *********************************/

// 用户登录
func UserLogin(accIP, accPort string, appID uint32, userID string, addr string, loginTime uint64) (userOnline *UserOnline) {

	userOnline = &UserOnline{
		AccIP:         accIP,
		AccPort:       accPort,
		AppID:         appID,
		UserID:        userID,
		ClientIP:      addr,
		LoginTime:     loginTime,
		HeartbeatTime: loginTime,
		IsLogoff:      false,
	}

	return
}

// 用户心跳
func (u *UserOnline) Heartbeat(currentTime uint64) {

	u.HeartbeatTime = currentTime
	u.IsLogoff = false

	return
}

// 用户退出登录
func (u *UserOnline) LogOut() {

	currentTime := uint64(time.Now().Unix())
	u.LogOutTime = currentTime
	u.IsLogoff = true

	return
}

/**********************  数据操作  *********************************/

// 用户是否在线
func (u *UserOnline) IsOnline() (online bool) {
	if u.IsLogoff {

		return
	}

	currentTime := uint64(time.Now().Unix())

	if u.HeartbeatTime < (currentTime - heartbeatTimeout) {
		fmt.Println("用户是否在线 心跳超时", u.AppID, u.UserID, u.HeartbeatTime)

		return
	}

	if u.IsLogoff {
		fmt.Println("用户是否在线 用户已经下线", u.AppID, u.UserID)

		return
	}

	return true
}

// 用户是否在本台机器上
func (u *UserOnline) UserIsLocal(localIP, localPort string) (result bool) {

	if u.AccIP == localIP && u.AccPort == localPort {
		result = true

		return
	}

	return
}
