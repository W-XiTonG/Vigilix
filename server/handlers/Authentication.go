package handlers

import (
	"log"
)

func Authentication(handlers int, agentId map[string]int) (string, bool) {
	for key, value := range agentId {
		if value == handlers {
			return key, true
		}
	}
	return "", false
}

//func ValidValues(inputKey string, keyValue map[string]string) (string, bool) {
//	for key, value := range keyValue {
//		if key == inputKey {
//			return value, true
//		}
//	}
//	return "", false
//}

func ValidValues(inputKey string, keyValue map[string]string) (string, bool) {
	value, ok := keyValue[inputKey]
	if ok {
		return value, true
	}
	return "", false
}

// AgentName 根据ID查找描述
func AgentName(config map[int][]string, targetID int) string {
	// 获取 key 对应的切片长度
	if values, exists := config[targetID]; exists {
		if len(values) > 0 {
			return values[0]
		}
	} else {
		log.Printf("[Agent] 键不存在 %d", targetID)
	}
	return ""
}

// AgentKey 根据ID查找鉴权密钥
func AgentKey(config map[int][]string, targetID int) (string, bool) {
	// 获取 key 对应的切片长度
	if values, exists := config[targetID]; exists {
		if len(values) > 2 {
			return values[2], true
		}
	} else {
		log.Printf("[Agent] 键不存在 %d", targetID)
		return "", false
	}
	return "", false
}

// AgentPartition 根据ID查找目录
func AgentPartition(config map[int][]string, targetID int) (string, bool) {
	// 获取 key 对应的切片长度
	if values, exists := config[targetID]; exists {
		if len(values) > 1 {
			return values[1], true
		}
	} else {
		log.Printf("[Agent] 键不存在 %d", targetID)
		return "", false
	}
	return "", false
}

// ClientPartition 查找Client数据
func ClientPartition(config map[string]string, clientUser string) (string, string) {
	value, exists := config[clientUser]
	if exists {
		log.Println("[Client] 匹配数据：", value)
		return clientUser, value
	} else {
		log.Printf("[Client] 键不存在 %s", clientUser)
		return "", ""
	}
}
