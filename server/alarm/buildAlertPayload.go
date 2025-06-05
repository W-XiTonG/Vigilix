package alarm

import (
	"log"
	"server/handlers"
)

// 构建推送消息列表
func buildAlertPayload(trigger Trigger, pushUrl string) {
	//var messages []AlertMessage
	//for _, trigger := range triggers {
	// 提取受影响主机列表
	hostNames := make([]string, 0, len(trigger.Hosts))
	for _, host := range trigger.Hosts {
		hostNames = append(hostNames, host.Name)
	}

	// 转换时间戳（根据Zabbix返回的时间格式调整）
	//eventTime, err := time.Parse("2006-01-02T15:04:05", trigger.LastEvent.Value)
	//if err != nil {
	//	log.Printf("时间解析失败: %v", err)
	//	eventTime = time.Now() // 失败时使用当前时间
	//}
	//eventTime := parseZabbixTimestamp(trigger.LastEvent.Value)

	// 构建消息体
	msg := AlertMessage{
		EventID:      trigger.LastEvent.EventID,
		TriggerID:    trigger.TriggerID,
		Description:  trigger.Description,
		Hosts:        hostNames,
		Status:       mapStatus(trigger.LastEvent.Value), // 需要实现状态映射
		Severity:     getPriority(trigger.Priority),      // 需要实现优先级映射
		Timestamp:    trigger.LastEvent.Clock,
		Acknowledged: getAckStatus(trigger.LastEvent.Acknowledged),
		//TriggerURL:   buildTriggerURL(trigger.TriggerID), // URL构造
	}
	err := handlers.Push(msg, pushUrl, 3)
	if err != nil {
		log.Println(err)
	}
	//messages = append(messages, msg)
	//}
}
