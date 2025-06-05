package Http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"time"
)

// sendMessageToClient 封装发送消息给客户端的逻辑
func sendMessageToClient(conn *websocket.Conn, Type int, Content string) bool {
	writeMu.Lock()         // 加锁
	defer writeMu.Unlock() // 解锁
	msg := Message{
		Type:    Type, // 自定义一个错误消息类型
		Content: Content,
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error: encoding JSON:", err)
		return false
	}
	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Println("Error: sending message to interactive:", err)
		return false
	}
	return true
}

// handleOrdinaryNotices 发送普通消息
func handleOrdinaryNotices(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数（如通知内容）
	content := r.URL.Query().Get("content")
	if content == "" {
		respondError(w, http.StatusBadRequest, "缺少通知内容参数（content）")
		return
	}

	// 构造普通通知消息
	msg := Message{
		Type:    OrdinaryNotices, // 类型3
		Time:    time.Now().Format(time.RFC3339),
		Content: content,
	}

	// 发送消息到所有客户端
	sendMessageToAllClients(msg)
	respondSuccess(w, "普通通知已发送")
}

// 定向发送函数（保持原广播功能）
func sendToClients(typeID int, tableStatus bool, deliverTime string, clientIDs []int) bool {
	if tableStatus {
		collectMutex.Lock()
		defer collectMutex.Unlock()
		// 初始化新批次
		//currentBatchID = fmt.Sprintf("batch_%d", time.Now().Unix())
		batchFile = excelize.NewFile()
		batchRow = 0
	}
	msg := Message{
		Type:    typeID, // 自定义消息类型
		Time:    deliverTime,
		Content: "CollectInformation.",
	}
	jsonData, _ := json.Marshal(msg)
	success := false
	for _, id := range clientIDs {
		if conn, ok := clientConnections.Load(id); ok {
			err := conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, jsonData)
			if err == nil {
				success = true // 至少一个发送成功
			} else {
				log.Printf("Error: ID %d 发送失败: %v", id, err)
				// 可选择移除无效连接（若确认连接已断开）
				clientConnections.Delete(id)
				continue
			}
		} else {
			log.Printf("Error: ID %d 未找到对应的连接", id)
		}
	}
	return success
}

// manualPushHandler 全量采集
func manualPushHandler(w http.ResponseWriter, r *http.Request, typeID int, tableStatus bool, deliverTime string) {
	if tableStatus {
		collectMutex.Lock()
		defer collectMutex.Unlock()
		// 初始化新批次
		//currentBatchID = fmt.Sprintf("batch_%d", time.Now().Unix())
		batchFile = excelize.NewFile()
		batchRow = 0
	}
	msg := Message{
		Type:    typeID, // 自定义消息类型
		Time:    deliverTime,
		Content: "CollectInformation.",
	}
	sendMessageToAllClients(msg)
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprint(w, "Message sent to all clients.")
	if err != nil {
		log.Printf("Error: writing message to all clients: %v", err)
	}
}
