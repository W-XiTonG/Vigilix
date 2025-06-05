package alarm

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"
)

// ZabbixClient 结构体定义
type ZabbixClient struct {
	authToken string
	mutex     sync.Mutex
	client    *http.Client
}

type Request struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Auth    string      `json:"auth,omitempty"`
	ID      int         `json:"id"`
}

type LoginResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error"`
	ID int `json:"id"`
}

type Trigger struct {
	TriggerID   string `json:"triggerid"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	LastChange  string `json:"lastchange"`
	Hosts       []Host `json:"hosts"`
	Value       string `json:"value"`
	LastEvent   Event  `json:"lastEvent"`
}

type Host struct {
	HostID string `json:"hostid"`
	Name   string `json:"name"`
}

type Event struct {
	EventID      string   `json:"eventid"`
	Acknowledged string   `json:"acknowledged"`
	Clock        UnixTime `json:"clock"`    // 事件时间戳（秒）
	Value        string   `json:"value"`    // 事件状态（1=问题，0=恢复）
	Name         string   `json:"name"`     // 事件名称
	Severity     string   `json:"severity"` // 严重等级（与触发器priority可能不同）
}

type TriggerResponse struct {
	Jsonrpc string    `json:"jsonrpc"`
	Result  []Trigger `json:"result"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error"`
	ID int `json:"id"`
}

// TriggerState 触发器状态
type TriggerState struct {
	Acknowledged bool
}

// 全局变量
var (
	zabbixClient = &ZabbixClient{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 跳过证书验证
			},
			Timeout: 15 * time.Second,
		},
	}
	shutdownSignal   = make(chan struct{})
	previousTriggers = make(map[string]TriggerState)
	previousMutex    sync.Mutex
	once             sync.Once
)

// UnixTime 自定义时间类型
type UnixTime struct {
	time.Time
}

// AlertMessage 告警消息结构体
type AlertMessage struct {
	EventID      string   `json:"event_id"`       // 事件ID
	TriggerID    string   `json:"trigger_id"`     // 触发器ID
	Description  string   `json:"description"`    // 触发器描述
	Hosts        []string `json:"affected_hosts"` // 主机
	Status       string   `json:"status"`         // 状态：PROBLEM/OK
	Severity     string   `json:"severity"`       // 严重等级
	Timestamp    UnixTime `json:"timestamp"`      // 事件时间戳
	Acknowledged string   `json:"acknowledged"`   // 是否已确认
	//TriggerURL   string   `json:"trigger_url"`    // 可选的链接
}
