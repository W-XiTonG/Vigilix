package interactive

import (
	"agent/collectors"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

// 发送客户端ID（仅在连接/重连成功时调用）
func sendClientID(conn *websocket.Conn, clientID int, typeID int, AuthenticationKey string) {
	msg := pushClientInfo{
		ClientID: clientID,
		Key:      AuthenticationKey,
		Type:     typeID,
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Client %d JSON encode error: %v", clientID, err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		log.Printf("Client %d failed to send ID: %v", clientID, err)
	} else {
		log.Printf("Client %d successfully sent ID", clientID)
	}
}
func sendData(conn *websocket.Conn, clientID int, deliverTime string, typeID int, content *collectors.SystemMetrics, AuthenticationKey string) {
	msg := pushClientInfo{
		ClientID: clientID,
		Type:     typeID,
		Key:      AuthenticationKey,
		Time:     deliverTime,
		Data:     content,
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Client %d JSON encode error: %v", clientID, err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		log.Printf("Client %d failed to send ID: %v", clientID, err)
	} else {
		log.Printf("Client successfully sent ID:%d,typeID: %d,content:%v", clientID, typeID, content)
	}
}
