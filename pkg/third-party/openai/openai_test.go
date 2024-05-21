package openai_test

import (
	"fmt"
	"testing"
	"weChatRobot-go/pkg/third-party/openai"
)

func TestGetRespMessage(t *testing.T) {
	chatGPT := openai.NewOpenAI()
	var respMessage = chatGPT.ProcessText("aaa", "bbb", "什么是GPT")
	fmt.Printf("respMessage:%v", respMessage)
}
