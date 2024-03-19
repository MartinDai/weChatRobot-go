package service

import (
	"github.com/tidwall/gjson"
	"weChatRobot-go/pkg/logger"
)

var keywordResultMap gjson.Result

func InitKeywordMap(keywordBytes []byte) {
	if !gjson.ValidBytes(keywordBytes) {
		logger.Warn("关键字JSON文件格式不正确，跳过解析")
		return
	}

	result := gjson.ParseBytes(keywordBytes)
	keywordResultMap = result

	for k, v := range keywordResultMap.Map() {
		logger.Info("初始化关键字map", k, v.Value())
	}
}

func GetResultByKeyword(keyword string) gjson.Result {
	result := keywordResultMap.Get(keyword)
	return result
}
