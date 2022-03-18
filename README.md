# weChatRobot
一个简单的智能聊天机器人项目，基于微信公众号和图灵机器人(V2)开发。

本项目还有Java实现的版本：https://github.com/MartinDai/weChatRobot

![qrcode](static/images/qrcode.jpg "扫码关注，体验智能机器人")

## 项目介绍：
  本项目是一个微信公众号项目，需配合微信公众号使用，在微信公众号配置本项目运行的服务器域名，用户关注公众号后，向公众号发送任意信息，公众号会根据用户发送的内容自动回复。
  
## 涉及框架及技术
+ go 1.17
+ gin
+ simplejson
+ koanf

## 支持的功能
* [x] 自动回复文本消息，回复内容来自于图灵机器人
* [x] 自定义关键字回复内容

## 使用说明：
1. 使用之前需要有微信公众号的帐号以及图灵机器人的帐号，没有的请戳[微信公众号申请](https://mp.weixin.qq.com/cgi-bin/readtemplate?t=register/step1_tmpl&lang=zh_CN)和[图灵机器人帐号注册](http://tuling123.com/register/email.jhtml)。
2. 在`config.yml`文件里面配置相关的key。
3. 在微信公众号后台配置回调URL为`http://robot.doodl6.com/weChat/receiveMessage`,其中`robot.doodl6.com`是你自己的域名，token与第2点文件里面配置的保持一致即可。

## 本地开发
本地开发时（以GoLand为例）需要配置Program Arguments为`-config ./config.yml`

编译运行：在根目录执行`go build -o weChatRobot-go main.go`，编译得到可执行文件`weChatRobot-go`。 

执行`./weChatRobot-go -config ./config.yml`启动项目

编译适合当前系统的可执行文件：
```
make weChatRobot
```

编译全平台的可执行文件：
```
make all
```

生成的可执行文件在`bin`目录下

## Docker运行

构建适用于当前操作系统/架构的镜像
```
docker build --no-cache -t wechatrobot-go:latest .
```

后台启动项目
```
docker run --name wechatrobot-go -p 8080:8080 -d wechatrobot-go:latest
```