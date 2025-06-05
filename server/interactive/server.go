package Http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/alarm"
	"server/handlers"
	"server/util"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// server类型
const (
	// AuthenticationFailed 剔除用户
	AuthenticationFailed = 0
	// HeartbeatDetection 发送心跳检测
	HeartbeatDetection = 1
	// CollectInformation 手动推送路由，后续用作采集
	CollectInformation = 2
	// OrdinaryNotices 普通推送通知
	OrdinaryNotices = 3
	// AuthenticationSuccess 鉴权成功
	AuthenticationSuccess = 4
	// CommandResult ClientTool命令执行成功
	CommandResult = 5
	// InvalidCommand ClientTool未知命令
	InvalidCommand = 6
	// ErrorCommand ClientTool失败命令
	ErrorCommand = 7
)

// interactive 类型
const (
	// 客户端断开连接
	clientDisconnect = 100
	// 客户端连接
	clientConnect = 101
	// 客户端推送数据
	clientData = 102
	// 客户端发生告警
	clientAlarm = 103
	// 客户端连接工具
	clientTools = 104
)

// 广播通道缓冲大小
const broadcastBufferSize = 256

var (
	collectMutex      sync.Mutex
	batchFile         *excelize.File                            // 新增，用于存储创建的 Excel 文件
	batchRow          int                                       // 新增，用于记录当前写入的行号
	broadcastChan     = make(chan Message, broadcastBufferSize) // 广播消息通道
	clients           sync.Map                                  // 替换原有map为sync.Map
	clientConnections sync.Map                                  // key: ClientID(int), value: *websocket.Conn
	clientIDs         sync.Map                                  // key: *websocket.Conn, value: ClientID(int)
	clientMutex       sync.Mutex
	clientConnCount   = make(map[string]int) // 记录客户端连接次数
	count             int
	writeMu           sync.Mutex
)

