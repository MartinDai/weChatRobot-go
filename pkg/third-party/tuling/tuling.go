package tuling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"os"
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

func NewTuling() *Tuling {
	apiKey := os.Getenv("TULING_API_KEY")
	return &Tuling{
		apiKey:        apiKey,
		userNameIdMap: make(map[string]int32),
		userIdAdder:   0,
	}
}

// ProcessText 从图灵机器人获取响应消息
func (t *Tuling) ProcessText(fromUserName, toUserName, content string) interface{} {
	if t.apiKey == "" {
		return nil
	}

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

	var respBytes []byte
	if respBytes, err = io.ReadAll(resp.Body); err != nil {
		logger.Error(err, "读取图灵机器人响应内容报错")
		return nil
	}

	respStr := string(respBytes)
	logger.Info("收到图灵机器人响应内容", "respStr", respStr)

	if !gjson.Valid(respStr) {
		logger.Warn("图灵机器人响应内容不是json格式，无法解析")
		return nil
	}

	respJson := gjson.Parse(respStr)
	code := respJson.Get("intent.code").Int()
	switch code {
	case model.ParamErrCode:
		return util.BuildRespTextMessage(fromUserName, toUserName, "我不是很理解你说的话")
	case model.NoResultCode:
		return util.BuildRespTextMessage(fromUserName, toUserName, "我竟无言以对！")
	case model.NoApiTimesCode:
		return util.BuildRespTextMessage(fromUserName, toUserName, "我今天已经说了太多话了，有点累，明天再来找我聊天吧！")
	case model.SuccessCode:
		var respTextMessage interface{}
		resultArray := respJson.Get("results").Array()
		for _, result := range resultArray {
			if result.Get("resultType").String() == model.TextResultType {
				valueMap := result.Get("values")
				respTextMessage = util.BuildRespTextMessage(fromUserName, toUserName, valueMap.Get("text").String())
				break
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
