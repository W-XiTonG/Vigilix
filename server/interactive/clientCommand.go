package Http

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"server/handlers"
	"strings"
)

// 独立处理已认证请求
func handleAuthenticatedRequest(conn *websocket.Conn, msg clientToolData, pushAddr, user string, pushStatus bool) {
	msgData := strings.Fields(msg.Command)
	log.Printf("Client -> Server %v", msg)
	switch true {
	case msg.Command == "do_something":
		sendMessageToClient(conn, CommandResult, "执行成功")
		log.Printf("Server -> Client 执行成功 Command: %s", msg.Command)
	case strings.HasPrefix(msg.Command, "push"):
		if pushStatus {
			if len(msgData) > 1 {
				pushData := strings.Join(msgData[1:], " ")
				// 使用map创建动态JSON结构
				//data := map[string]interface{}{
				//	user: pushData, // 使用传递进来的user作为键
				//}
				err := handlers.Push(pushData, pushAddr, 1)
				if err != nil {
					sendMessageToClient(conn, InvalidCommand, "执行失败")
					log.Printf("Server -> ClientTool: 执行失败 error: %v Command: %s", err, msg.Command)
				} else {
					sendMessageToClient(conn, CommandResult, "执行成功")
					log.Printf("Server -> ClientTool: 执行成功 Command: %s", msg.Command)
				}
			} else {
				sendMessageToClient(conn, InvalidCommand, "push命令用法: push <data>")
				log.Printf("Server -> ClientTool: 执行失败 Command: %s", msg.Command)
			}
		} else {
			revert := fmt.Sprintf("命令异常: Error status(%d)", ErrorCommand)
			sendMessageToClient(conn, ErrorCommand, revert)
		}
	default:
		sendMessageToClient(conn, InvalidCommand, "未知命令")
		log.Printf("Server -> ClientTool: 执行失败 Command: %s", msg.Command)
	}
}
