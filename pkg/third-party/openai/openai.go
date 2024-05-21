package openai

import (
	"context"
	"github.com/otiai10/openaigo"
	"net/http"
	"net/url"
	"os"
	"weChatRobot-go/pkg/logger"
	"weChatRobot-go/pkg/util"
)

type OpenAI struct {
	client *openaigo.Client
}

func NewOpenAI() *OpenAI {
	apiKey := os.Getenv("OPENAI_API_KEY")
	var client *openaigo.Client
	if apiKey != "" {
		baseDomain := os.Getenv("OPENAI_BASE_DOMAIN")
		if baseDomain != "" && !util.ValidateAddress(baseDomain) {
			logger.Fatal("OPENAI_BASE_DOMAIN is not valid", "openaiBaseDomain", baseDomain)
		}

		proxyAddress := os.Getenv("OPENAI_PROXY")
		if proxyAddress != "" && !util.ValidateAddress(proxyAddress) {
			logger.Fatal("OPENAI_PROXY is not valid", "openaiBaseDomain", baseDomain)
		}

		client = openaigo.NewClient(apiKey)
		if baseDomain != "" {
			client.BaseURL = "https://" + baseDomain + "/v1"
		}
		if proxyAddress != "" {
			proxyUrl, _ := url.Parse("http://" + proxyAddress)
			transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
			client.HTTPClient = &http.Client{Transport: transport}
		}
	}

	return &OpenAI{
		client: client,
	}
}

func (openai *OpenAI) ProcessText(fromUserName, toUserName, content string) interface{} {
	if openai.client == nil {
		return nil
	}

	request := openaigo.ChatCompletionRequestBody{
		Model: "gpt-3.5-turbo",
		Messages: []openaigo.ChatMessage{
			{Role: "system", Content: "你是一个AI助手，保持回复内容尽量简短"},
			{Role: "user", Content: content},
		},
	}
	ctx := context.Background()

	var response openaigo.ChatCompletionResponse
	var err error
	if response, err = openai.client.Chat(ctx, request); err != nil {
		logger.Error(err, "GPT Completion error")
		return nil
	}
	return util.BuildRespTextMessage(fromUserName, toUserName, response.Choices[0].Message.Content)
}
