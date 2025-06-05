package Http

import (
	"server/handlers"
	"server/util"
)

// Message 定义一个 JSON 数据
type Message struct {
	Type    int    `json:"type"`
	Time    string `json:"time"`
	Content string `json:"content"`
}

type Agent struct {
	Type        int                       `json:"type"`
	ClientID    int                       `json:"client_id"`
	Key         string                    `json:"key"`
	DeliverTime string                    `json:"deliverTime"`
	Data        *util.ClientSystemMetrics `json:"data"`
}

//type Alarm struct {
//	HostName   string  `json:"HostName"`
//	IpAddr     string  `json:"IpAddr"`
//	AlertTime  string  `json:"AlertTime"`
//	Message    string  `json:"Message"`
//	MetricType string  `json:"MetricType"`
//	Current    float64 `json:"Current"`
//	Partition  string  `json:"Partition"`
//	State      string  `json:"State"`
//}

type agentAlarm struct {
	Type        int            `json:"type"`
	ClientID    int            `json:"client_id"`
	Key         string         `json:"key"`
	DeliverTime string         `json:"deliverTime"`
	Data        handlers.Alarm `json:"data"`
}

// 客户端工具连接登录
type clientToolData struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Command  string `json:"Command"`
}
