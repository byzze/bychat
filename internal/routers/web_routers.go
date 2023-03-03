package routers

import (
	"bychat/api/v1/home"
	"bychat/api/v1/systems"
	"bychat/api/v1/user"

	"github.com/gin-gonic/gin"
)

// InitWeb 初始化路由
func InitWeb(router *gin.Engine) {
	router.LoadHTMLGlob("web/**/*")

	// 用户组
	userRouter := router.Group("/user")
	{
		userRouter.POST("/login", user.Login)
		userRouter.POST("/logout", user.LogOut)
		userRouter.GET("/list", user.GetRoomUserList)
		userRouter.POST("/enter", user.EnterRoom)
		userRouter.POST("/exit", user.ExitRoom)
		// userRouter.GET("/online", user.Online)
		// userRouter.POST("/sendMessage", user.SendMessage)
		userRouter.POST("/sendMessageAll", user.SendMessageAll)
		userRouter.GET("/historyMessageList", user.HistoryMessageList)
	}

	// 系统
	systemRouter := router.Group("/system")
	{
		systemRouter.GET("/state", systems.Status)
	}

	// home
	homeRouter := router.Group("/home")
	{
		homeRouter.GET("/index", home.Index)
	}
}
