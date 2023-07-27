package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"
	"weChatRobot-go/controller"
	"weChatRobot-go/logger"
	"weChatRobot-go/model"
	"weChatRobot-go/provider"
	"weChatRobot-go/service"
	"weChatRobot-go/third-party/chatgpt"
	"weChatRobot-go/third-party/tuling"
	"weChatRobot-go/util"
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
		logger.Fatal(err, "process config error")
	}
}

func runApp(configFile string) error {
	var config *model.Config
	var err error
	if config, err = getConfig(configFile); err != nil {
		return err
	}

	service.InitKeywordMap(keywordBytes)

	router := setupRouter(config)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.AppConfig.Port),
		Handler: router,
	}

	ctx := context.Background()
	go func() {
		logger.Info("Listening and serving HTTP on http://127.0.0.1%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err, "Server startup failed")
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutdown Server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		logger.Fatal(err, "Server Shutdown failed")
	}
	logger.Info("Server exited")

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

func setupRouter(config *model.Config) *gin.Engine {
	gin.SetMode(config.AppConfig.Mode)

	router := gin.Default()
	//模板文件
	templates := template.Must(template.New("").ParseFS(fs, "templates/*.html"))
	router.SetHTMLTemplate(templates)
	//静态文件
	router.StaticFS("/public", http.FS(fs))

	router.GET("/", controller.IndexHandler)

	openaiApiKey := os.Getenv("OPENAI_API_KEY")
	var chatGPT *chatgpt.ChatGPT
	if openaiApiKey != "" {
		openaiBaseDomain := os.Getenv("OPENAI_BASE_DOMAIN")
		if openaiBaseDomain != "" && !util.ValidateAddress(openaiBaseDomain) {
			logger.Fatalf("OPENAI_BASE_DOMAIN is not valid:%s", openaiBaseDomain)
		}

		openaiProxy := os.Getenv("OPENAI_PROXY")
		if openaiProxy != "" && !util.ValidateAddress(openaiProxy) {
			logger.Fatalf("OPENAI_PROXY is not valid:%v", openaiBaseDomain)
		}

		chatGPT = chatgpt.NewChatGPT(openaiApiKey, openaiBaseDomain, openaiProxy)
	}

	var tl *tuling.Tuling
	tulingApiKey := os.Getenv("TULING_API_KEY")
	if tulingApiKey != "" {
		tl = tuling.NewTuling(tulingApiKey)
	}

	ws := controller.NewMessageController(&config.WechatConfig, chatGPT, tl)
	weChatGroup := router.Group("weChat")
	{
		//签名回调
		weChatGroup.GET("/receiveMessage", ws.ReceiveMessage)
		//接收发送给公众号的消息
		weChatGroup.POST("/receiveMessage", ws.ReceiveMessage)
	}
	return router
}
