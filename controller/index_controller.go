package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"weChatRobot-go/model"
	"weChatRobot-go/service"
	"weChatRobot-go/third-party/chatgpt"
	"weChatRobot-go/third-party/tuling"
)

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

type MessageController struct {
	wechatService *service.WechatService
}

func NewMessageController(wc *model.WechatConfig, chatGPT *chatgpt.ChatGPT, tuling *tuling.Tuling) *MessageController {
	wechatService := service.NewWechatService(wc, chatGPT, tuling)
	return &MessageController{
		wechatService: wechatService,
	}
}

// ReceiveMessage 收到微信回调信息
func (mc *MessageController) ReceiveMessage(c *gin.Context) {
	if c.Request.Method == "GET" {
		signature := c.Query("signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		if mc.wechatService.CheckSignature(signature, timestamp, nonce) {
			_, _ = fmt.Fprint(c.Writer, c.Query("echostr"))
		} else {
			_, _ = fmt.Fprint(c.Writer, "你是谁？你想干嘛？")
		}
	} else {
		var reqMessage model.ReqMessage
		err := c.ShouldBindXML(&reqMessage)
		if err != nil {
			_, _ = fmt.Fprint(c.Writer, "系统处理消息异常")
			log.Printf("解析XML出错: %v\n", err)
			return
		}

		log.Printf("收到消息 %v\n", reqMessage)
		respXmlStr := mc.wechatService.GetResponseMessage(reqMessage)
		log.Printf("响应消息 %v\n", respXmlStr)

		_, _ = fmt.Fprint(c.Writer, respXmlStr)
	}
}
