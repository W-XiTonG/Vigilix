package SMS

import (
	"fmt"
	"log"
	"net"
	"os"
)

// SGIPSubmitSend submit发送
func SGIPSubmitSend(ContentSType int, conn net.Conn, UserNumber, Message, File, Content, SPNumber, ChargeNumber, CorpID, ServiceType string) error {
	// 短信内容
	var message string
	if ContentSType == 0 {
		message = Message
	} else if ContentSType == 1 {
		msgPath := File
		messageContent, err := os.ReadFile(msgPath)
		message = string(messageContent)
		if err != nil {
			return fmt.Errorf("短信：读取文件失败 %s: %w", msgPath, err)
		}
	} else if ContentSType == 2 {
		message = Content
	}

	messageCoding := uint8(8) // UTF-16BE 编码

	// 编码短信内容
	messageContent, err := encodeMessage(message, messageCoding)
	if err != nil {
		log.Printf("编码短信内容失败: %v\n", err)
		return err
	}
	messageLength := uint32(len(messageContent))

	// 获取当前时间
	//now := time.Now()
	//expireTime := getSGIPTime(now.Add(24 * time.Hour))
	//
	//// 定义 ScheduleTime（短消息定时发送的时间）
	//scheduleTime := getSGIPTime(now) // 立即发送

	// 创建 SGIP SUBMIT 请求消息
	req := SGIPSubmit{
		Header: SGIPHeader{
			CommandID:  SGIP_SUBMIT,           // SUBMIT 请求命令 ID
			SequenceID: generateSequenceID(2), // 生成 SequenceID
		},
		UserCount: 1, // 用户数量
		FeeType:   1, // 计费类型
		//FeeValue:   0, // 收费值
		//GivenValue: 0,               // 赠送花费
		AgentFlag:        1, // 代收费标志，0：应收；1：实收，字符
		MorelatetoMTFlag: 3, // 引起 MT 消息的原因
		Priority:         1,
		ReportFlag:       1, // 需要状态报告
		//TP_pid:           1,                // GSM 协议类型
		//TP_udhi:          1,                // GSM 协议类型
		MessageCoding:  messageCoding,  // 消息编码格式
		MessageType:    0,              // 消息类型（普通短信）
		MessageLength:  messageLength,  // 消息长度
		MessageContent: messageContent, // 消息内容
		//ExpireTime:     expireTime,     // 短消息寿命的终止时间
		//ScheduleTime:   scheduleTime,   // 短消息定时发送的时间
	}

	// 填充固定长度字段
	copy(req.SPNumber[:], []byte(SPNumber))
	copy(req.ChargeNumber[:], []byte(ChargeNumber))
	copy(req.UserNumber[:], []byte(UserNumber))
	copy(req.CorpID[:], []byte(CorpID))
	copy(req.ServiceType[:], []byte(ServiceType))

	// 初始化 Reserve 字段
	for i := 0; i < len(req.Reserve); i++ {
		req.Reserve[i] = 0x00
	}

	// 编码请求
	data, err := req.Encode()
	if err != nil {
		log.Printf("编码失败: %v\n", err)
		return err
	}

	// 发送 SUBMIT 请求
	_, err = conn.Write(data)
	if err != nil {
		log.Printf("发送 SUBMIT 请求失败: %v\n", err)
		return err
	}
	log.Println(UserNumber, "SUBMIT 请求发送成功")
	// 接收响应
	buf := make([]byte, 1024)
	//log.Printf("buf: %v\n", buf)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("接收响应失败: %v\n", err)
		return err
	}
	log.Printf("接收到 SUBMIT 响应: %x\n", buf[:n])
	// 打印状态码
	log.Println("SUBMIT状态码:", buf[20])
	return nil
}
