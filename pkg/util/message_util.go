package util

import (
	"time"
	"weChatRobot-go/pkg/model"
)

func BuildRespTextMessage(fromUserName, toUserName, content string) model.RespTextMessage {
	respMessage := model.RespTextMessage{
		Content: content,
	}
	respMessage.FromUserName = fromUserName
	respMessage.ToUserName = toUserName
	respMessage.CreateTime = time.Now().Unix()
	respMessage.MsgType = "text"
	return respMessage
}

func BuildRespNewsMessage(fromUserName, toUserName string, articles []model.ArticleItem) model.RespNewsMessage {
	respMessage := model.RespNewsMessage{
		ArticleCount: len(articles),
		Articles:     articles,
	}
	respMessage.FromUserName = fromUserName
	respMessage.ToUserName = toUserName
	respMessage.CreateTime = time.Now().Unix()
	respMessage.MsgType = "news"
	return respMessage
}
