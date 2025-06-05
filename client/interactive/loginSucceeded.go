package interactive

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// 登录成功的
func loginSucceeded(conn *websocket.Conn, user string) {
	fmt.Println("已成功连接到服务器，输入[exit]退出")

	// 处理系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go handleSystemSignals(sigChan, conn)

	// 启动读取消息的goroutine
	messageChan := make(chan string)
	go readServerMessages(conn, messageChan)

	// 初始化基础数据
	baseData := clientToolData{
		Username: user,
	}

	// 用户输入处理
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("[VIGILIX]> ")

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		if input == "exit" {
			break
		}

		// 复制基础数据并设置当前命令
		dataToSend := baseData
		dataToSend.Command = input

		// 将结构体转换为JSON
		messageJSON, err := json.Marshal(dataToSend)
		if err != nil {
			fmt.Printf("转换JSON失败: %v\n", err)
			fmt.Print("[VIGILIX]> ")
			continue
		}

		// 发送JSON消息到服务器
		if err := conn.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
			fmt.Printf("发送消息失败: %v\n", err)
			break
		}
		// 等待并显示服务器响应
		response := <-messageChan
		var msg ServerMessage
		msg, err = JsonToStruct(response)
		if err != nil {
			fmt.Printf("< %v\n", err)
		}
		fmt.Println("< " + msg.Content)
		fmt.Print("[VIGILIX]> ")
	}
}
