package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"weChatRobot-go/chatgpt"
	"weChatRobot-go/config"
	"weChatRobot-go/controller"
	"weChatRobot-go/models"
	"weChatRobot-go/service"
	"weChatRobot-go/tuling"

	"github.com/gin-gonic/gin"
)

//go:embed static/images templates
var fs embed.FS

//go:embed static/keyword.json
var keywordBytes []byte

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "配置文件")
	flag.Parse()

	if err := runApp(configFile); err != nil {
		log.Fatal(err)
	}
}

func runApp(configFile string) error {
	configSettings, err := getConfigSettings(configFile)
	if err != nil {
		return err
	}

	chatgpt.ApiKey = os.Getenv("OPENAI_API_KEY")
	tuling.ApiKey = os.Getenv("TULING_API_KEY")

	go service.InitKeywordMap(keywordBytes)

	// 注册路由
	router := setupRouter(configSettings)
	if err := router.Run(fmt.Sprintf(":%d", configSettings.AppConfig.Port)); err != nil {
		return fmt.Errorf("[ERROR] server startup failed, err:%v", err)
	}

	return nil
}

func getConfigSettings(configFile string) (*models.ConfigSettings, error) {
	if configFile == "" {
		return nil, fmt.Errorf("[ERROR] config file not specified")
	}

	fileExt := path.Ext(configFile)
	if fileExt == ".yml" || fileExt == ".yaml" {
		return config.NewFile(configFile).RetrieveConfig()
	} else {
		return nil, fmt.Errorf("[ERROR] config file only support .yml or .yaml format")
	}
}

func setupRouter(cs *models.ConfigSettings) *gin.Engine {
	gin.SetMode(cs.AppConfig.Mode)

	router := gin.Default()
	//模板文件
	templates := template.Must(template.New("").ParseFS(fs, "templates/*.html"))
	router.SetHTMLTemplate(templates)
	//静态文件
	router.StaticFS("/public", http.FS(fs))

	router.GET("/", controller.IndexHandler)

	ws := controller.MessageController{
		WechatService: struct{ Config models.WechatConfig }{Config: cs.WechatConfig},
	}
	weChatGroup := router.Group("weChat")
	{
		//签名回调
		weChatGroup.GET("/receiveMessage", ws.ReceiveMessage)
		//接收发送给公众号的消息
		weChatGroup.POST("/receiveMessage", ws.ReceiveMessage)
	}
	return router
}
