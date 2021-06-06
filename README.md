# weChatRobot
一个简单的智能聊天机器人项目，基于微信公众号和图灵机器人(V2)开发。（本项目还有Java实现的版本https://github.com/MartinDai/weChatRobot）

![qrcode](static/images/qrcode.jpg "扫码关注，体验智能机器人")

## 功能介绍：
  本项目是一个微信公众号项目，需配合微信公众号使用，在微信公众号配置本项目运行的服务器域名，用户关注公众号后，向公众号发送任意信息，公众号会根据用户发送的内容自动回复。
  
## 涉及框架及技术
+ go 1.16
+ gin
+ simplejson

## 使用说明：
1. 使用之前需要有微信公众号的帐号以及图灵机器人的帐号，没有的请戳[微信公众号申请](https://mp.weixin.qq.com/cgi-bin/readtemplate?t=register/step1_tmpl&lang=zh_CN)和[图灵机器人帐号注册](http://tuling123.com/register/email.jhtml)。
2. 在conf目录下的config.ini文件里面配置app相关的key。
3. 微信公众号URL配置为`http://robot.doodl6.com/weChat/receiveMessage`,其中`robot.doodl6.com`是你自己的域名，token与`config.ini`文件配置一致即可。
4. 命令启动：`go run main.go conf/config.ini`
5. 编译启动：先执行`go build`命令编译，再执行`./weChatRobot-go conf/config.ini`启动

## 支持的功能
* [x] 自动回复文本消息，回复内容来自于图灵机器人
* [x] 自定义关键字回复内容
