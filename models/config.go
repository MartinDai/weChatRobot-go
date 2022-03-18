package models

type ConfigSettings struct {
	AppConfig    AppConfig    `koanf:"application"`
	WechatConfig WechatConfig `koanf:"wechat"`
	TulingConfig TulingConfig `koanf:"tuling"`
}

type AppConfig struct {
	Port int    `koanf:"port"`
	Mode string `koanf:"mode"`
}

type WechatConfig struct {
	Token string `koanf:"token"`
}

type TulingConfig struct {
	AppKey string `koanf:"api_key"`
}
