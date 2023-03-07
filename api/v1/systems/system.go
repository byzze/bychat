package systems

import (
	"bychat/api/v1/base"
	"bychat/internal/common"
	"bychat/internal/websocket"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Status 查询系统状态
func Status(c *gin.Context) {
	isDebug := c.Query("isDebug")
	logrus.Info("http_request 查询系统状态", isDebug)

	data := make(map[string]interface{})

	numGoroutine := runtime.NumGoroutine()
	numCPU := runtime.NumCPU()

	// goroutine 的数量
	data["numGoroutine"] = numGoroutine
	data["numCPU"] = numCPU

	// ClientManager 信息
	data["managerInfo"] = websocket.GetManagerInfo(isDebug)

	base.Response(c, common.OK, "", data)
}
