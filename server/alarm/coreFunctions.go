package alarm

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// 初始化连接
func initializeConnection(user, pass string, zabbixURL string, enableDebug bool) error {
	if _, err := loginAndUpdateToken(user, pass, zabbixURL, enableDebug); err != nil {
		return fmt.Errorf("[Zabbix]初始登录失败: %w", err)
	}
	return nil
}

// 启动监控
func startMonitoring(IntervalMin, IntervalMax time.Duration, pushUrl, user, pass, zabbixURL string, enableDebug bool) {
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())

	// 生成首个随机间隔
	initialInterval := generateRandomInterval(IntervalMin, IntervalMax)
	ticker := time.NewTicker(initialInterval)
	defer ticker.Stop()

	for {
		select {
		case <-shutdownSignal:
			log.Println("[Zabbix]接收到关闭信号，停止Zabbix监控")
			return
		case <-ticker.C:
			// 执行检查逻辑
			processTriggers(pushUrl, user, pass, zabbixURL, enableDebug)

			// 生成新的随机间隔并重置定时器
			newInterval := generateRandomInterval(IntervalMin, IntervalMax)
			ticker.Reset(newInterval)
			if enableDebug {
				log.Printf("[Zabbix]下次检查将在 %v 后执行", newInterval)
			}
		}
	}
}

// 处理触发器
func processTriggers(pushUrl, user, pass, zabbixURL string, enableDebug bool) {
	startTime := time.Now()
	// 获取触发器
	currentToken := zabbixClient.getToken()
	triggers, err := getTriggers(zabbixURL, currentToken, enableDebug)

	// 错误处理
	if err != nil {
		if isAuthError(err) {
			log.Println("[Zabbix]检测到认证失效，尝试重新登录...")
			if newToken, err := loginAndUpdateToken(user, pass, zabbixURL, enableDebug); err == nil {
				triggers, err = getTriggers(zabbixURL, newToken, enableDebug)
			}
		}

		if err != nil {
			log.Printf("[Zabbix]获取触发器失败: %v", err)
			return
		}
	}

	// 处理当前触发器状态
	currentTriggers := make(map[string]TriggerState)
	for _, t := range triggers {
		ack := t.LastEvent.Acknowledged == "1"
		currentTriggers[t.TriggerID] = TriggerState{Acknowledged: ack}
	}

	previousMutex.Lock()
	defer previousMutex.Unlock()

	// 检测消失的告警
	for id, prevState := range previousTriggers {
		if _, exists := currentTriggers[id]; !exists {
			if prevState.Acknowledged {
				log.Printf("[Zabbix - 清除告警] ID: %s - 已确认告警解除", id)
				// 找到对应的Trigger（需要从triggers中查找）
				if trigger, ok := findTriggerByID(triggers, id); ok {
					buildAlertPayload(trigger, pushUrl)
				}
			} else {
				log.Printf("[Zabbix - 自动清除] ID: %s - 未确认告警自动解除", id)
				// 找到对应的Trigger（需要从triggers中查找）
				if trigger, ok := findTriggerByID(triggers, id); ok {
					buildAlertPayload(trigger, pushUrl)
				}
			}
		}
	}

	// 检测新告警和状态变化（重构逻辑）
	for _, t := range triggers {
		id := t.TriggerID
		currentAck := t.LastEvent.Acknowledged == "1"

		prevState, exists := previousTriggers[id]
		if !exists {
			// 全新告警
			log.Printf("[Zabbix - 新告警] ID: %s - %s (主机: %v)",
				id, t.Description, getHostNames(t.Hosts))
			// 找到对应的Trigger（需要从triggers中查找）
			if trigger, ok := findTriggerByID(triggers, id); ok {
				buildAlertPayload(trigger, pushUrl)
			}
			continue
		}

		// 状态变化检测
		if prevState.Acknowledged != currentAck {
			if currentAck {
				log.Printf("[Zabbix - 确认告警] ID: %s - 已人工确认", id)
				// 找到对应的Trigger（需要从triggers中查找）
				if trigger, ok := findTriggerByID(triggers, id); ok {
					buildAlertPayload(trigger, pushUrl)
				}
			} else {
				log.Printf("[Zabbix - 重新激活] ID: %s - 已重新变为未确认状态", id)
				// 找到对应的Trigger（需要从triggers中查找）
				if trigger, ok := findTriggerByID(triggers, id); ok {
					buildAlertPayload(trigger, pushUrl)
				}
			}
		}
	}

	// 更新历史记录
	newPrevious := make(map[string]TriggerState)
	for id, state := range currentTriggers {
		newPrevious[id] = state
	}
	previousTriggers = newPrevious
	// 打印触发器列表日志
	if enableDebug {
		log.Printf("[Zabbix]当前异常触发器ID列表：%v", previousTriggers)
		log.Println("[Zabbix]异常触发器描述:")
		for i, t := range triggers {
			fmt.Printf("%d. %s - %s\n", i+1, t.TriggerID, t.Description)
		}
	}
	if enableDebug {
		// 处理结果
		if len(triggers) > 0 {
			log.Printf("[Zabbix]发现 %d 个活跃告警 [耗时 %v]", len(triggers), time.Since(startTime))
			//printTriggers(triggers)
		} else {
			log.Printf("[Zabbix]当前无活跃告警 [耗时 %v]", time.Since(startTime))
		}
	}
}
