package ark_test

import (
	"fmt"
	"testing"
	"weChatRobot-go/pkg/third-party/ark"
)

func TestGetRespMessage(t *testing.T) {
	doubao := ark.NewDoubao()
	var respMessage = doubao.ProcessText("aaa", "bbb", "什么是GPT")
	fmt.Printf("respMessage:%v", respMessage)
}
