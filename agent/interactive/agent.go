package interactive

import (
	"agent/alarm"
	"agent/collectors"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	// 连接服务端
	connect = 101
	// 推送数据
	data = 102
	// 告警推送
	alarmPush = 103
)

// pushClientInfo 定义客户端消息结构体
type pushClientInfo struct {
	ClientID int                       `json:"client_id"`
	Type     int                       `json:"type"`
	Key      string                    `json:"key"`
	Time     string                    `json:"deliverTime"`
	Data     *collectors.SystemMetrics `json:"data"`
}

// Message 定义一个示例的 JSON 数据结构体
type Message struct {
	Type          int    `json:"type"`
	Time          string `json:"time"`
	DiskPartition string `json:"disk_partition"`
	Content       string `json:"content"`
}

func Client(ctx context.Context, cancel context.CancelFunc, clientID int, serverAddr, AuthenticationKey,
	networkCard string, retryInterval time.Duration,
	alarmStatus bool, checkInterval time.Duration, queue int32, threshold float64) {
	wsAddr := serverAddr + "/login"
	// 添加告警控制变量
	var alarmCancelMutex sync.Mutex
	var alarmCancel context.CancelFunc

	// 外层循环用于连接断开后重连
	for {
		select {
		case <-ctx.Done(): // 退出点
			log.Printf("Client %d received shutdown signal", clientID)
			return
		default:
			var conn *websocket.Conn
			var err error

			// 内层循环处理首次连接和重试
		retryLoop:
			for {
				select {
				case <-ctx.Done():
					if conn != nil {
						err = conn.Close()
						if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
							log.Printf("Client %d 关闭连接错误: %v", clientID, err)
						}
					}
					// 关闭时终止告警
					alarmCancelMutex.Lock()
					if alarmCancel != nil {
						alarmCancel()
						alarmCancel = nil
					}
					alarmCancelMutex.Unlock()
					time.Sleep(retryInterval)
					return
				default:
					log.Printf("Client %d attempting to connect to %s", clientID, wsAddr)
					conn, _, err = websocket.DefaultDialer.Dial(wsAddr, nil)
					if err != nil {
						log.Printf("Connection failed: %v, retrying in %v", err, retryInterval)
						select {
						case <-ctx.Done():
							if conn != nil {
								func() {
									if err = conn.Close(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
										log.Printf("%v", err)
									}
								}()
							}
							return
						case <-time.After(retryInterval):
						}
						continue
					}
					log.Printf("Client %d connected successfully", clientID)
					break retryLoop
				}
			}
			//// 创建子 context（控制本次连接的生命周期）
			//connCtx, cancelConn := context.WithCancel(ctx)
			// 创建连接级context
			connCtx, connCancel := context.WithCancel(ctx)

			// 连接成功后立即发送ID
			if conn != nil {
				sendClientID(conn, clientID, connect, AuthenticationKey)
			}

			// 创建断开通知通道
			disconnected := make(chan struct{})

			// 启动一个协程监听 context 取消
			go func(conn *websocket.Conn) {
				<-ctx.Done()
				log.Printf("Client %d 正在关闭连接...", clientID)

				// 等待消息处理协程退出（通过 disconnected 通道）
				select {
				case <-disconnected:
					log.Printf("Client %d 消息处理协程已退出", clientID)
				case <-time.After(1 * time.Second):
					log.Printf("Client %d 强制关闭连接", clientID)
				}
				// 关闭时终止告警
				alarmCancelMutex.Lock()
				if alarmCancel != nil {
					alarmCancel()
					alarmCancel = nil
				}
				alarmCancelMutex.Unlock()
				// 关闭连接
				if conn != nil {
					err = conn.Close()
					if err != nil {
						if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
							log.Printf("Client %d 关闭连接错误: %v", clientID, err)
						}
					}
				}
			}(conn) // 传递当前连接的副本
			// 消息接收协程（添加重连触发机制）
			go func() {
				defer close(disconnected) // 关闭通道通知连接断开
				defer func() {
					//<-ctx.Done()
					alarmCancelMutex.Lock()
					if alarmCancel != nil {
						alarmCancel()
						alarmCancel = nil
					}
					alarmCancelMutex.Unlock()
					if conn != nil {
						if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
							log.Printf("Client %d 关闭连接错误: %v", clientID, err)
						}
					}
				}() // 确保连接关闭

				for {
					_, message, err := conn.ReadMessage()
					if err != nil {
						log.Printf("Client %d connection lost: %v", clientID, err)
						return
					}

					var msg Message
					if err := json.Unmarshal(message, &msg); err != nil {
						log.Printf("Error unmarshalling message: %v", err)
						continue
					}
					//log.Printf("Server -> Client : %v ", msg)
					switch msg.Type {
					case 0:
						log.Printf("Error: Authentication failed: clientId %d", clientID)
						time.Sleep(retryInterval)
					case 1:
						// 服务端心跳检测
						//log.Printf("msgType %d: HeartbeatDetection", msg.Type)
					case 2:
						sendData(conn, clientID, msg.Time, data, collectors.Metrics(msg.DiskPartition), AuthenticationKey)
					case 3:
						log.Printf("interactive received: %s", string(message))
					case 4:
						// 鉴权成功
						if alarmStatus {
							// 终止之前的告警
							alarmCancelMutex.Lock()
							if alarmCancel != nil {
								alarmCancel()
								alarmCancel = nil
							}
							alarmCancelMutex.Unlock()
							// 使用连接级context
							var alarmCtx context.Context
							alarmCtx, alarmCancel = context.WithCancel(connCtx)
							go alarm.ThresholdAlarm(alarmCtx, cancel, conn, alarmPush, clientID, checkInterval, queue, networkCard, AuthenticationKey, threshold)
						}
					default:
						log.Printf("Client %d: Unknown message type: %d", clientID, msg.Type)
					}
				}
			}()

			// 阻塞等待断开通知
			<-disconnected
			log.Printf("Client %d initiating reconnect...", clientID)
			alarmCancelMutex.Lock()
			if alarmCancel != nil {
				alarmCancel()
				alarmCancel = nil
			}
			alarmCancelMutex.Unlock()
			//cancel() // 连接关闭时释放资源
			connCancel()
		}
	}
}
