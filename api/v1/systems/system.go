package systems

import (
	"bychat/api/v1/base"
	"bychat/internal/common"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 查询系统状态
func Status(c *gin.Context) {

	isDebug := c.Query("isDebug")
	logrus.Info("http_request 查询系统状态", isDebug)

	data := make(map[string]interface{})

	numGoroutine := runtime.NumGoroutine()
	numCPU := runtime.NumCPU()

	// goroutine 的数量
	data["numGoroutine"] = numGoroutine
	data["numCPU"] = numCPU
	base.Response(c, common.OK, "", data)
}
