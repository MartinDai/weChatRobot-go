package service

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"log"
	"sort"
	"strings"
	"weChatRobot-go/chatgpt"
	"weChatRobot-go/message"
	"weChatRobot-go/models"
	"weChatRobot-go/tuling"
)

type WechatService struct {
	Config models.WechatConfig
}

// CheckSignature 校验签名
func (ws *WechatService) CheckSignature(signature, timestamp, nonce string) bool {
	if signature == "" || timestamp == "" || nonce == "" {
		return false
	}

	arr := []string{ws.Config.Token, timestamp, nonce}
	// 将token、timestamp、nonce三个参数进行字典序排序
	sort.Strings(arr)
	//拼接字符串
	content := strings.Join(arr, "")
	//sha1签名
	sha := sha1.New()
	sha.Write([]byte(content))
	sha1Value := hex.EncodeToString(sha.Sum(nil))

	return signature == sha1Value
}

func GetResponseMessage(reqMessage models.ReqMessage) string {
	var respMessage interface{}
	if reqMessage.MsgType == models.MsgTypeEvent {
		respMessage = getRespMessageByEvent(reqMessage.ToUserName, reqMessage.FromUserName, reqMessage.Event)
	} else if reqMessage.MsgType == models.MsgTypeText {
		respMessage = getRespMessageByKeyword(reqMessage.ToUserName, reqMessage.FromUserName, reqMessage.Content)

		//优先使用chatgpt响应
		if respMessage == nil && chatgpt.ApiKey != "" {
			respMessage = chatgpt.GetRespMessage(reqMessage.ToUserName, reqMessage.FromUserName, reqMessage.Content)
		}

		if respMessage == nil && tuling.ApiKey != "" {
			respMessage = tuling.GetRespMessage(reqMessage.ToUserName, reqMessage.FromUserName, reqMessage.Content)
		}
	} else {
		respMessage = message.BuildRespTextMessage(reqMessage.ToUserName, reqMessage.FromUserName, "我只对文字感兴趣[悠闲]")
	}

	if respMessage == nil {
		//最后兜底，如果没有响应，则返回输入的文字
		respMessage = message.BuildRespTextMessage(reqMessage.ToUserName, reqMessage.FromUserName, reqMessage.Content)
	}

	respXmlStr, err := xml.Marshal(&respMessage)
	if err != nil {
		log.Printf("XML编码出错: %v\n", err)
		return ""
	}

	return string(respXmlStr)
}

func getRespMessageByEvent(fromUserName, toUserName, event string) interface{} {
	if event == models.EventTypeSubscribe {
		return message.BuildRespTextMessage(fromUserName, toUserName, "谢谢关注！可以开始跟我聊天啦😁")
	} else if event == models.EventTypeUnsubscribe {
		log.Printf("用户[%v]取消了订阅", fromUserName)
	}
	return nil
}

func getRespMessageByKeyword(fromUserName, toUserName, keyword string) interface{} {
	v, ok := keywordMessageMap[keyword]
	if ok {
		msgType, err := v.Get("type").String()
		if err != nil {
			return nil
		}

		if msgType == models.MsgTypeText {
			content, _ := v.Get("Content").String()
			return message.BuildRespTextMessage(fromUserName, toUserName, content)
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
			return message.BuildRespNewsMessage(fromUserName, toUserName, articles)
		}
	}
	return nil
}
