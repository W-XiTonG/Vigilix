package interactive

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
)

// 连接关闭
func readServerMessages(conn *websocket.Conn, messageChan chan<- string) {
	defer close(messageChan) // 确保通道关闭
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Fatalf("连接关闭: %v\n", err)
			}
			return
		}
		messageChan <- string(message)
	}
}

// 退出
func handleSystemSignals(sigChan <-chan os.Signal, conn *websocket.Conn) {
	<-sigChan
	fmt.Println("\n收到退出信号，关闭连接...")
	err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		fmt.Println(err)
	}
	if err = conn.Close(); err != nil {
		fmt.Println(err)
	}
	os.Exit(0)
}

// 发送登录指令
func sendClientLogin(conn *websocket.Conn, username, password string) {
	msg := clientToolData{
		Username: username,
		Password: password,
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("JSON encode error: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		fmt.Printf("failed to send ID: %v", err)
	} else {
		fmt.Println("successfully sent ID")
	}
}
