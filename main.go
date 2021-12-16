package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"weChatRobot-go/config"
	"weChatRobot-go/controller"
	"weChatRobot-go/service"
)

//go:embed static/images templates
var f embed.FS

//go:embed static/keyword.json
var keywordBytes []byte

func main() {
	var port int
	flag.IntVar(&port, "p", config.Port, "端口号")
	flag.Parse()

	go service.InitKeywordMap(keywordBytes)

	// 注册路由
	router := SetupRouter()
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		fmt.Printf("server startup failed, err:%v\n", err)
	}
}

func SetupRouter() *gin.Engine {
	if config.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	//模板文件
	templates := template.Must(template.New("").ParseFS(f, "templates/*.html"))
	router.SetHTMLTemplate(templates)
	//静态文件
	router.StaticFS("/public", http.FS(f))

	router.GET("/", controller.IndexHandler)

	weChatGroup := router.Group("weChat")
	{
		//签名回调
		weChatGroup.GET("/receiveMessage", controller.ReceiveMessage)
		//接收发送给公众号的消息
		weChatGroup.POST("/receiveMessage", controller.ReceiveMessage)
	}
	return router
}
