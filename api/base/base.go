package base

import (
	"bychat/pkg/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Base 基础路由
type Base struct {
	gin.Context
}

// Response 获取全部请求解析到map
func Response(c *gin.Context, code uint32, msg string, data interface{}) {
	message := common.Response(code, msg, data)
	/*
		这段代码的作用是对服务器的所有 HTTP 响应做如下处理：

		设置 "Access-Control-Allow-Origin" 头，告诉浏览器允许任意域访问这个服务器。
		设置 "Access-Control-Allow-Methods" 头，告诉浏览器服务器支持的所有 CORS 请求的方法。
		设置 "Access-Control-Allow-Headers" 头，告诉浏览器允许服务器接收哪些请求头。
		设置 "Access-Control-Expose-Headers" 头，告诉浏览器允许浏览器解析的响应头。
		设置 "Access-Control-Allow-Credentials" 头，告诉浏览器服务器是否允许返回的数据包含 cookie 信息。
		设置 "Content-Type" 头为 "application/json"，告诉浏览器服务器返回的数据格式为 JSON。
	*/
	// 允许跨域
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Origin", "*")                                        // 这是允许访问所有域
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE") // 服务器支持的所有跨域请求的方法,为了避免浏览器请求的多次'预检'请求
	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
	c.Header("Access-Control-Allow-Credentials", "true")                                                                                                                                                   //  跨域请求是否需要带cookie信息 默认设置为true
	c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json

	c.JSON(http.StatusOK, message)
	return
}
