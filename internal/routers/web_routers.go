package routers

import (
	"bychat/api/fileserver"
	"bychat/api/home"
	"bychat/api/systems"
	"bychat/api/user"
	"net/http"

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
		userRouter.POST("/enter", user.EnterChatRoom)
		userRouter.POST("/exit", user.ExitChatRoom)
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
	// file
	file := router.Group("/upload")
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	{
		file.POST("/file", fileserver.UploadFile)
	}
	router.StaticFS("/fileserver/bychat", http.Dir("./api/v1/fileserver/bychat"))
}
