package dingding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// DingTalkMessage 定义钉钉消息结构体
type DingTalkMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func SendDing(Secret, WebhookURL, Message, File, Content string, Type int) error {
	// 获取当前时间戳，毫秒级
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	// 计算签名
	secret := Secret
	signature, err := Calculate(timestamp, secret)
	if err != nil {
		return fmt.Errorf("钉钉：计算签名失败: %w", err)
	}
	// 将签名和时间戳添加到 Webhook URL 中
	webhookURL := fmt.Sprintf("%s&timestamp=%s&sign=%s", WebhookURL, timestamp, signature)
	msg := DingTalkMessage{
		MsgType: "text",
	}

	// Type : 1：文件传递，2：统一内容
	if Type == 0 {
		msg.Text.Content = Message
	} else if Type == 1 {
		msgPath := File
		content, err := os.ReadFile(msgPath)
		msg.Text.Content = string(content)
		if err != nil {
			return fmt.Errorf("钉钉：读取文件失败 %s: %w", msgPath, err)
		}
	} else if Type == 2 {
		msg.Text.Content = Content
	}

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("钉钉：封装消息失败: %w", err)
	}
	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(msgJSON))
	if err != nil {
		return fmt.Errorf("钉钉：创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("钉钉：发送请求失败: %w", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Printf("resp.Body err = %v", err)
		}
	}()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}
	// 检查相应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("钉钉：请求失败，状态码 %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func DingDing(Secret, WebhookURL, Message, File, Content string, Type int) {
	log.Println("正在推送钉钉消息……")
	err := SendDing(Secret, WebhookURL, Message, File, Content, Type)
	if err != nil {
		log.Printf("钉钉：发送消息失败: %v", err)
		return
	}
	log.Println("钉钉：发送成功")
}
