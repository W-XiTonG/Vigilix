package Http

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"server/handlers"
)

func clientTool(w http.ResponseWriter, r *http.Request, ClientTool map[string]string, upGrader websocket.Upgrader, pushAddr string, pushStatus bool) {
	// 升级HTTP连接到WebSocket
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("升级连接失败:", err)
		return
	}
	defer func() {
		if err = conn.Close(); err != nil {
			log.Printf("clientTool conn.Close Error: %v", err)
			return
		}
	}()
	log.Printf("ClientTool conn.Close %v:", r.RemoteAddr)
	// 认证状态标志
	var authenticated = false
	for {
		// 读取消息
		_, message, err := conn.ReadMessage()
		if err != nil {
			// 区分正常关闭和异常错误
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("异常关闭: %v", err)
			} else {
				log.Println("连接已正常关闭")
			}
			return
		}
		var clientToolMsg clientToolData
		if err := json.Unmarshal(message, &clientToolMsg); err == nil {
			// 正常处理结构化消息
			log.Printf("ClientTool -> Server Data: %v", clientToolMsg)
			if !authenticated {
				clientUser, clientPass := handlers.ClientPartition(ClientTool, clientToolMsg.Username)
				if clientToolMsg.Username != clientUser || clientPass != clientToolMsg.Password {
					// 鉴权失败，发送消息给客户端
					sendMessageToClient(conn, AuthenticationFailed, "账号或密码错误")
					log.Printf("Server -> ClientTool: AuthenticationFailed")
					return
				}
				authenticated = true // 标记为已认证
				log.Printf("ClientTool Authentication succeeded")
				sendMessageToClient(conn, AuthenticationSuccess, "Authentication succeeded")
				//continue
			} else {
				handleAuthenticatedRequest(conn, clientToolMsg, pushAddr, clientToolMsg.Username, pushStatus)
			}
		} else {
			sendMessageToClient(conn, AuthenticationFailed, "消息格式错误")
			log.Printf("Server -> ClientTool: 消息格式错误")
			continue // 或 return 根据需求
		}
	}
}
