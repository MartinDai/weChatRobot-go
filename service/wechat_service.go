package service

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"log"
	"sort"
	"strings"
	"weChatRobot-go/model"
	"weChatRobot-go/third-party/chatgpt"
	"weChatRobot-go/third-party/tuling"
	"weChatRobot-go/util"
)

type WechatService struct {
	config  *model.WechatConfig
	chatGPT *chatgpt.ChatGPT
	tuling  *tuling.Tuling
}

func NewWechatService(wc *model.WechatConfig, chatGPT *chatgpt.ChatGPT, tuling *tuling.Tuling) *WechatService {
	return &WechatService{
		config:  wc,
		chatGPT: chatGPT,
		tuling:  tuling,
	}
}

// CheckSignature 校验签名
func (ws *WechatService) CheckSignature(signature, timestamp, nonce string) bool {
	if signature == "" || timestamp == "" || nonce == "" {
		return false
	}

	arr := []string{ws.config.Token, timestamp, nonce}
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

func (ws *WechatService) GetResponseMessage(reqMessage model.ReqMessage) string {
	var respMessage interface{}
	if reqMessage.MsgType == model.MsgTypeEvent {
		respMessage = getRespMessageByEvent(reqMessage.ToUserName, reqMessage.FromUserName, reqMessage.Event)
	} else if reqMessage.MsgType == model.MsgTypeText {
		respMessage = getRespMessageByKeyword(reqMessage.ToUserName, reqMessage.FromUserName, reqMessage.Content)

		//优先使用ChatGPT响应
		if respMessage == nil && ws.chatGPT != nil {
			respMessage = ws.chatGPT.GetRespMessage(reqMessage.ToUserName, reqMessage.FromUserName, reqMessage.Content)
		}

		if respMessage == nil && ws.tuling != nil {
			respMessage = ws.tuling.GetRespMessage(reqMessage.ToUserName, reqMessage.FromUserName, reqMessage.Content)
		}
	} else {
		respMessage = util.BuildRespTextMessage(reqMessage.ToUserName, reqMessage.FromUserName, "我只对文字感兴趣[悠闲]")
	}

	if respMessage == nil {
		//最后兜底，如果没有响应，则返回输入的文字
		respMessage = util.BuildRespTextMessage(reqMessage.ToUserName, reqMessage.FromUserName, reqMessage.Content)
	}

	var respXmlStr []byte
	var err error
	if respXmlStr, err = xml.Marshal(&respMessage); err != nil {
		log.Printf("XML编码出错: %v\n", err)
		return ""
	}

	return string(respXmlStr)
}

func getRespMessageByEvent(fromUserName, toUserName, event string) interface{} {
	if event == model.EventTypeSubscribe {
		return util.BuildRespTextMessage(fromUserName, toUserName, "谢谢关注！可以开始跟我聊天啦😁")
	} else if event == model.EventTypeUnsubscribe {
		log.Printf("用户[%v]取消了订阅", fromUserName)
	}
	return nil
}

func getRespMessageByKeyword(fromUserName, toUserName, keyword string) interface{} {
	v, ok := keywordMessageMap[keyword]
	if ok {
		var msgType string
		var err error
		if msgType, err = v.Get("type").String(); err != nil {
			return nil
		}

		if msgType == model.MsgTypeText {
			content, _ := v.Get("Content").String()
			return util.BuildRespTextMessage(fromUserName, toUserName, content)
		} else if msgType == model.MsgTypeNews {
			var articleArray []interface{}
			if articleArray, err = v.Get("Articles").Array(); err != nil {
				return nil
			}

			var articleLength = len(articleArray)
			var articles = make([]model.ArticleItem, articleLength)
			for i, articleJson := range articleArray {
				if eachArticle, ok := articleJson.(map[string]interface{}); ok {
					var article model.Article
					article.Title = eachArticle["Title"].(string)
					article.Description = eachArticle["Description"].(string)
					article.PicUrl = eachArticle["PicUrl"].(string)
					article.Url = eachArticle["Url"].(string)

					var articleItem model.ArticleItem
					articleItem.Article = article
					articles[i] = articleItem
				}
			}
			return util.BuildRespNewsMessage(fromUserName, toUserName, articles)
		}
	}
	return nil
}
