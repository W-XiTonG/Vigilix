package Http

import (
	"encoding/json"
	"sync"
	"time"
)

// AsyncConfig 异步处理器配置
type AsyncConfig struct {
	MaxWorkers    int           // 最大并发工作协程数
	QueueSize     int           // 任务队列容量
	WorkerTimeout time.Duration // 单个任务超时时间
}

// AsyncProcessor 异步处理器对象
type AsyncProcessor struct {
	taskQueue   chan string    // 任务队列
	stopChan    chan struct{}  // 停止信号
	wg          sync.WaitGroup // 协程同步
	pushHandler func(string)   // 实际推送函数
	asyncConfig AsyncConfig    // 配置参数
	stopOnce    sync.Once
}

// 自定义消息
type customData struct {
	Type    int             `json:"Type"`
	Content json.RawMessage `json:"Content"`
}

// UnixTime 自定义时间类型
type UnixTime struct {
	time.Time
}
