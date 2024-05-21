package dashscope

import (
	"fmt"
	"testing"
	"weChatRobot-go/pkg/util"
)

func TestGetRespMessage(t *testing.T) {
	dashscope := NewDashscope()
	var respMessage = dashscope.ProcessText("aaa", "bbb", "什么是GPT")
	fmt.Printf("respMessage:%v", util.ToJsonString(respMessage))
}
