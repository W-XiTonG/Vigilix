package alarm

import (
	"agent/util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

var (
	droppedAlertCount   uint64 // 被丢弃告警计数器
	processedAlertCount uint64 // 已处理告警计数器
	wg                  sync.WaitGroup
)

// 全局告警状态（使用线程安全的map）
var alertStatus = struct {
	sync.RWMutex
	m map[string]bool
}{m: make(map[string]bool)}

func ThresholdAlarm(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, alarmPush, clientID int, checkInterval time.Duration, queue int32,
	networkCard, authenticationKey string, threshold float64) {
	// 使用传入的 ctx 直接控制生命周期
	defer log.Println("[Alarm]告警监控退出")
	if conn == nil {
		log.Printf("[Alarm] websocket 连接为空,告警未启动")
	}
	// 子上下文用于协调退出
	//ctx, cancel := context.WithCancel(ctx)
	//defer cancel()
	var alertChan = make(chan Alert, queue)
	// 确保通道只关闭一次
	var once sync.Once
	closeAlertChan := func() {
		once.Do(func() {
			close(alertChan)
			log.Println("[Alarm]告警通道已安全关闭")
		})
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Alarm] ThresholdAlarm recovered: %v", r)
		}
	}()
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	ip, err := util.GetIPByInterfaceName(networkCard)
	if err != nil {
		fmt.Printf("[Alarm] 获取 %s 接口的 IP 地址时出错: %v\n", networkCard, err)
		return
	}
	log.Printf("[Alarm] 监控成功开启")

	// 启动一个 goroutine 来处理告警发送
	wg.Add(1) // 在goroutine启动前增加计数器
	go func() {
		defer wg.Done() // 确保最终释放 WaitGroup
		for {
			select {
			case alert, ok := <-alertChan:
				if !ok { // 通道已关闭
					return
				}
				sendAlert(conn, alert, clientID, alarmPush, authenticationKey)
			case <-ctx.Done(): // 监听上下文取消
				log.Println("[Alarm] 告警处理协程收到退出信号，停止处理")
				return
			}
		}
	}()
	//
	//// 创建一个信号通道，用于接收系统信号
	//sigs := make(chan os.Signal, 1)
	//// 监听 SIGINT 和 SIGTERM 信号
	//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//
	//// 启动一个 goroutine 来处理接收到的信号
	//go func() {
	//	// 阻塞等待信号
	//	sig := <-sigs
	//	log.Printf("[Alarm]接收到信号 %v，准备退出...", sig)
	//	// 接收到信号后，关闭退出通道，通知主循环退出
	//	cancel()
	//	// close(alertChan)
	//	closeAlertChan()
	//}()
	// 给主循环添加标签
mainLoop:
	// 主循环，使用 select 语句监听 ticker.C 和 ctx 通道
	for {
		select {
		case <-ticker.C:
			// 定时器触发，执行指标采集和告警检查
			metrics, err := collectMetrics()
			if err != nil {
				log.Printf("[Alarm] 采集指标失败: %v", err)
				continue
			}
			checkAndSendAlert(alertChan, ip, threshold, metrics)
		case <-ctx.Done():
			// 接收到退出信号，打印日志并跳出整个主循环
			log.Println("[Alarm] 收到上下文取消信号，停止采集...")
			closeAlertChan()
			break mainLoop
		}
	}

	// 优雅关闭处理
	// 关闭告警通道，停止接收新的告警
	//close(alertChan)
	log.Println("[Alarm] 等待剩余告警发送...")

	// 创建一个通道用于通知等待完成
	done := make(chan struct{})
	// 启动一个 goroutine 等待所有告警发送完成
	go func() {
		wg.Wait() // 等待所有进行中的发送完成
		close(done)
	}()

	// 等待告警发送完成或超时
	select {
	case <-done:
		log.Println("[Alarm] 所有告警已处理")
	case <-time.After(3 * time.Second):
		log.Println("[Alarm] 等待超时，强制退出")
	}
}

// 采集性能指标
func collectMetrics() (*SystemMetrics, error) {
	hostInfo, _ := host.Info()

	// 获取CPU使用率（1秒内的平均使用率）
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("[Alarm]获取CPU使用率失败: %v", err)
	}

	// 获取内存使用率
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("[Alarm]获取内存信息失败: %v", err)
	}

	// 获取磁盘使用率
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, fmt.Errorf("[Alarm]获取磁盘分区信息失败: %v", err)
	}
	diskUsage := make(map[string]float64)
	diskInode := make(map[string]float64)
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			log.Printf("[Alarm]获取 %s 磁盘使用率失败: %v", partition.Mountpoint, err)
			continue
		}
		diskUsage[partition.Mountpoint] = usage.UsedPercent
		operatingSystem := runtime.GOOS
		if operatingSystem == "linux" {
			diskInode[partition.Mountpoint] = usage.InodesUsedPercent
		}
	}

	return &SystemMetrics{
		HostName:   hostInfo.Hostname,
		CPUUsage:   cpuPercent[0],
		MemUsage:   memInfo.UsedPercent,
		DiskUsage:  diskUsage,
		ReportTime: time.Now().Format(time.RFC3339),
		DiskInode:  diskInode,
	}, nil
}

func enqueueAlert(alertChan chan<- Alert, alert Alert) {
	select {
	case alertChan <- alert:
		atomic.AddUint64(&processedAlertCount, 1)
		log.Printf("[Alarm - SUCCESS] 告警入队: %s %s", alert.MetricType, alert.Partition)
	default:
		atomic.AddUint64(&droppedAlertCount, 1)
		alertJSON, _ := json.Marshal(alert)
		log.Printf("[Alarm - WARNING] 队列已满丢弃告警: %s", string(alertJSON))
	}
}
