package alarm

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

// 发送告警到Push服务
func sendAlert(conn *websocket.Conn, alert Alert, clientID int, typeID int, authenticationKey string) {
	msg := pushClientInfo{
		ClientID: clientID,
		Type:     typeID,
		Key:      authenticationKey,
		Data:     alert,
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Client %d JSON encode error: %v", clientID, err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		log.Printf("Client %d failed to send ID: %v", clientID, err)
	} else {
		log.Printf("Client successfully sent ID:%d,typeID: %d,content:%v", clientID, typeID, alert)
	}
}
