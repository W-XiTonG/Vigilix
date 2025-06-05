package enterpriseWeChat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// SendWechatMessage 发送消息到企业微信机器人
func SendWechatMessage(webhookURL string, message interface{}) error {
	// 序列化消息内容
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("JSON marshal failed: %v", err)
	}
	log.Printf("[Push -> 企业微信机器人] : %s", string(msgBytes))
	// 发送POST请求
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(msgBytes))
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Printf("close body failed: %v", err)
		}
	}()

	// 处理响应
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API return error status: %s", resp.Status)
	}

	// 解析响应内容
	var result struct {
		Code    int    `json:"errcode"`
		Message string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response failed: %v", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("API error: %d - %s", result.Code, result.Message)
	}

	return nil
}
