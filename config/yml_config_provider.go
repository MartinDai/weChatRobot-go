package config

import (
	"errors"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"weChatRobot-go/models"
)

type ymlConfigProvider struct {
	filePath string
}

// NewFile returns a new Provider that reads the configuration from the given file.
func NewFile(filePath string) Provider {
	return &ymlConfigProvider{
		filePath: filePath,
	}
}

func (ycp *ymlConfigProvider) RetrieveConfig() (*models.ConfigSettings, error) {
	if ycp.filePath == "" {
		return nil, errors.New("config file not specified")
	}

	var config models.ConfigSettings
	k := koanf.New("::")
	_ = k.Load(file.Provider(ycp.filePath), yaml.Parser())
	err := k.Unmarshal("", &config)
	if err != nil {
		return nil, err
	}
	return &config, nil

}
