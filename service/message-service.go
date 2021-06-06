package service

import (
	"encoding/xml"
	"log"
	"time"
	"weChatRobot-go/models"
)

func GetResponseMessage(reqMessage models.ReqMessage) string {
	var respMessage interface{}
	if reqMessage.MsgType == models.MsgTypeEvent {
		respMessage = GetRespMessageByEvent(reqMessage.ToUserName, reqMessage.FromUserName, reqMessage.Event)
	} else if reqMessage.MsgType == models.MsgTypeText {
		respMessage = GetRespMessageByKeyword(reqMessage.ToUserName, reqMessage.FromUserName, reqMessage.Content)
		if respMessage == nil {
			respMessage = GetRespMessageFromTuling(reqMessage.ToUserName, reqMessage.FromUserName, reqMessage.Content)
		}
	} else {
		respMessage = BuildRespTextMessage(reqMessage.ToUserName, reqMessage.FromUserName, "我只对文字感兴趣[悠闲]")
	}

	if respMessage == nil {
		return ""
	} else {
		respXmlStr, err := xml.Marshal(&respMessage)
		if err != nil {
			log.Printf("XML编码出错: %v\n", err)
			return ""
		}

		return string(respXmlStr)
	}
}

func GetRespMessageByEvent(fromUserName, toUserName, event string) interface{} {
	if event == models.EventTypeSubscribe {
		return BuildRespTextMessage(fromUserName, toUserName, "谢谢关注！可以开始跟我聊天啦😁")
	} else if event == models.EventTypeUnsubscribe {
		log.Printf("用户[%v]取消了订阅", fromUserName)
	}
	return nil
}

func GetRespMessageByKeyword(fromUserName, toUserName, keyword string) interface{} {
	v, ok := keywordMessageMap[keyword]
	if ok {
		msgType, err := v.Get("type").String()
		if err != nil {
			return nil
		}

		if msgType == models.MsgTypeText {
			content, _ := v.Get("Content").String()
			return BuildRespTextMessage(fromUserName, toUserName, content)
		} else if msgType == models.MsgTypeNews {
			articleArray, err := v.Get("Articles").Array()
			if err != nil {
				return nil
			}

			var articleLength = len(articleArray)
			var articles = make([]models.ArticleItem, articleLength)
			for i, articleJson := range articleArray {
				if eachArticle, ok := articleJson.(map[string]interface{}); ok {
					var article models.Article
					article.Title = eachArticle["Title"].(string)
					article.Description = eachArticle["Description"].(string)
					article.PicUrl = eachArticle["PicUrl"].(string)
					article.Url = eachArticle["Url"].(string)

					var articleItem models.ArticleItem
					articleItem.Article = article
					articles[i] = articleItem
				}
			}
			return BuildRespNewsMessage(fromUserName, toUserName, articles)
		}
	}
	return nil
}

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
