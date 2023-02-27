package user

import (
	"bychat/api/v1/base"
	"bychat/internal/common"
	"bychat/internal/servers/websocket"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// List 查看全部在线用户
func List(c *gin.Context) {
	appIDStr := c.Query("appID")
	appIDUint64, _ := strconv.ParseInt(appIDStr, 10, 32)
	appID := uint32(appIDUint64)

	fmt.Println("http_request 查看全部在线用户", appID)

	data := make(map[string]interface{})

	userList := websocket.UserList(appID)
	data["userList"] = userList
	data["userCount"] = len(userList)

	base.Response(c, common.OK, "", data)
}
