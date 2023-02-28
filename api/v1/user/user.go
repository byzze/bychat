package user

import (
	"bychat/api/v1/base"
	"bychat/internal/common"
	"bychat/internal/models"
	"bychat/internal/servers/websocket"
	"bychat/lib/cache"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// List 查看全部在线用户
func List(ctx *gin.Context) {
	roomIDStr := ctx.Query("roomID")
	roomIDUint64, _ := strconv.ParseInt(roomIDStr, 10, 32)
	roomID := uint32(roomIDUint64)

	logrus.Info("http_request 查看全部在线用户 roomID:", roomID)

	data := make(map[string]interface{})

	userList := websocket.GetUserList(roomID)
	data["userList"] = userList
	data["userCount"] = len(userList)

	base.Response(ctx, common.OK, "", data)
}

// SendMessageAll 发送所有人消息
func SendMessageAll(ctx *gin.Context) {
	appIDStr := ctx.PostForm("appID")
	userID := ctx.PostForm("userID")
	msgID := ctx.PostForm("msgID")
	message := ctx.PostForm("message")

	appIDUint64, _ := strconv.ParseInt(appIDStr, 10, 32)

	//TODO
	appID := websocket.GetDefaultAppID()
	roomID := uint32(appIDUint64)

	logrus.WithFields(logrus.Fields{
		"roomID":  roomID,
		"userID":  userID,
		"msgID":   msgID,
		"message": message,
	}).Info("SendMessageAll")

	data := make(map[string]interface{})

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
	appIDStr := ctx.Query("appID")
	appIDUint64, _ := strconv.ParseInt(appIDStr, 10, 32)
	appID := uint32(appIDUint64)

	data := make(map[string]interface{})
	res, err := cache.ZGetMessageAll(appID)
	if err != nil {
		logrus.Error("HistoryMessageList", err)
		base.Response(ctx, common.OK, "", data)
		return
	}
	data["data"] = res
	base.Response(ctx, common.OK, "", data)
}
