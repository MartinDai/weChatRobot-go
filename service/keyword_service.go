package service

import (
	"github.com/bitly/go-simplejson"
	"log"
)

var keywordMessageMap = make(map[string]*simplejson.Json)

func InitKeywordMap(keywordBytes []byte) {
	var keywordJson *simplejson.Json
	var err error
	if keywordJson, err = simplejson.NewJson(keywordBytes); err != nil {
		log.Printf("解析关键字JSON文件报错:%v", err)
		return
	}

	var keywordMap map[string]interface{}
	if keywordMap, err = keywordJson.Map(); err != nil {
		log.Printf("转换关键字JSON为Map报错:%v", err)
		return
	}

	for k, v := range keywordMap {
		log.Printf("初始化关键字map %v %v\n", k, v)
		keywordMessageMap[k] = keywordJson.Get(k)
	}
}
