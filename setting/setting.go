package setting

import (
	"gopkg.in/ini.v1"
)

var Conf = new(AppConfig)

// AppConfig 应用程序配置
type AppConfig struct {
	Release bool `ini:"release"`
	Port    int  `ini:"port"`
	//微信公众号配置的token
	Token string `ini:"app.token"`
	//图灵机器人应用key
	ApiKey string `ini:"app.apiKey"`
	//关键字文件路径
	KeywordLocation string `ini:"keyword.location"`
}

func Init(file string) error {
	return ini.MapTo(Conf, file)
}
