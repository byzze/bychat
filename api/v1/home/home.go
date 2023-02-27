package home

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 聊天页面
func Index(c *gin.Context) {

	appIDStr := c.Query("appID")
	appIDUint64, _ := strconv.ParseInt(appIDStr, 10, 32)
	appID := uint32(appIDUint64)
	// if !websocket.InAppIds(appId) {
	// 	appId = websocket.GetDefaultAppId()
	// }

	fmt.Println("http_request 聊天首页", appID)

	data := gin.H{
		"title":        "聊天首页",
		"appId":        appID,
		"httpUrl":      viper.GetString("app.httpUrl"),
		"webSocketUrl": viper.GetString("app.webSocketUrl"),
	}
	c.HTML(http.StatusOK, "index.tpl", data)
}