func handleWebSocket(w http.ResponseWriter, r *http.Request, DetectionTime time.Duration, upGrader websocket.Upgrader,
	tableStatus, PushStatus, AgentAuthenticationStatus,
	DetectionStatus bool, tablePaths, PushAddr string,
	DetectionFrequency int, agentsID map[int][]string) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection to WebSocket:", err)
		return
	}

	log.Printf("Client connected: %s", r.RemoteAddr)
	// 添加客户端连接到列表
	//clients[conn] = true
	// 连接时存储
	clients.Store(conn, true)
	defer func() {
		// 连接关闭时从列表中移除
		// 使用 sync.Map 删除连接
		clients.Delete(conn)
		if clientID, ok := clientIDs.Load(conn); ok {
			clientConnections.Delete(clientID)
			clientIDs.Delete(conn)
			if PushStatus {
				clientsID, ok := util.ConvertToInt(clientID)
				if !ok {
					log.Println("Client ID not found in clients map")
				}
				handlers.PushServer(clientDisconnect, r.RemoteAddr, clientsID, "断开连接", PushAddr, agentsID, nil)
			}
		}

		if err = conn.Close(); err != nil {
			log.Println("Error closing WebSocket connection:", err)
		}
	}()

	// 持续接收客户端发送的 JSON 数据
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from interactive:", err)
			return
		}

		var clientMsg Agent
		if err := json.Unmarshal(message, &clientMsg); err == nil {
			// 正常处理结构化消息
			log.Printf("Received structured message: Type:%d,ClientID:%d,DeliverTime:%s,Data:%v",
				clientMsg.Type, clientMsg.ClientID, clientMsg.DeliverTime, clientMsg.Data)
			// 鉴权ID是否在server配置
			clientID := clientMsg.ClientID
			key, ok := handlers.AgentKey(agentsID, clientID)
			if !ok {
				// 鉴权失败，发送消息给客户端
				sendMessageToClient(conn, AuthenticationFailed, "Authentication failed")
				log.Println("Authentication failed, closing connection: Id NO")
				clients.Delete(conn)
				if clientID, ok := clientIDs.Load(conn); ok {
					clientConnections.Delete(clientID)
					clientIDs.Delete(conn)
					if PushStatus {
						clientsID, ok := util.ConvertToInt(clientID)
						if !ok {
							log.Println("Client ID not found in clients map")
						}
						handlers.PushServer(clientDisconnect, r.RemoteAddr, clientsID, "断开连接", PushAddr, agentsID, nil)
					}
				}
				if err = conn.Close(); err != nil {
					log.Println("Error closing WebSocket connection:", err)
				}
				return
			} else {
				if AgentAuthenticationStatus {
					if clientMsg.Key != key {
						// 鉴权失败，发送消息给客户端
						sendMessageToClient(conn, AuthenticationFailed, "Authentication failed")
						log.Println("Authentication failed, closing connection: Key NO")
						clients.Delete(conn)
						if clientID, ok := clientIDs.Load(conn); ok {
							clientConnections.Delete(clientID)
							clientIDs.Delete(conn)
							if PushStatus {
								clientsID, ok := util.ConvertToInt(clientID)
								if !ok {
									log.Println("Client ID not found in clients map")
								}
								handlers.PushServer(clientDisconnect, r.RemoteAddr, clientsID, "断开连接", PushAddr, agentsID, nil)
							}
						}
						if err = conn.Close(); err != nil {
							log.Println("Error closing WebSocket connection:", err)
						}
						return
					}
				}
				log.Printf("Authentication succeeded:%s", key)
				// 获取客户端IP（去除端口）
				clientIP := util.ExtractIP(r.RemoteAddr)
				// 拼接IP加ID
				reClientIP := fmt.Sprintf("%s:%d", clientIP, clientMsg.ClientID)
				// 判断连接类型
				clientMutex.Lock()
				count = clientConnCount[reClientIP]
				clientConnCount[reClientIP] = count + 1 // 递增次数
				clientMutex.Unlock()
				switch clientMsg.Type {
				case clientConnect:
					// 存储映射关系
					clientConnections.Store(clientMsg.ClientID, conn)
					clientIDs.Store(conn, clientMsg.ClientID)
					if PushStatus {
						if count >= 1 {
							handlers.PushServer(clientMsg.Type, r.RemoteAddr, clientMsg.ClientID, "恢复连接",
								PushAddr, agentsID, nil)
						}
					}
					//sendMessageToClient(conn, AuthenticationSuccess, "Authentication Success")
					// 定时向客户端发送 JSON 数据
					if DetectionStatus {
						heartbeatDone := startHeartbeatCheck(conn, DetectionTime, DetectionFrequency)
						go func() {
							<-heartbeatDone
							// 执行清理操作
							if clientID, ok := clientIDs.Load(conn); ok {
								clientConnections.Delete(clientID)
								clientIDs.Delete(conn)
								if PushStatus {
									clientsID, ok := util.ConvertToInt(clientID)
									if !ok {
										log.Println("Client ID not found in clients map")
									}
									handlers.PushServer(clientDisconnect, r.RemoteAddr, clientsID, "断开连接", PushAddr, agentsID, nil)
								}
							}
							clients.Delete(conn)
							if err = conn.Close(); err != nil {
								log.Println("Error clients closing WebSocket connection:", err)
							}
							log.Printf("连接关闭: %s", conn.RemoteAddr())
						}()
						log.Printf("Server -> Agent : Detection Start succeeded")
					}
					sendMessageToClient(conn, AuthenticationSuccess, "Authentication succeeded")
				case clientData:
					if tableStatus {
						collectMutex.Lock()
						// 双重检查 batchFile 有效性
						if batchFile == nil {
							log.Println("batchFile 未初始化，丢弃数据")
							break
						}
						//if _, ok := clients[conn]; ok {
						if _, ok := clients.Load(conn); ok {
							if batchFile == nil {
								batchFile = excelize.NewFile()
								batchRow = 0
							}
							batchRow++
							diskPartition, ok := handlers.AgentPartition(agentsID, clientMsg.ClientID)
							if !ok {
								log.Printf("Error: Disk partition value not found: %s", diskPartition)
							}
							clientDiskPartition, ok := util.FindDiskPartition(clientMsg.Data.Disk.Partitions,
								diskPartition)
							if !ok {
								log.Printf("Error: Partition not found: %v", clientDiskPartition)
							}
							agentName := handlers.AgentName(agentsID, clientMsg.ClientID)
							handlers.TableMaking(batchFile, tablePaths, clientMsg.DeliverTime, agentName, clientMsg.Data,
								batchRow, clientDiskPartition)
						}
						collectMutex.Unlock()
					}
				case clientAlarm:
					var agentsAlarm agentAlarm
					err = json.Unmarshal(message, &agentsAlarm)
					sendMessageToClient(conn, OrdinaryNotices, "Successfully received the alarm")
					if PushStatus {
						if count >= 1 {
							agentsAlarmData := handlers.Alarm{
								HostName:   agentsAlarm.Data.HostName,
								IpAddr:     agentsAlarm.Data.IpAddr,
								AlertTime:  agentsAlarm.Data.AlertTime,
								Message:    agentsAlarm.Data.Message,
								MetricType: agentsAlarm.Data.MetricType,
								Current:    agentsAlarm.Data.Current,
								Partition:  agentsAlarm.Data.Partition,
								State:      agentsAlarm.Data.State,
							}
							handlers.PushServer(clientAlarm, r.RemoteAddr, agentsAlarm.ClientID, "恢复连接",
								PushAddr, agentsID, &agentsAlarmData)
						}
					}
				case clientTools:
					sendMessageToClient(conn, clientTools, "clientTools succeeded")

				}
			}
		} else {
			// 处理纯文本消息（客户端ID）
			log.Printf("Received interactive : %s", string(message))
		}
		if err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}
		log.Printf("Received JSON message from interactive: %+v", clientMsg)
	}
}

