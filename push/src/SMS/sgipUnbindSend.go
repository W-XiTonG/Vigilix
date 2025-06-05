package SMS

import (
	"log"
	"net"
)

// SGIPUnbindSend unbind
func SGIPUnbindSend(conn net.Conn) error {
	// 创建 SGIP SUBMIT 请求消息
	req := SGIPUnbind{
		Header: SGIPHeader{
			CommandID:  SGIP_UNBIND,           // SUBMIT 请求命令 ID
			SequenceID: generateSequenceID(3), // 生成 SequenceID
		},
	}

	// 编码请求
	data, err := req.Encode()
	if err != nil {
		log.Printf("编码失败: %v\n", err)
		return err
	}
	// 发送UNBIND
	_, err = conn.Write(data)
	if err != nil {
		log.Printf("发送 UNBIND 请求失败: %v\n", err)
		return err
	}
	log.Println("UNBIND 请求发送成功")
	// 接收响应
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("接收响应失败: %v\n", err)
		return err
	}
	log.Printf("接收到 UNBIND 响应: %x\n", buf[:n])
	return nil
}
