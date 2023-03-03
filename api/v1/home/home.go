package home

import (
	"bychat/internal/servers/websocket"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Index 聊天页面
func Index(c *gin.Context) {
	roomIDStr := c.Query("roomID")
	roomIDUint64, _ := strconv.ParseInt(roomIDStr, 10, 32)
	roomID := uint32(roomIDUint64)
	if !websocket.InRoomIDs(roomID) {
		roomID = websocket.GetDefaultRoomID()
	}

	logrus.Info("http_request 聊天首页", roomID)

	data := gin.H{
		"title": "聊天首页",
		// "roomID":       roomID,
		"appID":        websocket.GetDefaultAppID(),
		"httpUrl":      viper.GetString("app.httpUrl"),
		"webSocketUrl": viper.GetString("app.webSocketUrl"),
	}
	c.HTML(http.StatusOK, "index.html", data)
}
