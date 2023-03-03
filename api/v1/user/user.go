package user

import (
	"bychat/api/v1/base"
	"bychat/internal/common"
	"bychat/internal/models"
	"bychat/internal/servers/websocket"
	"bychat/lib/cache"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Param 参数
type Param struct {
	AppID  uint32 `form:"appID" json:"appID"  binding:"-"`
	RoomID uint32 `form:"roomID" json:"roomID"  binding:"required"`
	UserID string `form:"userID" json:"userID" `
	Start  int64  `form:"start"`
	Limit  int64  `form:"limit"`
}

// Login 登录
func Login(ctx *gin.Context) {
	data := make(map[string]interface{})
	id := ctx.PostForm("id")
	name := ctx.PostForm("name")

	websocket.Login(0, id, name)
	// 放入缓存 map 后续可以redis
	data["token"] = fmt.Sprintf("%d", time.Now().Unix())
	base.Response(ctx, common.OK, "登陆成功", data)
}

// LogOut 退出
func LogOut(ctx *gin.Context) {
	data := make(map[string]interface{})
	id := ctx.PostForm("id")

	websocket.LogOut(0, id)

	base.Response(ctx, common.OK, "退出成功", data)
}

// GetRoomUserList 查看全部在线用户
func GetRoomUserList(ctx *gin.Context) {
	data := make(map[string]interface{})

	roomIDStr := ctx.Query("roomID")
	roomIDUint64, _ := strconv.ParseInt(roomIDStr, 10, 32)
	roomID := uint32(roomIDUint64)
	appID := websocket.GetDefaultAppID()

	if roomID == 0 {
		base.Response(ctx, common.OK, "参数错误", data)
		return
	}

	logrus.WithFields(logrus.Fields{
		"roomID": roomID,
		"appID":  appID,
	}).Info("http_request 查看全部在线用户 roomID:", roomID)

	userList := websocket.GetRoomUserList(appID, roomID)
	data["userList"] = userList
	data["userCount"] = len(userList)

	base.Response(ctx, common.OK, "", data)
}

// SendMessageAll 发送所有人消息
func SendMessageAll(ctx *gin.Context) {
	data := make(map[string]interface{})

	roomIDStr := ctx.PostForm("roomID")
	userID := ctx.PostForm("userID")
	msgID := ctx.PostForm("msgID")
	message := ctx.PostForm("message")

	roomIDUint64, _ := strconv.ParseInt(roomIDStr, 10, 32)
	roomID := uint32(roomIDUint64)
	appID := websocket.GetDefaultAppID()

	if roomID == 0 {
		base.Response(ctx, common.OK, "参数错误", data)
		return
	}

	logrus.WithFields(logrus.Fields{
		"roomID":  roomID,
		"userID":  userID,
		"msgID":   msgID,
		"message": message,
	}).Info("SendMessageAll")

	if cache.SeqDuplicates(msgID) {
		logrus.Info("数据重复：", msgID)
		base.Response(ctx, common.OK, "", data)
		return
	}

	sendResults, err := websocket.SendUserMessageAll(appID, roomID, userID, msgID, models.MessageCmdMsg, message)
	if err != nil {
		data["sendResultsErr"] = err.Error()
	}

	data["sendResults"] = sendResults

	base.Response(ctx, common.OK, "", data)
}

// HistoryMessageList 获取聊天消息
func HistoryMessageList(ctx *gin.Context) {
	var param Param
	data := make(map[string]interface{})

	if err := ctx.ShouldBindQuery(&param); err != nil {
		base.Response(ctx, common.OK, "", data)
		return
	}

	if param.RoomID == 0 {
		base.Response(ctx, common.OK, "参数错误", data)
		return
	}

	res, err := cache.ZGetMessageByOffset(param.RoomID, param.Start, param.Limit)
	if err != nil {
		logrus.Error("ZGetMessageByOffset", err)
		base.Response(ctx, common.OK, "", data)
		return
	}

	data["data"] = res
	data["start"] = param.Start + param.Limit
	data["limit"] = param.Limit
	base.Response(ctx, common.OK, "", data)
}

// EnterRoom 进入房间
func EnterRoom(ctx *gin.Context) {
	var param Param
	data := make(map[string]interface{})

	if err := ctx.Bind(&param); err != nil {
		data["err"] = err.Error()
		logrus.Error("EnterRoom Param Failed", err)
		base.Response(ctx, common.OK, "参数错误", data)
		return
	}

	if param.RoomID == 0 {
		base.Response(ctx, common.OK, "参数错误", data)
		return
	}

	websocket.EnterRoom(param.AppID, param.RoomID, param.UserID)

	base.Response(ctx, common.OK, "", data)
}

// ExitRoom 离开房间
func ExitRoom(ctx *gin.Context) {
	var param Param
	data := make(map[string]interface{})

	if err := ctx.Bind(&param); err != nil {
		data["err"] = err.Error()
		logrus.Error("EnterRoom Param Failed", err)
		base.Response(ctx, common.OK, "参数错误", data)
		return
	}

	websocket.ExitRoom(param.AppID, param.UserID)

	base.Response(ctx, common.OK, "", data)
}