func sendMessageToAllClients(msg Message) {
	select {
	case broadcastChan <- msg: // 非阻塞写入通道
	default:
		log.Println("Warning: The broadcast queue is full, and the message type is dropped", msg.Type)
	}
}

func Server(port, DetectionFrequency int, DetectionTime, zabbixIntervalMin, zabbixIntervalMax time.Duration,
	upGrader websocket.Upgrader, tableStatus, PushStatus, AgentAuthenticationStatus,
	DetectionStatus, zabbixAlarmStatus, zabbixDebug bool, tablePaths, PushAddr, zabbixUrl, zabbixUser, zabbixPass string,
	agentIDs map[int][]string, ClientToolStatus bool, ClientTool map[string]string) {
	if zabbixAlarmStatus {
		go alarm.ZabbixAlarm(PushAddr, zabbixUrl, zabbixUser, zabbixPass, zabbixIntervalMin, zabbixIntervalMax, zabbixDebug)
	}
	// 启动广播协程 (只需启动一次)
	go func() {
		for msg := range broadcastChan {
			jsonData, err := json.Marshal(msg)
			if err != nil {
				log.Println("Error: Broadcast message encoding failed:", err)
				continue
			}
			// 遍历所有客户端发送
			clients.Range(func(key, _ interface{}) bool {
				conn := key.(*websocket.Conn)
				if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
					log.Printf("Error: interactive %s Failed to send: %v", conn.RemoteAddr(), err)
					if err = conn.Close(); err != nil {
						log.Println("Error: closing WebSocket connection:", err)
					}
					clients.Delete(conn)
					if clientID, ok := clientIDs.Load(conn); ok {
						clientConnections.Delete(clientID)
						clientIDs.Delete(conn)
						if PushStatus {
							clientsID, ok := util.ConvertToInt(clientID)
							if !ok {
								log.Println("Error: Client ID not found in clients map")
							}
							handlers.PushServer(clientDisconnect, "未知", clientsID, "断开连接",
								PushAddr, agentIDs, nil)
						}
					}
				}
				return true // 继续遍历
			})
		}
	}()
	// 创建路由器
	mux := http.NewServeMux()

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, DetectionTime, upGrader, tableStatus, PushStatus,
			AgentAuthenticationStatus, DetectionStatus, tablePaths,
			PushAddr, DetectionFrequency, agentIDs)
	})
	// 触发普通通知
	mux.HandleFunc("/OrdinaryNotices", handleOrdinaryNotices)
	// 添加手动推送的路由
	mux.HandleFunc("/CollectInformation", func(w http.ResponseWriter, r *http.Request) {
		// 记录原始请求参数
		log.Printf("请求参数: %+v", r.URL.Query())
		// 新增解析clientID参数
		ids := r.URL.Query()["clientID"]
		var targetIDs []int
		var invalidIDs []string // 记录非法ID
		for _, sid := range ids {
			id, err := strconv.Atoi(strings.TrimSpace(sid)) // 增加去除空白字符
			if err != nil {
				invalidIDs = append(invalidIDs, sid)
				continue
			}
			if id <= 0 { // 添加ID有效性验证
				invalidIDs = append(invalidIDs, sid)
				continue
			}
			targetIDs = append(targetIDs, id)
		}
		// 处理非法ID情况
		if len(invalidIDs) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, err := fmt.Fprintf(w, "包含非法ID: %v", invalidIDs)
			if err != nil {
				log.Printf("Fprintf: Error writing response: %v", err)
			}
			log.Printf("非法ID拒绝: %v", invalidIDs)
			return
		}
		// 获取当前时间
		currentTime := time.Now()
		// 按照指定格式格式化时间，精确到分钟
		formattedTime := currentTime.Format("20060102150405.000")
		if len(targetIDs) > 0 {
			log.Printf("开始定向采集到IDs: %v", targetIDs)
			// 定向发送
			if ok := sendToClients(CollectInformation, tableStatus, formattedTime, targetIDs); ok {
				// 返回 200 状态码表示成功
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("定向采集成功"))
				if err != nil {
					log.Printf("Error writing response: %v", err)
				}
				log.Printf("Return value:%d", http.StatusOK)
			} else {
				respondError(w, http.StatusInternalServerError, "定向采集失败")
			}
		} else {
			log.Println("执行全量采集")
			manualPushHandler(w, r, CollectInformation, tableStatus, formattedTime)
			// 返回 200 状态码表示成功
			respondSuccess(w, "全量采集推送成功")
			log.Printf("Return value:%d", http.StatusOK)
		}
	})
	mux.HandleFunc("/ClientTool", func(w http.ResponseWriter, r *http.Request) {
		if ClientToolStatus {
			clientTool(w, r, ClientTool, upGrader, PushAddr, PushStatus)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	// 应用全局 CORS
	handler := corsMiddleware(mux)

	ServerPort := fmt.Sprintf(":%d", port)
	//err := http.ListenAndServe(ServerPort, handler)
	//if err != nil {
	//	log.Fatal("Error: starting interactive:", err)
	//}
	// 创建支持优雅关闭的HTTP Server
	srv := &http.Server{
		Addr:    ServerPort,
		Handler: handler,
	}

	// 创建全局取消上下文
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动服务器监听（在goroutine中运行）
	go func() {
		log.Printf("WebSocket interactive started on :%d", port)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// 统一信号处理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待终止信号
	<-quit
	log.Println("Shutting down server...")

	if zabbixAlarmStatus {
		// 先关闭Zabbix告警
		//alarm.Shutdown()
		shutdownChan := alarm.InitSignal()
		<-shutdownChan
	}

	// 然后关闭HTTP服务（设置5秒超时）
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Server gracefully stopped")

}

// 心跳检测函数
func startHeartbeatCheck(conn *websocket.Conn, interval time.Duration, DetectionFrequency int) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		errorCount := 0
		for {
			select {
			case <-ticker.C:
				ok := sendMessageToClient(conn, HeartbeatDetection, "HeartbeatDetection.")
				if !ok {
					errorCount++
					if errorCount >= DetectionFrequency {
						log.Printf("Error: 心跳检测%d次，关闭连接: %s", DetectionFrequency, conn.RemoteAddr())
						return
					}
				} else {
					errorCount = 0
				}
			}
		}
	}()
	return done
}
