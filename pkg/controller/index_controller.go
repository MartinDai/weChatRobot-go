package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"weChatRobot-go/pkg/logger"
	"weChatRobot-go/pkg/model"
	"weChatRobot-go/pkg/service"
)

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

type MessageController struct {
	wechatService *service.WechatService
}

func NewMessageController(wc *model.WechatConfig) *MessageController {
	wechatService := service.NewWechatService(wc)
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
		if err := c.ShouldBindXML(&reqMessage); err != nil {
			_, _ = fmt.Fprint(c.Writer, "系统处理消息异常")
			logger.Error(err, "解析XML出错")
			return
		}

		logger.Info("收到消息", "reqMessage", reqMessage)
		respXmlStr := mc.wechatService.GetResponseMessage(reqMessage)
		logger.Info("响应消息", "respXmlStr", respXmlStr)

		_, _ = fmt.Fprint(c.Writer, respXmlStr)
	}
}
