package chatgpt

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"weChatRobot-go/util"

	"github.com/otiai10/openaigo"
)

type ChatGPT struct {
	client *openaigo.Client
}

func NewChatGPT(apiKey, baseDomain, proxyAddress string) *ChatGPT {
	client := openaigo.NewClient(apiKey)
	if baseDomain != "" {
		client.BaseURL = "https://" + baseDomain + "/v1"
	}
	if proxyAddress != "" {
		proxyUrl, _ := url.Parse("http://" + proxyAddress)
		transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		client.HTTPClient = &http.Client{Transport: transport}
	}
	return &ChatGPT{
		client: client,
	}
}

func (gpt *ChatGPT) GetRespMessage(fromUserName, toUserName, content string) interface{} {
	request := openaigo.ChatCompletionRequestBody{
		Model: "gpt-3.5-turbo",
		Messages: []openaigo.ChatMessage{
			{Role: "user", Content: content},
		},
	}
	ctx := context.Background()
	response, err := gpt.client.Chat(ctx, request)
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return nil
	}
	return util.BuildRespTextMessage(fromUserName, toUserName, response.Choices[0].Message.Content)
}
