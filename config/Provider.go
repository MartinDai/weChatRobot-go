package config

import "weChatRobot-go/models"

type Provider interface {
	RetrieveConfig() (*models.ConfigSettings, error)
}
