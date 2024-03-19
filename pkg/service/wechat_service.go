package service

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"sort"
	"strings"
	"weChatRobot-go/pkg/logger"
	"weChatRobot-go/pkg/model"
	"weChatRobot-go/pkg/third-party/chatgpt"
	"weChatRobot-go/pkg/third-party/tuling"
	"weChatRobot-go/pkg/util"
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
		logger.Error(err, "XML编码出错")
		return ""
	}

	return string(respXmlStr)
}

func getRespMessageByEvent(fromUserName, toUserName, event string) interface{} {
	if event == model.EventTypeSubscribe {
		return util.BuildRespTextMessage(fromUserName, toUserName, "谢谢关注！可以开始跟我聊天啦😁")
	} else if event == model.EventTypeUnsubscribe {
		logger.Info("用户取消了订阅", "fromUserName", fromUserName)
	}
	return nil
}

func getRespMessageByKeyword(fromUserName, toUserName, keyword string) interface{} {
	v := GetResultByKeyword(keyword)
	if v.Exists() {
		msgType := v.Get("type").String()
		if msgType == "" {
			return nil
		}

		if msgType == model.MsgTypeText {
			content := v.Get("Content").String()
			return util.BuildRespTextMessage(fromUserName, toUserName, content)
		} else if msgType == model.MsgTypeNews {
			articleArray := v.Get("Articles").Array()
			var articleLength = len(articleArray)
			if articleLength == 0 {
				return nil
			}

			var articles = make([]model.ArticleItem, articleLength)
			for i, articleJson := range articleArray {
				var article model.Article
				article.Title = articleJson.Get("Title").String()
				article.Description = articleJson.Get("Description").String()
				article.PicUrl = articleJson.Get("PicUrl").String()
				article.Url = articleJson.Get("Url").String()

				var articleItem model.ArticleItem
				articleItem.Article = article
				articles[i] = articleItem
			}
			return util.BuildRespNewsMessage(fromUserName, toUserName, articles)
		}
	}
	return nil
}
