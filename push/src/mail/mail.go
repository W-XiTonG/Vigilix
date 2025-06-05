package mail

import (
	"log"
	"net"
	"net/smtp"
	"os"
)

func Mail(Type int, Body, File, Contents, ReceiverEmail, SenderEmail, Subject, SmtpServer, SenderPassword string, CcEmails []string) {
	log.Println("正在发送邮件……")
	// Type : 1：文件传递，2：统一内容
	var Content string
	if Type == 0 {
		Content = Body
	} else if Type == 1 {
		msgPath := File
		content, err := os.ReadFile(msgPath)
		Content = string(content)
		if err != nil {
			log.Printf("邮件：读取文件失败 %s: %v", msgPath, err)
		}
	} else if Type == 2 {
		Content = Contents
	}
	var ccHeader string
	if len(CcEmails) > 0 {
		ccHeader = "Cc: " + jsonEmails(CcEmails) + "\r\n"
		//log.Println(ccHeader)
	}
	// 构建邮件内容
	message := []byte("To: " + ReceiverEmail + "\r\n" +
		ccHeader +
		"From: " + SenderEmail + "\r\n" +
		"Subject: " + Subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"" + "\r\n" +
		"\r\n" + Content)
	// 创建认证信息
	host, _, err := net.SplitHostPort(SmtpServer)
	auth := smtp.PlainAuth("", SenderEmail, SenderPassword, host)
	// 创建收件信息
	recipients := []string{ReceiverEmail}
	recipients = append(recipients, CcEmails...)
	// 发送邮件
	err = smtp.SendMail(SmtpServer, auth, SenderEmail, recipients, message)
	if err != nil {
		log.Printf("邮件发送状态未知，具体请检查邮箱是否收到，错误信息如下：\n%v", err)
		print("")
	} else {
		log.Println("邮件：发送成功！")
	}
}
