package service

import (
	"github.com/bitly/go-simplejson"
	"weChatRobot-go/pkg/logger"
)

var keywordMessageMap = make(map[string]*simplejson.Json)

func InitKeywordMap(keywordBytes []byte) {
	var keywordJson *simplejson.Json
	var err error
	if keywordJson, err = simplejson.NewJson(keywordBytes); err != nil {
		logger.Error(err, "解析关键字JSON文件报错")
		return
	}

	var keywordMap map[string]interface{}
	if keywordMap, err = keywordJson.Map(); err != nil {
		logger.Error(err, "转换关键字JSON为Map报错")
		return
	}

	for k, v := range keywordMap {
		logger.Info("初始化关键字map", k, v)
		keywordMessageMap[k] = keywordJson.Get(k)
	}
}
