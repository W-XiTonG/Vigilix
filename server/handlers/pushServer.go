package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// AgentStatus Pusher 告警推送器（封装推送逻辑）
type AgentStatus struct {
	ClientID     int    `json:"AgentID"`
	ClientName   string `json:"Agent描述"`
	ClientAddr   string `json:"Agent地址"`
	ClientStatus string `json:"Agent状态"`
}

// Alarm Alert 告警信息结构体
type Alarm struct {
	HostName   string  `json:"HostName"`
	IpAddr     string  `json:"IpAddr"`
	AlertTime  string  `json:"AlertTime"`
	Message    string  `json:"Message"`
	MetricType string  `json:"MetricType"`
	Current    float64 `json:"Current"`
	Partition  string  `json:"Partition"`
	State      string  `json:"State"`
}

type Data struct {
	Type    int         `json:"Type"`
	Content interface{} `json:"Content"`
}

const (
	customMessage = 1
	agentAlarm    = 2
	zabbixAlarm   = 3
	agentStatus   = 4
)

// Push 将AgentStatus的信息封装为JSON并推送到指定URL
func Push(jsonData interface{}, url string, pushType int) error {
	log.Printf("[ -> Server] %s", jsonData)
	//jsonData, err := json.Marshal(interactive)
	//if err != nil {
	//	return fmt.Errorf("JSON编码失败: %v", err)
	//}
	// 包装到Content字段
	wrappedData := Data{
		Type:    pushType,
		Content: jsonData}
	pushJsonData, err := json.Marshal(wrappedData)
	if err != nil {
		log.Printf("包装数据失败: %v", err)
		return err
	}
	log.Printf("[Server -> Push] %s", string(pushJsonData))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(pushJsonData))
	if err != nil {
		return fmt.Errorf("推送失败: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error: resp.Body.Close %v", err)
		}
	}()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("服务端返回异常状态码: %d", resp.StatusCode)
	}

	return nil
}

// PushServer 断连告警使用
func PushServer(Type int, RemoteAddr string, clientID int, status, PushAddr string, agentsID map[int][]string, agentData *Alarm) {
	switch Type {
	case 100, 101:
		ClientName := AgentName(agentsID, clientID)
		agent := AgentStatus{
			ClientID:     clientID,
			ClientName:   ClientName,
			ClientAddr:   RemoteAddr,
			ClientStatus: status,
		}
		//jsonData, err := json.Marshal(agent)
		//if err != nil {
		//	log.Printf("JSON编码失败: %v", err)
		//}
		err := Push(agent, PushAddr, agentStatus)
		if err != nil {
			log.Printf("推送Agent状态失败: %v", err)
		} else {
			log.Printf("PushAddr:%s;interactive:%v", PushAddr, agent)
		}

	case 103:
		//ClientName := util.FindKey(agentsID, clientID)
		agent := Alarm{
			HostName:   agentData.HostName,
			IpAddr:     agentData.IpAddr,
			AlertTime:  agentData.AlertTime,
			Message:    agentData.Message,
			MetricType: agentData.MetricType,
			Current:    agentData.Current,
			Partition:  agentData.Partition,
			State:      agentData.State,
		}
		//jsonData, err := json.Marshal(agent)
		//if err != nil {
		//	log.Printf("JSON编码失败: %v", err)
		//}
		err := Push(agent, PushAddr, agentAlarm)
		if err != nil {
			log.Printf("推送Agent状态失败: %v", err)
		} else {
			log.Printf("PushAddr:%s;interactive:%v", PushAddr, agent)
		}
	}
}
