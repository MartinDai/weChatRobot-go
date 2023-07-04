package chatgpt_test

import (
	"fmt"
	"testing"
	"weChatRobot-go/third-party/chatgpt"
)

func TestGetRespMessage(t *testing.T) {
	chatgpt.ApiKey = "sk-*****"

	var respMessage = chatgpt.GetRespMessage("aaa", "bbb", "什么是GPT")
	fmt.Printf("respMessage:%v", respMessage)
}
