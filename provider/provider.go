package provider

import "weChatRobot-go/model"

type Provider interface {
	RetrieveConfig() (*model.Config, error)
}
