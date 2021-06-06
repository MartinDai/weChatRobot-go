package routers

import (
	"github.com/gin-gonic/gin"
	"weChatRobot-go/controller"
	"weChatRobot-go/setting"
)

func SetupRouter() *gin.Engine {
	if setting.Conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	// 告诉gin框架模板文件引用的静态文件去哪里找
	r.Static("/static", "static")
	// 告诉gin框架去哪里找模板文件
	r.LoadHTMLGlob("templates/*")
	r.GET("/", controller.IndexHandler)

	weChatGroup := r.Group("weChat")
	{
		//签名回调
		weChatGroup.GET("/receiveMessage", controller.ReceiveMessage)
		//接收发送给公众号的消息
		weChatGroup.POST("/receiveMessage", controller.ReceiveMessage)
	}
	return r
}
