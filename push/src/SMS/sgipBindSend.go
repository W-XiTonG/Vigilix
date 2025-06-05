package SMS

import (
	"log"
	"net"
)

func SGIPBindSend(LoginName, LoginPassword, Message, File, Content, SPNumber, ChargeNumber, CorpID, ServiceType string, UserNumber []string, ContentSType int, conn net.Conn) {
	// 创建 SGIP BIND 请求消息
	req := SGIPBind{
		Header: SGIPHeader{
			CommandID:  SGIP_BIND,             // BIND 请求命令 ID
			SequenceID: generateSequenceID(1), // 生成 SequenceID
		},
		LoginType: 0x01, // 登录类型，0x01 表示客户端
	}

	// 填充 LoginName 和 LoginPassword
	copy(req.LoginName[:], []byte(LoginName))
	copy(req.LoginPassword[:], []byte(LoginPassword))

	// 初始化 Reserve 字段
	for i := 0; i < len(req.Reserve); i++ {
		req.Reserve[i] = 0x00
	}

	// 编码请求
	data, err := req.Encode()
	if err != nil {
		log.Printf("编码失败: %v\n", err)
		return
	}
	// 建立TCP链接
	//conn, err := net.Dial("tcp", provider.GetMailConfig().SMS.Sgip.SmgIpPort)
	//if err != nil {
	//	fmt.Printf("连接失败: %v\n", err)
	//	return
	//}
	//defer conn.Close()

	// 发送bind请求
	_, err = conn.Write(data)
	if err != nil {
		log.Printf("发送消息失败: %v\n", err)
		return
	}
	log.Println("Bind消息发送成功")

	// 接收响应
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("接收响应失败: %v\n", err)
		return
	}

	// 解析响应
	log.Printf("接收到 BIND 响应: %x\n", buf[:n])
	// 打印状态码
	log.Println("BIND状态码:", buf[20])
	if len(UserNumber) > 0 {
		for _, UserNumber := range UserNumber {
			// 在 BIND 成功后，发送 SUBMIT 请求
			err = SGIPSubmitSend(ContentSType, conn, UserNumber, Message, File, Content, SPNumber, ChargeNumber, CorpID, ServiceType)
			//fmt.Printf("%v\n", req.Header.SequenceID)
			if err != nil {
				log.Printf("%s发送 SUBMIT 请求失败: %v\n", UserNumber, err)
				return
			}
		}
	}

	if err = SGIPUnbindSend(conn); err != nil {
		log.Printf("SGIPUnbindSend: %v\n", err)
	}
	// 断开连接（defer conn.Close() 会自动关闭连接）
	log.Println("所有操作完成，连接已断开")
}
