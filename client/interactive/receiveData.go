package interactive

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

// 处理登录时接收的类型
func receive(conn *websocket.Conn) int {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("connection lost: %v", err)
			return 0
		}

		var msg ServerMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Printf("Error unmarshalling message: %v", err)
			continue
		}
		//log.Printf("Server -> Client : %v ", msg)
		switch msg.Type {
		case 0:
			fmt.Printf("%v", msg.Content)
			return 0
		case 3:
			fmt.Printf("%v", msg.Content)
		case 4:
			return 4
		default:
			fmt.Printf("Unknown message type: %v", msg.Content)
		}
	}
}
