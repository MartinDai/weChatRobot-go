package chatgpt_test

import (
	"fmt"
	"testing"
	"weChatRobot-go/chatgpt"
)

func TestGetRespMessage(t *testing.T) {
	chatgpt.ApiKey = "sk-*****"

	var respMessage = chatgpt.GetRespMessage("aaa", "bbb", "Suggest one name for a horse.")
	fmt.Printf("respMessage:%v", respMessage)
}
