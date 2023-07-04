package chatgpt_test

import (
	"fmt"
	"testing"
	"weChatRobot-go/third-party/chatgpt"
)

func TestGetRespMessage(t *testing.T) {
	chatGPT := chatgpt.NewChatGPT("sk-*****", "", "")
	var respMessage = chatGPT.GetRespMessage("aaa", "bbb", "什么是GPT")
	fmt.Printf("respMessage:%v", respMessage)
}
