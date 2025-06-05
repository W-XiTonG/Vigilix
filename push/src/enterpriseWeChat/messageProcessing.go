package enterpriseWeChat

import (
	"log"
	"os"
)

// PushEnterpriseWeChat 消息处理
func PushEnterpriseWeChat(webhookURL, message, Content, File string, AteSpecifyStatus bool,
	mentionedMobileList []string, contentSType int) {
	// TextMessage 文本消息结构体
	type TextMessage struct {
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
			//MentionedList       []string `json:"mentioned_list,omitempty"`
			MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
		} `json:"text"`
	}
	// 创建文本消息
	msg := TextMessage{
		MsgType: "text",
	}
	// Type :0:各自配置， 1：文件传递，2：统一内容
	if contentSType == 0 {
		msg.Text.Content = message
	} else if contentSType == 1 {
		msgPath := File
		content, err := os.ReadFile(msgPath)
		msg.Text.Content = string(content)
		if err != nil {
			log.Printf("企业微信：读取文件失败 %s: %v", msgPath, err)
		}
	} else if contentSType == 2 {
		msg.Text.Content = Content
	}

	// 初始化MentionedMobileList为nil
	msg.Text.MentionedMobileList = nil

	// 处理@人逻辑
	if AteSpecifyStatus && len(mentionedMobileList) > 0 {
		// @指定人（仅当AteSpecifyStatus为true且有手机号时）
		msg.Text.MentionedMobileList = mentionedMobileList
	}

	// 发送消息
	if err := SendWechatMessage(webhookURL, msg); err != nil {
		log.Printf("Error sending message: %v\n", err)
	} else {
		log.Println("Message sent successfully")
	}
}
