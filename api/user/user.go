package user

import (
	"bychat/api/base"
	"bychat/infra/models"
	"bychat/internal/api/user"
	"bychat/internal/cache"
	"bychat/pkg/common"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CommonParam 公共参数
type CommonParam struct {
	ID       uint32 `json:"id"`
	NickName string `json:"nickname"`
	AppID    uint32 `json:"appID" form:"appID"`
	RoomID   uint32 `json:"roomID" form:"roomID"`
	UserID   uint32 `json:"userID" form:"userID" `
}

// HistoryMessageParam 历史消息参数
type HistoryMessageParam struct {
	CommonParam
	Start int64 `form:"start"`
	Limit int64 `form:"limit"`
}

// SendMessageParam 发送消息参数
type SendMessageParam struct {
	CommonParam
	MsgSeq     string `json:"msgSeq"`
	MsgType    string `json:"msgType"`
	MsgContent string `json:"msgContent"`
	URL        string `json:"url"`
	Name       string `json:"name"`
	Format     string `json:"format"`
	Size       int64  `json:"size"`
	Second     int32  `json:"second"`
	Width      int32  `json:"width"`
	Height     int32  `json:"height"`
}

// Login 登录
func Login(ctx *gin.Context) {
	data := make(map[string]interface{})
	var param CommonParam
	if err := ctx.BindJSON(&param); err != nil {
		logrus.Error("websocket Login BindJSON:", err)
		base.Response(ctx, common.ParameterIllegal, "", data)
		return
	}

	err := user.Login(param.AppID, param.ID, param.NickName)
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
	var param CommonParam
	if err := ctx.BindJSON(&param); err != nil {
		logrus.Error("websocket Login BindJSON:", err)
		base.Response(ctx, common.ParameterIllegal, "", nil)
		return
	}

	user.LogOut(param.AppID, param.ID)

	base.Response(ctx, common.OK, "退出成功", nil)
}

// GetRoomUserList 查看房间全部在线用户
func GetRoomUserList(ctx *gin.Context) {
	data := make(map[string]interface{})

	var param CommonParam
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

	userList := user.GetRoomUserList(param.AppID, param.RoomID)

	data["userList"] = userList
	data["userCount"] = len(userList)

	base.Response(ctx, common.OK, "", data)
}

// SendMessageAll 发送所有人消息
func SendMessageAll(ctx *gin.Context) {
	data := make(map[string]interface{})

	var param SendMessageParam
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
		"roomID":         param.RoomID,
		"userID":         param.UserID,
		"msgSeq":         param.MsgSeq,
		"messageContetn": param.MsgContent,
		"messageType":    param.MsgType,
	}).Info("SendMessageAll Param")

	if cache.SeqDuplicates(param.MsgSeq) {
		logrus.Info("数据重复：", param.MsgSeq)
		base.Response(ctx, common.OK, "", data)
		return
	}
	uo, err := cache.GetUserOnlineInfo(param.UserID)
	if err != nil {
		logrus.Error("给全体用户发消息", err)
		return
	}

	var message string
	msgType := models.MessageType(param.MsgType)
	switch msgType {
	case models.MessageTypeText:
		message = models.GetTextMsgData(uo.NickName, "", param.MsgSeq, param.MsgContent)
	case models.MessageTypeImage:
		message = models.GetImgMsgData(uo.NickName, "", param.MsgSeq, param.URL, param.Name, param.Size, param.Width, param.Height)
	case models.MessageTypeFile:
		message = models.GetFileMsgData(uo.NickName, "", param.MsgSeq, param.URL, param.Name, param.Size)
	case models.MessageTypeVedio:
		message = models.GetVedioMsgData(uo.NickName, "", param.MsgSeq, param.URL, param.Name, param.Format, param.Size, param.Second)
	case models.MessageTypeAudio:
		message = models.GetSoundMsgData(uo.NickName, "", param.MsgSeq, param.URL, param.Size, param.Second)
	default:
		base.Response(ctx, common.ParameterIllegal, "未知数据格式", data)
		return
	}

	// 缓存聊天数据
	cache.ZSetMessage(param.RoomID, message)

	sendResults, err := user.UserSendMessageAll(param.AppID, param.RoomID, param.UserID, message)
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
	var param HistoryMessageParam

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
	var param CommonParam

	if err := ctx.BindJSON(&param); err != nil {
		logrus.Error("EnterChatRoom Param Failed", err)
		base.Response(ctx, common.ParameterIllegal, err.Error(), nil)
		return
	}

	if param.RoomID == 0 {
		base.Response(ctx, common.ParameterIllegal, "", nil)
		return
	}

	err := user.EnterChatRoom(param.AppID, param.RoomID, param.UserID)
	if err != nil {
		logrus.Error("EnterChatRoom Failed", err)
		base.Response(ctx, common.OperationFailure, err.Error(), nil)
		return
	}

	base.Response(ctx, common.OK, "", nil)
}

// ExitChatRoom 离开房间
func ExitChatRoom(ctx *gin.Context) {
	var param CommonParam

	if err := ctx.BindJSON(&param); err != nil {
		logrus.Error("EnterChatRoom Param Failed", err)
		base.Response(ctx, common.ParameterIllegal, err.Error(), nil)
		return
	}

	err := user.ExitChatRoom(param.AppID, param.RoomID, param.UserID)
	if err != nil {
		logrus.Error("ExitChatRoom Failed", err)
		base.Response(ctx, common.OperationFailure, err.Error(), nil)
		return
	}

	base.Response(ctx, common.OK, "", nil)
}
