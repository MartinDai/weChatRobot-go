GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

.PHONY: all
all: weChatRobot-darwin_amd64 weChatRobot-darwin_arm64 weChatRobot-linux_amd64 weChatRobot-linux_arm64 weChatRobot-windows_amd64

.PHONY: weChatRobot-darwin_amd64
weChatRobot-darwin_amd64:
	GOOS=darwin  GOARCH=amd64 $(MAKE) weChatRobot

.PHONY: weChatRobot-darwin_arm64
weChatRobot-darwin_arm64:
	GOOS=darwin  GOARCH=arm64 $(MAKE) weChatRobot

.PHONY: weChatRobot-linux_amd64
weChatRobot-linux_amd64:
	GOOS=linux   GOARCH=amd64 $(MAKE) weChatRobot

.PHONY: weChatRobot-linux_arm64
weChatRobot-linux_arm64:
	GOOS=linux   GOARCH=arm64 $(MAKE) weChatRobot

.PHONY: weChatRobot-windows_amd64
weChatRobot-windows_amd64:
	GOOS=windows GOARCH=amd64 EXTENSION=.exe $(MAKE) weChatRobot

.PHONY: weChatRobot
weChatRobot:
	CGO_ENABLED=0 go build -trimpath -o ./bin/weChatRobot_$(GOOS)_$(GOARCH)$(EXTENSION) .