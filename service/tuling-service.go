package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"weChatRobot-go/config"
	"weChatRobot-go/models"
)

const TulingApiUrl = "http://openapi.tuling123.com/openapi/api/v2"

//对微信传过来的userName做映射，因为有些userName的格式是图灵API不支持的
var userNameIdMap = make(map[string]int32)
var userIdAdder int32 = 0

// GetRespMessageFromTuling 从图灵机器人获取响应消息
func GetRespMessageFromTuling(fromUserName, toUserName, content string) interface{} {
	userId := GetUserId(toUserName)
	req := models.ReqParam{
		ReqType: 0,
		Perception: models.Perception{InputText: models.InputText{
			Text: content,
		}},
		UserInfo: models.UserInfo{
			ApiKey: config.ApiKey,
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
	if code == models.ParamErrCode {
		return BuildRespTextMessage(fromUserName, toUserName, "我不是很理解你说的话")
	} else if code == models.NoResultCode {
		return BuildRespTextMessage(fromUserName, toUserName, "我竟无言以对！")
	} else if code == models.NoApiTimesCode {
		return BuildRespTextMessage(fromUserName, toUserName, "我今天已经说了太多话了，有点累，明天再来找我聊天吧！")
	} else if code == models.TextCode {
		var respTextMessage interface{}
		resultArray, _ := resultJson.Get("results").Array()
		for _, result := range resultArray {
			//转换成map结构
			if resultMap, ok := result.(map[string]interface{}); ok {
				if resultMap["resultType"].(string) == models.TextResultType {
					valueMap := resultMap["values"].(map[string]interface{})
					respTextMessage = BuildRespTextMessage(fromUserName, toUserName, valueMap["text"].(string))
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

func GetUserId(userName string) int32 {
	if userId, ok := userNameIdMap[userName]; ok {
		return userId
	} else {
		userId := atomic.AddInt32(&userIdAdder, 1)
		userNameIdMap[userName] = userId
		return userId
	}
}
