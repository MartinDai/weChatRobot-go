# weChatRobot

一个基于微信公众号的智能聊天机器人项目，支持图灵机器人(V2)和ChatGPT对话模式回复内容

本项目还有Java实现的版本：<https://github.com/MartinDai/weChatRobot>

![qrcode](cmd/static/images/qrcode.jpg "扫码关注，体验智能机器人")

## 项目介绍

  本项目是一个微信公众号项目，需配合微信公众号使用，在微信公众号配置本项目运行的服务器域名，用户关注公众号后，向公众号发送任意信息，公众号会根据用户发送的内容自动回复。
  
## 第三方依赖

- [gin](https://github.com/gin-gonic/gin)
- [simplejson](https://github.com/bitly/go-simplejson)
- [yaml](https://gopkg.in/yaml.v3)
- [openaigo](https://github.com/otiai10/openaigo)

## 支持的功能

+ [x] 自定义关键字回复内容
+ [x] 调用ChatGPT接口回复内容（需配置环境变量：`OPENAI_API_KEY`）
+ [x] 调用图灵机器人(V2)接口回复内容（需配置环境变量：`TULING_API_KEY`）

## 使用说明

1. 使用之前需要有微信公众号的帐号，没有的请戳[微信公众号申请](https://mp.weixin.qq.com/cgi-bin/readtemplate?t=register/step1_tmpl&lang=zh_CN)
2. 如果需要使用图灵机器人的回复内容则需要[注册图灵机器人帐号](http://tuling123.com/register/email.jhtml)获取相应的ApiKey并配置在环境变量中
3. 如果需要使用ChatGPT的回复内容则需要[创建OpenAI的API Key](https://platform.openai.com/account/api-keys)并配置在环境变量中
4. 可以通过配置环境变量`OPENAI_BASE_DOMAIN`更换访问OpenAI的域名
5. 可以通过配置环境变量`OPENAI_PROXY`使用代理服务访问OpenAI
6. 内容响应来源的优先级`自定义关键 > ChatGPT > 图灵机器人`
7. 在微信公众号后台配置回调URL为<https://wechatrobot.doodl6.com/weChat/receiveMessage>，其中`wechatrobot.doodl6.com`是你自己的域名，token与`config.yml`里面配置的保持一致即可

## 本地开发

### GoLand

需要配置Program Arguments为`-config config.yml`，然后运行main.go

### VS Code

可以直接使用`launch.json`的配置，里面还包含了环境变量的配置直接设置即可

## 编译运行

### 通过Makefile构建

构建适合当前系统的可执行文件

```shell
make
```

构建指定平台架构的可执行文件

```shell
make linux_amd64
```

编译全平台的可执行文件

```shell
make all
```

生成的可执行文件在`bin`目录下，执行`./bin/weChatRobot_darwin_arm64 -config config.yml`启动运行，其中`_darwin_arm64`后缀不同系统架构不一样

## Docker运行

构建适用于当前操作系统架构的镜像

```shell
docker build --no-cache -t wechatrobot-go:latest .
```

构建指定架构的镜像

```shell
docker buildx build --no-cache -t wechatrobot-go:latest --platform=linux/amd64 -o type=docker .
```

后台启动镜像

```shell
docker run --name wechatrobot-go -p 8080:8080 -d wechatrobot-go:latest
```
