# weChatRobot

一个基于微信公众号的智能聊天机器人项目，支持根据关键字或者调用OpenAI、通义千问等大语言模型服务回复内容。

本项目还有Java实现的版本：<https://github.com/MartinDai/weChatRobot>

![qrcode](cmd/static/images/qrcode.jpg "扫码关注，体验智能机器人")

## 项目介绍

本项目是一个微信公众号项目，需配合微信公众号使用，在微信公众号配置本项目运行的服务器域名，用户关注公众号后，向公众号发送任意信息，公众号会根据用户发送的内容自动回复。
  
## 第三方依赖

- [gin](https://github.com/gin-gonic/gin)
- [gjson](https://github.com/tidwall/gjson)
- [yaml](https://gopkg.in/yaml.v3)
- [openaigo](https://github.com/otiai10/openaigo)

## 支持的功能

+ [x] 自定义关键字回复内容
+ [x] 调用OpenAI接口回复内容（需配置环境变量：`OPENAI_API_KEY`）
+ [x] 调用通义千问接口回复内容（需配置环境变量：`DASHSCOPE_API_KEY`）
+ [x] 调用图灵机器人(V2)接口回复内容（需配置环境变量：`TULING_API_KEY`）

## 使用说明

需要有微信公众号的帐号，没有的请戳[微信公众号申请](https://mp.weixin.qq.com/cgi-bin/readtemplate?t=register/step1_tmpl&lang=zh_CN)

内容响应来源的优先级`自定义关键字 > OpenAI > 通义千问 > 图灵机器人`

在微信公众号后台配置回调URL为<https://<your.domain>/weChat/receiveMessage>，其中`<your.domain>`替换成你自己的域名，token与`config.yml`里面配置的保持一致即可

### OpenAI

1. 如果需要使用OpenAI的回复内容则需要[创建OpenAI的API Key](https://platform.openai.com/account/api-keys)并配置在启动参数或者环境变量中
2. 可以通过配置启动参数或者环境变量`OPENAI_SERVER_URL`指定访问OpenAI服务的baseUrl
3. 可以通过配置启动参数或者环境变量`OPENAI_PROXY`使用代理服务访问OpenAI

### 通义千问

如果需要使用通义千问的回复内容则需要[创建通义千问的API Key](https://bailian.console.aliyun.com/#/api_key)并配置在启动参数或者环境变量中

### 图灵机器人

如果需要使用图灵机器人的回复内容则需要[注册图灵机器人帐号](http://tuling123.com/register/email.jhtml)获取相应的ApiKey并配置在启动参数或者环境变量中

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
