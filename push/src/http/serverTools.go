package Http

import (
	"Push/util"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// NewAsyncProcessor 创建异步处理器
func NewAsyncProcessor(pushHandler func(string), config AsyncConfig) *AsyncProcessor {
	return &AsyncProcessor{
		taskQueue:   make(chan string, config.QueueSize),
		stopChan:    make(chan struct{}),
		pushHandler: pushHandler,
		asyncConfig: config,
	}
}

// Start 启动 worker 池
func (ap *AsyncProcessor) Start() {
	ap.wg.Add(ap.asyncConfig.MaxWorkers)
	for i := 0; i < ap.asyncConfig.MaxWorkers; i++ {
		go ap.worker()
	}
}

// Stop 停止 worker 池
func (ap *AsyncProcessor) Stop() {
	ap.stopOnce.Do(func() { // 确保只执行一次
		close(ap.stopChan)
		ap.wg.Wait()

		// 安全关闭任务队列
		if ap.taskQueue != nil {
			close(ap.taskQueue)
			ap.taskQueue = nil // 防止重复关闭
		}
	})
}

func (ap *AsyncProcessor) worker() {
	defer ap.wg.Done()

	for {
		select {
		case task := <-ap.taskQueue:
			// 带超时控制的推送
			done := make(chan struct{})
			go func() {
				defer close(done)
				ap.pushHandler(task) // 执行实际推送逻辑
			}()

			select {
			case <-done:
			case <-time.After(ap.asyncConfig.WorkerTimeout):
				log.Printf("任务超时: %s", task)
			}

		case <-ap.stopChan:
			return // 收到停止信号退出
		}
	}
}

// Conversion 是否需要换行
func Conversion(LineBreaksStatus bool, jsonStr string, LineBreaks string) string {
	// 将请求体内容转换为字符串
	//jsonStr := string(body)
	if LineBreaksStatus {
		Content := util.AddLineBreaks(jsonStr, LineBreaks)
		return Content
	}
	return jsonStr
}

// formatAlertMessage 格式化告警消息为人类可读的字符串
func formatAlertMessage(hostname, ip, time, message, metricType, state, partition string, current float64) string {
	return fmt.Sprintf("[告警信息]\n主机: %s\nIP: %s\n时间: %s\n描述: %s\n类型: %s\n当前值: %.2f\n分区: %s\n状态: %s",
		hostname, ip, time, message, metricType, current, partition, state)
}

// 处理agent状态信息
func formatStatusMessage(ClientID int, ClientName, ClientAddr, ClientStatus string) string {
	return fmt.Sprintf("[告警信息]\nAgentID: %d\n描述: %s\n地址: %s\n状态: %s",
		ClientID, ClientName, ClientAddr, ClientStatus)
}

// 处理zabbix告警信息
func formatZabbixStatusMessage(EventID, TriggerID, Description, Status, Severity, Acknowledged string, Hosts []string, Timestamp UnixTime) string {
	return fmt.Sprintf("[告警信息]\n事件ID: %s\n触发器ID: %s\n触发器描述: %s\n主机: %s\n状态: %s\n严重等级: %s\n事件时间戳: %v\n是否已确认: %s",
		EventID, TriggerID, Description, Hosts, Status, Severity, Timestamp, Acknowledged)
}

// UnmarshalJSON 时间解析逻辑
func (u *UnixTime) UnmarshalJSON(b []byte) error {
	// 处理可能的引号包裹
	s := strings.Trim(string(b), `"`)

	// 转换为int64
	t, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("解析时间戳失败: %w，原始值: %s", err, s)
	}

	// Zabbix时间戳为秒级
	u.Time = time.Unix(t, 0)
	return nil
}

// MarshalJSON 添加MarshalJSON方法
func (u UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", u.Unix())), nil
}
