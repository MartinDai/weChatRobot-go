package main

import (
	"fmt"
	"os"
	"weChatRobot-go/routers"
	"weChatRobot-go/service"
	"weChatRobot-go/setting"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage：./weChatRobot-go conf/config.ini")
		return
	}

	// 加载配置文件
	if err := setting.Init(os.Args[1]); err != nil {
		fmt.Printf("load config from file failed, err:%v\n", err)
		return
	}

	go service.InitKeywordMap()

	// 注册路由
	r := routers.SetupRouter()
	if err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port)); err != nil {
		fmt.Printf("server startup failed, err:%v\n", err)
	}
}
