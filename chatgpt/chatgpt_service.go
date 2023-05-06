package chatgpt

import (
	"context"
	"fmt"
	"weChatRobot-go/message"

	"github.com/otiai10/openaigo"
)

var ApiKey string

func GetRespMessage(fromUserName, toUserName, content string) interface{} {
	client := openaigo.NewClient(ApiKey)
	// proxy_url, _ := url.Parse("http://127.0.0.1:7890")
	// transport := &http.Transport{Proxy: http.ProxyURL(proxy_url)}
	// client.HTTPClient = &http.Client{Transport: transport}

	request := openaigo.ChatCompletionRequestBody{
		Model: "gpt-3.5-turbo",
		Messages: []openaigo.ChatMessage{
			{Role: "user", Content: content},
		},
	}
	ctx := context.Background()
	response, err := client.Chat(ctx, request)
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return nil
	}
	return message.BuildRespTextMessage(fromUserName, toUserName, response.Choices[0].Message.Content)
}
