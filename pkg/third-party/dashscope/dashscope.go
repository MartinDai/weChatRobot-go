package dashscope

import (
	"crypto/tls"
	"net/http"
	"os"
	"weChatRobot-go/pkg/logger"
	"weChatRobot-go/pkg/util"
)

type Dashscope struct {
	client *Client
}

func NewDashscope() *Dashscope {
	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	var client *Client
	if apiKey != "" {
		client = &Client{
			apiKey: apiKey,
			httpClient: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			},
		}
	}

	return &Dashscope{
		client: client,
	}
}

func (d *Dashscope) ProcessText(fromUserName, toUserName, content string) interface{} {
	if d.client == nil {
		return nil
	}

	param := &GenerationParam{
		Model: "qwen-turbo",
		Input: Input{
			Messages: []Message{
				{Role: "system", Content: "你是一个AI助手，尽量保证回复内容在200个字符以内"},
				{Role: "user", Content: content},
			},
		},
	}

	var result *GenerationResult
	var err error
	if result, err = d.client.call(param); err != nil {
		logger.Error(err, "dashscope generation error")
		return nil
	}

	return util.BuildRespTextMessage(fromUserName, toUserName, result.Output.Text)
}
