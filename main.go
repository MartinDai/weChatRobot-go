package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"
	"weChatRobot-go/controller"
	"weChatRobot-go/model"
	"weChatRobot-go/provider"
	"weChatRobot-go/service"
	"weChatRobot-go/third-party/chatgpt"
	"weChatRobot-go/third-party/tuling"
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
	conf, err := getConfig(configFile)
	if err != nil {
		return err
	}

	chatgpt.ApiKey = os.Getenv("OPENAI_API_KEY")
	tuling.ApiKey = os.Getenv("TULING_API_KEY")

	service.InitKeywordMap(keywordBytes)

	router := setupRouter(conf)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.AppConfig.Port),
		Handler: router,
	}

	ctx := context.Background()
	go func() {
		log.Printf("[INFO] Listening and serving HTTP on http://127.0.0.1%s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(fmt.Errorf("[ERROR] Server startup failed, Cause:%w", err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("[INFO] Shutdown Server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("[ERROR] Server Shutdown:", err)
	}
	log.Println("[INFO] Server exiting")

	return nil
}

func getConfig(configFile string) (*model.Config, error) {
	if configFile == "" {
		return nil, fmt.Errorf("[ERROR] config file not specified")
	}

	fileExt := path.Ext(configFile)
	if fileExt == ".yml" || fileExt == ".yaml" {
		return provider.NewFile(configFile).RetrieveConfig()
	} else {
		return nil, fmt.Errorf("[ERROR] config file only support .yml or .yaml format")
	}
}

func setupRouter(cs *model.Config) *gin.Engine {
	gin.SetMode(cs.AppConfig.Mode)

	router := gin.Default()
	//模板文件
	templates := template.Must(template.New("").ParseFS(fs, "templates/*.html"))
	router.SetHTMLTemplate(templates)
	//静态文件
	router.StaticFS("/public", http.FS(fs))

	router.GET("/", controller.IndexHandler)

	ws := controller.MessageController{
		WechatService: struct{ Config model.WechatConfig }{Config: cs.WechatConfig},
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
