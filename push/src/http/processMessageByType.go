package Http

import (
	"encoding/json"
	"fmt"
	"log"
)

// processMessageByType 根据消息类型处理消息内容
func processMessageByType(data customData) (string, error) {
	switch data.Type {
	case customMessage:
		// 自定义消息直接返回内容
		return string(data.Content), nil
	case agentAlarm:
		// 处理代理告警消息
		var alarmData struct {
			HostName   string  `json:"HostName"`
			IpAddr     string  `json:"IpAddr"`
			AlertTime  string  `json:"AlertTime"`
			Message    string  `json:"Message"`
			MetricType string  `json:"MetricType"`
			Current    float64 `json:"Current"`
			Partition  string  `json:"Partition"`
			State      string  `json:"State"`
		}

		log.Printf("%s", string(data.Content))

		if err := json.Unmarshal(data.Content, &alarmData); err != nil {
			return "", fmt.Errorf("解析agent告警失败: %v", err)
		}
		// 使用格式化函数生成告警消息
		return formatAlertMessage(
			alarmData.HostName,
			alarmData.IpAddr,
			alarmData.AlertTime,
			alarmData.Message,
			alarmData.MetricType,
			alarmData.State,
			alarmData.Partition,
			alarmData.Current,
		), nil

	case zabbixAlarm:
		// AlertMessage zabbix告警消息结构体
		var AlertMessage struct {
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
		if err := json.Unmarshal(data.Content, &AlertMessage); err != nil {
			return "", fmt.Errorf("解析zabbix告警失败: %v", err)
		}
		return formatZabbixStatusMessage(
			AlertMessage.EventID,
			AlertMessage.TriggerID,
			AlertMessage.Description,
			AlertMessage.Status,
			AlertMessage.Severity,
			AlertMessage.Acknowledged,
			AlertMessage.Hosts,
			AlertMessage.Timestamp,
		), nil

	case agentStatus:
		// 处理代理状态消息
		var AgentStatus struct {
			ClientID     int    `json:"AgentID"`
			ClientName   string `json:"Agent描述"`
			ClientAddr   string `json:"Agent地址"`
			ClientStatus string `json:"Agent状态"`
		}
		if err := json.Unmarshal(data.Content, &AgentStatus); err != nil {
			return "", fmt.Errorf("解析agent状态失败: %v", err)
		}
		// 使用格式化函数生成告警消息
		return formatStatusMessage(
			AgentStatus.ClientID,
			AgentStatus.ClientName,
			AgentStatus.ClientAddr,
			AgentStatus.ClientStatus,
		), nil

	default:
		return "", fmt.Errorf("未知消息: %v", data)
	}
}
