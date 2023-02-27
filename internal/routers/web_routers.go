package routers

import (
	"bychat/api/v1/home"
	"bychat/api/v1/systems"
	"bychat/api/v1/user"

	"github.com/gin-gonic/gin"
)

// 初始化路由
func Init(router *gin.Engine) {
	router.LoadHTMLGlob("web/**/*")

	// 用户组
	userRouter := router.Group("/user")
	{
		userRouter.GET("/list", user.List)
		// userRouter.GET("/online", user.Online)
		// userRouter.POST("/sendMessage", user.SendMessage)
		// userRouter.POST("/sendMessageAll", user.SendMessageAll)
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
