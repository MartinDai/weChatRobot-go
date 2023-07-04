package tuling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"weChatRobot-go/model"
	"weChatRobot-go/util"

	"github.com/bitly/go-simplejson"
)

const TulingApiUrl = "https://openapi.tuling123.com/openapi/api/v2"

// 对微信传过来的userName做映射，因为有些userName的格式是图灵API不支持的
var userNameIdMap = make(map[string]int32)
var userIdAdder int32 = 0
var ApiKey string

// GetRespMessage 从图灵机器人获取响应消息
func GetRespMessage(fromUserName, toUserName, content string) interface{} {
	userId := getUserId(toUserName)
	req := model.ReqParam{
		ReqType: 0,
		Perception: model.Perception{InputText: model.InputText{
			Text: content,
		}},
		UserInfo: model.UserInfo{
			ApiKey: ApiKey,
			UserId: fmt.Sprintf("%d", userId),
		},
	}

	reqJsonBytes, _ := json.Marshal(req)
	reqJson := string(reqJsonBytes)
	log.Printf("请求图灵机器人参数 %v", reqJson)

	resp, err := http.Post(TulingApiUrl, "application/json", bytes.NewReader(reqJsonBytes))
	if err != nil {
		log.Printf("从图灵机器人获取响应内容报错,err:%v", err)
		return nil
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取图灵机器人响应内容报错,err:%v", err)
		return nil
	}

	resultJson, err := simplejson.NewJson(result)
	if err != nil {
		log.Printf("解析图灵机器人响应JSON报错:%v", err)
		return nil
	}

	log.Printf("收到图灵机器人响应内容 %v", resultJson)

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

func getUserId(userName string) int32 {
	if userId, ok := userNameIdMap[userName]; ok {
		return userId
	} else {
		userId := atomic.AddInt32(&userIdAdder, 1)
		userNameIdMap[userName] = userId
		return userId
	}
}
