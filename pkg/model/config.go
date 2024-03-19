package model

type Config struct {
	AppConfig    AppConfig    `koanf:"application"`
	WechatConfig WechatConfig `koanf:"wechat"`
	LoggerConfig LoggerConfig `koanf:"logger"`
}

type AppConfig struct {
	Port int    `koanf:"port"`
	Mode string `koanf:"mode"`
}

type WechatConfig struct {
	Token string `koanf:"token"`
}

type LoggerConfig struct {
	Level string `koanf:"level"`
}
