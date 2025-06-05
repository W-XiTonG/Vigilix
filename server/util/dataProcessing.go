package util

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// SplitMapValues 函数用于将 map[string]string 中的值拆分为 map[string]int 和 map[string]string
func SplitMapValues(inputMap map[string]string) (map[string]int, map[string]string, map[string]string) {
	intMap := make(map[string]int)
	stringMap := make(map[string]string)
	stringKey := make(map[string]string)

	for key, value := range inputMap {
		// 使用逗号分割值
		parts := strings.SplitN(value, ",", 3)
		if len(parts) == 3 {
			// 尝试将第一部分转换为整数
			if num, err := strconv.Atoi(parts[0]); err == nil {
				intMap[key] = num
			}
			stringMap[key] = parts[1]
			stringKey[key] = parts[2]
		}
	}
	return intMap, stringMap, stringKey
}

// BytesToMB 将字节转换为兆字节
func BytesToMB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024)
}

// ConvertToInt 尝试将 interface{} 类型转换为 int 类型
func ConvertToInt(data interface{}) (int, bool) {
	//var result int
	//var ok bool

	switch v := data.(type) {
	case int:
		return v, true
	case int8:
		return int(v), true
	case int16:
		return int(v), true
	case int32:
		return int(v), true
	case int64:
		// 注意：可能会截断超出 int 范围的值
		return int(v), true
	case uint:
		// 注意：可能会截断超出 int 范围的值
		return int(v), true
	case uint8:
		return int(v), true
	case uint16:
		return int(v), true
	case uint32:
		// 注意：可能会截断超出 int 范围的值
		return int(v), true
	case uint64:
		// 注意：可能会截断超出 int 范围的值
		return int(v), true
	case float32:
		// 注意：会丢失小数部分
		return int(v), true
	case float64:
		// 注意：会丢失小数部分
		return int(v), true
	case string:
		result, err := strconv.Atoi(v)
		return result, err == nil
	default:
		return 0, false
	}
}

// ConvertToString 尝试将 interface{} 类型转换为 string 类型
func ConvertToString(data interface{}) (string, bool) {
	switch v := data.(type) {
	case string:
		return v, true
	case int:
		return strconv.Itoa(v), true
	case int8:
		return strconv.FormatInt(int64(v), 10), true
	case int16:
		return strconv.FormatInt(int64(v), 10), true
	case int32:
		return strconv.FormatInt(int64(v), 10), true
	case int64:
		return strconv.FormatInt(v, 10), true
	case uint:
		return strconv.FormatUint(uint64(v), 10), true
	case uint8:
		return strconv.FormatUint(uint64(v), 10), true
	case uint16:
		return strconv.FormatUint(uint64(v), 10), true
	case uint32:
		return strconv.FormatUint(uint64(v), 10), true
	case uint64:
		return strconv.FormatUint(v, 10), true
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), true
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), true
	case []byte:
		return string(v), true
	default:
		return fmt.Sprintf("%v", v), false
	}
}

// ExtractIP 从RemoteAddr中提取IP（处理IPv6带端口的情况）
func ExtractIP(fullAddr string) string {
	host, _, err := net.SplitHostPort(fullAddr)
	if err != nil {
		return fullAddr // 降级处理（无端口情况）
	}
	return host
}
