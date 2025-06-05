package alarm

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func parseTimestamp(s string) int64 {
	var ts int64
	_, err := fmt.Sscanf(s, "%d", &ts)
	if err != nil {
		log.Println(err)
	}
	return ts
}

func getPriority(p string) string {
	switch p {
	case "0":
		return "未分类"
	case "1":
		return "信息"
	case "2":
		return "警告"
	case "3":
		return "一般严重"
	case "4":
		return "严重"
	case "5":
		return "灾难"
	default:
		return "未知"
	}
}

func getAckStatus(a string) string {
	if a == "1" {
		return "已确认"
	}
	return "未确认"
}

// 状态映射（示例）
func mapStatus(value string) string {
	switch value {
	case "1":
		return "异常"
	case "0":
		return "正常"
	default:
		return "UNKNOWN"
	}
}

func getHostNames(hosts []Host) []string {
	names := make([]string, len(hosts))
	for i, h := range hosts {
		names[i] = h.Name
	}
	return names
}

// 生成指定范围内的随机时间间隔
func generateRandomInterval(min, max time.Duration) time.Duration {
	// 参数校验
	if min <= 0 || max <= 0 || min > max {
		log.Printf("[Zabbix]无效的时间范围 [%v-%v]，使用默认值 15-30s", min, max)
		min = 15 * time.Second
		max = 30 * time.Second
	}

	// 转换为纳秒计算避免精度丢失
	minNs := min.Nanoseconds()
	maxNs := max.Nanoseconds()

	// 生成随机数（包含边界值）
	randomNs := rand.Int63n(maxNs-minNs+1) + minNs
	return time.Duration(randomNs) * time.Nanosecond
}

// 根据TriggerID查找Trigger
func findTriggerByID(triggers []Trigger, id string) (Trigger, bool) {
	for _, t := range triggers {
		if t.TriggerID == id {
			return t, true
		}
	}
	return Trigger{}, false
}

// UnmarshalJSON 时间解析逻辑
func (u *UnixTime) UnmarshalJSON(b []byte) error {
	// 处理可能的引号包裹
	s := strings.Trim(string(b), `"`)

	// 转换为int64
	t, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("[Zabbix]解析时间戳失败: %w，原始值: %s", err, s)
	}

	// Zabbix时间戳为秒级
	u.Time = time.Unix(t, 0)
	return nil
}

// MarshalJSON 添加MarshalJSON方法
func (u UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", u.Unix())), nil
}
