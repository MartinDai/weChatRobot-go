package ark

import (
	"context"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"os"
	"weChatRobot-go/pkg/logger"
	"weChatRobot-go/pkg/util"
)

type Doubao struct {
	client     *arkruntime.Client
	endpointId string
}

func NewDoubao() *Doubao {
	apiKey := os.Getenv("ARK_API_KEY")
	endpointId := os.Getenv("ARK_ENDPOINT_ID")
	var client *arkruntime.Client
	if apiKey != "" && endpointId != "" {
		client = arkruntime.NewClientWithApiKey(
			apiKey,
			arkruntime.WithBaseUrl("https://ark.cn-beijing.volces.com/api/v3"),
			arkruntime.WithRegion("cn-beijing"),
		)
	}

	return &Doubao{
		client:     client,
		endpointId: endpointId,
	}
}

func (doubao *Doubao) ProcessText(fromUserName, toUserName, content string) interface{} {
	if doubao.client == nil {
		return nil
	}

	request := model.ChatCompletionRequest{
		Model: doubao.endpointId,
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("你是一个AI助手，尽量保证回复内容在200个字符以内"),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(content),
				},
			},
		},
	}

	var response model.ChatCompletionResponse
	var err error
	if response, err = doubao.client.CreateChatCompletion(context.Background(), request); err != nil {
		logger.Error(err, "doubao Completion error")
		return nil
	}
	return util.BuildRespTextMessage(fromUserName, toUserName, *response.Choices[0].Message.Content.StringValue)
}
