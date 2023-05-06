package message

import (
	"time"
	"weChatRobot-go/models"
)

func BuildRespTextMessage(fromUserName, toUserName, content string) models.RespTextMessage {
	respMessage := models.RespTextMessage{
		Content: content,
	}
	respMessage.FromUserName = fromUserName
	respMessage.ToUserName = toUserName
	respMessage.CreateTime = time.Now().Unix()
	respMessage.MsgType = "text"
	return respMessage
}

func BuildRespNewsMessage(fromUserName, toUserName string, articles []models.ArticleItem) models.RespNewsMessage {
	respMessage := models.RespNewsMessage{
		ArticleCount: len(articles),
		Articles:     articles,
	}
	respMessage.FromUserName = fromUserName
	respMessage.ToUserName = toUserName
	respMessage.CreateTime = time.Now().Unix()
	respMessage.MsgType = "news"
	return respMessage
}
