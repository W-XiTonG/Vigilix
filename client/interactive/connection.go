package interactive

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

// 客户端工具连接登录
type clientToolData struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Command  string `json:"Command"`
}

// ServerMessage Message 定义一个示例的 JSON 数据结构体
type ServerMessage struct {
	Type    int    `json:"type"`
	Time    string `json:"time"`
	Content string `json:"content"`
}

func Connect(ip, username, password string, port int) {
	serverURL := "ws://" + ip + ":" + strconv.Itoa(port) + "/ClientTool"
	// 创建WebSocket连接
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatalf("无法连接到服务器: %v", err)
	}
	defer func() {
		if err = conn.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 发送登录指令
	sendClientLogin(conn, username, password)
	// 接收数据
	returnValue := receive(conn)
	if returnValue == 4 {
		loginSucceeded(conn, username)
	}
	fmt.Println("\n客户端已退出")
}

// JsonToStruct 将JSON字符串转换为结构体的函数
func JsonToStruct(jsonStr string) (ServerMessage, error) {
	var p ServerMessage
	err := json.Unmarshal([]byte(jsonStr), &p)
	if err != nil {
		return p, err
	}
	return p, nil
}
