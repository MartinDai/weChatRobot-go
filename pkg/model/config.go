package model

type Config struct {
	AppConfig    AppConfig    `yaml:"application"`
	WechatConfig WechatConfig `yaml:"wechat"`
	LoggerConfig LoggerConfig `yaml:"logger"`
}

type AppConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type WechatConfig struct {
	Token string `yaml:"token"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}
