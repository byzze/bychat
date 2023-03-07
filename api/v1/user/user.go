package user

import (
	"bychat/api/v1/base"
	"bychat/internal/common"
	"bychat/internal/models"
	"bychat/internal/websocket"
	"bychat/lib/cache"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Param 参数
type Param struct {
	ID       uint32 `json:"id"`
	NickName string `json:"nickname"`
	MsgID    string `json:"msgID"`
	Message  string `json:"message"`
	AppID    uint32 `form:"appID" json:"appID"  binding:"-"`
	RoomID   uint32 `form:"roomID" json:"roomID"  binding:"-"`
	UserID   uint32 `form:"userID" json:"userID" `
	Start    int64  `form:"start"`
	Limit    int64  `form:"limit"`
}

// Login 登录
func Login(ctx *gin.Context) {
	data := make(map[string]interface{})
	var param Param
	if err := ctx.BindJSON(&param); err != nil {
		logrus.Error("websocket Login BindJSON:", err)
		base.Response(ctx, common.ParameterIllegal, "", data)
		return
	}

	err := websocket.Login(param.AppID, param.ID, param.NickName)
	if err != nil {
		logrus.Error("websocket Login", err)
		base.Response(ctx, common.OperationFailure, "", data)
		return
	}
	logrus.WithFields(logrus.Fields{
		"id":   param.ID,
		"name": param.NickName,
	}).Info("Login Info")

	// 放入缓存 map 后续可以redis
	data["token"] = fmt.Sprintf("%d", time.Now().Unix())
	base.Response(ctx, common.OK, "登陆成功", data)
}

// LogOut 退出
func LogOut(ctx *gin.Context) {
	var param Param
	if err := ctx.BindJSON(&param); err != nil {
		logrus.Error("websocket Login BindJSON:", err)
		base.Response(ctx, common.ParameterIllegal, "", nil)
		return
	}

	websocket.LogOut(param.AppID, param.ID)

	base.Response(ctx, common.OK, "退出成功", nil)
}

// GetRoomUserList 查看房间全部在线用户
func GetRoomUserList(ctx *gin.Context) {
	data := make(map[string]interface{})

	var param Param
	if err := ctx.BindQuery(&param); err != nil {
		logrus.Error("websocket GetRoomUserList:", err)
		base.Response(ctx, common.ParameterIllegal, "", nil)
		return
	}

	if param.RoomID == 0 {
		base.Response(ctx, common.ParameterIllegal, "", data)
		return
	}

	logrus.WithFields(logrus.Fields{
		"roomID": param.RoomID,
		"appID":  param.AppID,
	}).Info("http_request 查看全部在线用户 roomID:", param.RoomID)

	userList := websocket.GetRoomUserList(param.AppID, param.RoomID)

	data["userList"] = userList
	data["userCount"] = len(userList)

	base.Response(ctx, common.OK, "", data)
}

// SendMessageAll 发送所有人消息
func SendMessageAll(ctx *gin.Context) {
	data := make(map[string]interface{})

	var param Param
	if err := ctx.BindJSON(&param); err != nil {
		logrus.Error("websocket SendMessageAll:", err)
		base.Response(ctx, common.ParameterIllegal, "", nil)
		return
	}

	if param.RoomID == 0 {
		base.Response(ctx, common.ParameterIllegal, "", data)
		return
	}

	logrus.WithFields(logrus.Fields{
		"roomID":  param.RoomID,
		"userID":  param.UserID,
		"msgID":   param.MsgID,
		"message": param.Message,
	}).Info("SendMessageAll")

	if cache.SeqDuplicates(param.MsgID) {
		logrus.Info("数据重复：", param.MsgID)
		base.Response(ctx, common.OK, "", data)
		return
	}

	sendResults, err := websocket.SendUserMessageAll(param.AppID, param.RoomID, param.UserID, param.MsgID, models.MessageCmdMsg, param.Message)
	if err != nil {
		data["sendResultsErr"] = err.Error()
		base.Response(ctx, common.OperationFailure, err.Error(), data)
		return
	}

	data["sendResults"] = sendResults

	base.Response(ctx, common.OK, "", data)
}

// HistoryMessageList 获取聊天消息
func HistoryMessageList(ctx *gin.Context) {
	var param Param
	data := make(map[string]interface{})

	if err := ctx.BindQuery(&param); err != nil {
		base.Response(ctx, common.ParameterIllegal, "", data)
		return
	}

	if param.RoomID == 0 {
		base.Response(ctx, common.ParameterIllegal, "", data)
		return
	}

	res, err := cache.ZGetMessageByOffset(param.RoomID, param.Start, param.Limit)
	if err != nil {
		logrus.Error("ZGetMessageByOffset", err)
		base.Response(ctx, common.OperationFailure, "", data)
		return
	}

	data["data"] = res
	data["start"] = param.Start + param.Limit
	data["limit"] = param.Limit
	base.Response(ctx, common.OK, "", data)
}

// EnterChatRoom 进入房间
func EnterChatRoom(ctx *gin.Context) {
	var param Param

	if err := ctx.BindJSON(&param); err != nil {
		logrus.Error("EnterChatRoom Param Failed", err)
		base.Response(ctx, common.ParameterIllegal, err.Error(), nil)
		return
	}

	if param.RoomID == 0 {
		base.Response(ctx, common.ParameterIllegal, "", nil)
		return
	}

	err := websocket.EnterChatRoom(param.AppID, param.RoomID, param.UserID)
	if err != nil {
		logrus.Error("EnterChatRoom Failed", err)
		base.Response(ctx, common.OperationFailure, err.Error(), nil)
		return
	}

	base.Response(ctx, common.OK, "", nil)
}

// ExitChatRoom 离开房间
func ExitChatRoom(ctx *gin.Context) {
	var param Param

	if err := ctx.BindJSON(&param); err != nil {
		logrus.Error("EnterChatRoom Param Failed", err)
		base.Response(ctx, common.ParameterIllegal, err.Error(), nil)
		return
	}

	err := websocket.ExitChatRoom(param.AppID, param.RoomID, param.UserID)
	if err != nil {
		logrus.Error("ExitChatRoom Failed", err)
		base.Response(ctx, common.OperationFailure, err.Error(), nil)
		return
	}

	base.Response(ctx, common.OK, "", nil)
}
