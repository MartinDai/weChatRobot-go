package provider

import (
	"weChatRobot-go/pkg/model"
)

type Provider interface {
	RetrieveConfig() (*model.Config, error)
}
