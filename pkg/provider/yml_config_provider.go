package provider

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"weChatRobot-go/pkg/model"
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

func (ycp *ymlConfigProvider) RetrieveConfig() (*model.Config, error) {
	if ycp.filePath == "" {
		return nil, errors.New("config file not specified")
	}

	var config model.Config

	// 读取 YAML 文件内容
	data, err := os.ReadFile(ycp.filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil

}
