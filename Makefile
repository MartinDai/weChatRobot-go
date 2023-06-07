GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

GOBUILD=CGO_ENABLED=0 go build -trimpath

.DEFAULT_GOAL := build

.PHONY: all
all: darwin_amd64 darwin_arm64 linux_amd64 linux_arm64 windows_amd64

.PHONY: darwin_amd64
darwin_amd64:
	GOOS=darwin GOARCH=amd64 $(MAKE) build

.PHONY: darwin_arm64
darwin_arm64:
	GOOS=darwin GOARCH=arm64 $(MAKE) build

.PHONY: linux_amd64
linux_amd64:
	GOOS=linux GOARCH=amd64 $(MAKE) build

.PHONY: linux_arm64
linux_arm64:
	GOOS=linux GOARCH=arm64 $(MAKE) build

.PHONY: windows_amd64
windows_amd64:
	GOOS=windows GOARCH=amd64 EXTENSION=.exe $(MAKE) build

.PHONY: build
build:
	$(GOBUILD) -o ./bin/weChatRobot_$(GOOS)_$(GOARCH)$(EXTENSION) .