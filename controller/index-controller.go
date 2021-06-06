package controller

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sort"
	"strings"
	"weChatRobot-go/models"
	"weChatRobot-go/service"
	"weChatRobot-go/setting"
)

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// ReceiveMessage 收到微信回调信息
func ReceiveMessage(c *gin.Context) {
	if c.Request.Method == "GET" {
		signature := c.Query("signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		if CheckSignature(signature, timestamp, nonce) {
			_, _ = fmt.Fprint(c.Writer, c.Query("echostr"))
		} else {
			_, _ = fmt.Fprint(c.Writer, "你是谁？你想干嘛？")
		}
	} else {
		var reqMessage models.ReqMessage
		err := c.ShouldBindXML(&reqMessage)
		if err != nil {
			_, _ = fmt.Fprint(c.Writer, "系统处理消息异常")
			log.Printf("解析XML出错: %v\n", err)
			return
		}

		log.Printf("收到消息 %v\n", reqMessage)
		respXmlStr := service.GetResponseMessage(reqMessage)
		log.Printf("响应消息 %v\n", respXmlStr)

		_, _ = fmt.Fprint(c.Writer, respXmlStr)
	}
}

// CheckSignature 校验签名
func CheckSignature(signature, timestamp, nonce string) bool {
	if signature == "" || timestamp == "" || nonce == "" {
		return false
	}

	arr := []string{setting.Conf.Token, timestamp, nonce}
	// 将token、timestamp、nonce三个参数进行字典序排序
	sort.Strings(arr)
	//拼接字符串
	content := strings.Join(arr, "")
	//sha1签名
	sha := sha1.New()
	sha.Write([]byte(content))
	sha1Value := hex.EncodeToString(sha.Sum(nil))

	return signature == sha1Value
}
