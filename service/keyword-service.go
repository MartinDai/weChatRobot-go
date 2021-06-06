package service

import (
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"weChatRobot-go/setting"
)

var keywordMessageMap = make(map[string]*simplejson.Json)

func InitKeywordMap() {
	keywordMessageBytes, err := ioutil.ReadFile(setting.Conf.KeywordLocation)
	if err != nil {
		log.Printf("读取关键字JSON文件报错:%v", err)
		return
	}

	keywordJson, err := simplejson.NewJson(keywordMessageBytes)
	if err != nil {
		log.Printf("解析关键字JSON文件报错:%v", err)
		return
	}

	keywordMap, err := keywordJson.Map()
	if err != nil {
		log.Printf("转换关键字JSON为Map报错:%v", err)
		return
	}

	for k, v := range keywordMap {
		log.Printf("初始化关键字map %v %v\n", k, v)
		keywordMessageMap[k] = keywordJson.Get(k)
	}
}
