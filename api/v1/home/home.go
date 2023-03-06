package home

import (
	"bychat/internal/servers/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Index 聊天页面
func Index(c *gin.Context) {
	logrus.Info("http_request 聊天首页")
	data := gin.H{
		"title":        "聊天首页",
		"appID":        websocket.GetDefaultAppID(),
		"httpUrl":      viper.GetString("app.httpUrl"),
		"webSocketUrl": viper.GetString("app.webSocketUrl"),
	}
	c.HTML(http.StatusOK, "index.html", data)
}
