package dashscope

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"weChatRobot-go/pkg/logger"
)

const generationUrl = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

type GenerationParam struct {
	Model string `json:"model"`
	Input Input  `json:"input"`
}

type Input struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GenerationResult struct {
	Output struct {
		Text         string `json:"text"`
		FinishReason string `json:"finish_reason"`
	} `json:"output"`
}

func (client *Client) call(param *GenerationParam) (*GenerationResult, error) {
	// 转换为JSON格式
	payload, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	// 构建请求
	req, err := http.NewRequest("POST", generationUrl, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	// 设置请求头部
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}

		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("API error: %s - %s", errorResponse.Code, errorResponse.Message)
	}

	result := &GenerationResult{}
	var respBytes []byte
	if respBytes, err = io.ReadAll(resp.Body); err != nil {
		logger.Error(err, "读取通义千问响应内容报错")
		return nil, nil
	}
	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
