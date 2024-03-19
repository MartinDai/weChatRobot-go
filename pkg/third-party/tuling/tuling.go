package tuling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io"
	"net/http"
	"sync/atomic"
	"weChatRobot-go/pkg/logger"
	"weChatRobot-go/pkg/model"
	"weChatRobot-go/pkg/util"
)

const tulingApiUrl = "https://openapi.tuling123.com/openapi/api/v2"

type Tuling struct {
	apiKey string
	// 对微信传过来的userName做映射，因为有些userName的格式是图灵API不支持的
	userNameIdMap map[string]int32
	userIdAdder   int32
}

func NewTuling(apiKey string) *Tuling {
	return &Tuling{
		apiKey:        apiKey,
		userNameIdMap: make(map[string]int32),
		userIdAdder:   0,
	}
}

// GetRespMessage 从图灵机器人获取响应消息
func (t *Tuling) GetRespMessage(fromUserName, toUserName, content string) interface{} {
	userId := t.getUserId(toUserName)
	req := model.ReqParam{
		ReqType: 0,
		Perception: model.Perception{InputText: model.InputText{
			Text: content,
		}},
		UserInfo: model.UserInfo{
			ApiKey: t.apiKey,
			UserId: fmt.Sprintf("%d", userId),
		},
	}

	reqJsonBytes, _ := json.Marshal(req)
	reqJson := string(reqJsonBytes)
	logger.Info("请求图灵机器人", "reqJson", reqJson)

	var resp *http.Response
	var err error
	if resp, err = http.Post(tulingApiUrl, "application/json", bytes.NewReader(reqJsonBytes)); err != nil {
		logger.Error(err, "从图灵机器人获取响应内容报错")
		return nil
	}

	var result []byte
	if result, err = io.ReadAll(resp.Body); err != nil {
		logger.Error(err, "读取图灵机器人响应内容报错")
		return nil
	}

	var resultJson *simplejson.Json
	if resultJson, err = simplejson.NewJson(result); err != nil {
		logger.Error(err, "解析图灵机器人响应JSON报错")
		return nil
	}

	logger.Info("收到图灵机器人响应内容", "resultJson", resultJson)

	code, _ := resultJson.Get("intent").Get("code").Int()
	switch code {
	case model.ParamErrCode:
		return util.BuildRespTextMessage(fromUserName, toUserName, "我不是很理解你说的话")
	case model.NoResultCode:
		return util.BuildRespTextMessage(fromUserName, toUserName, "我竟无言以对！")
	case model.NoApiTimesCode:
		return util.BuildRespTextMessage(fromUserName, toUserName, "我今天已经说了太多话了，有点累，明天再来找我聊天吧！")
	case model.SuccessCode:
		var respTextMessage interface{}
		resultArray, _ := resultJson.Get("results").Array()
		for _, result := range resultArray {
			//转换成map结构
			if resultMap, ok := result.(map[string]interface{}); ok {
				if resultMap["resultType"].(string) == model.TextResultType {
					valueMap := resultMap["values"].(map[string]interface{})
					respTextMessage = util.BuildRespTextMessage(fromUserName, toUserName, valueMap["text"].(string))
					break
				}
			}
		}
		if respTextMessage != nil {
			return respTextMessage
		}
	}

	return nil
}

func (t *Tuling) getUserId(userName string) int32 {
	if userId, ok := t.userNameIdMap[userName]; ok {
		return userId
	} else {
		userId := atomic.AddInt32(&t.userIdAdder, 1)
		t.userNameIdMap[userName] = userId
		return userId
	}
}
