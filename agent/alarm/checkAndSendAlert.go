package alarm

import (
	"fmt"
	"time"
)

func checkAndSendAlert(alertChan chan Alert, ip string, threshold float64, metrics *SystemMetrics) {
	now := time.Now().Format(time.RFC3339)
	hostName := metrics.HostName

	checkMetric := func(value float64, metricType, partition string) {
		key := fmt.Sprintf("%s-%s-%s", hostName, metricType, partition)

		alertStatus.RLock()
		isAlerting := alertStatus.m[key]
		alertStatus.RUnlock()

		// 触发告警逻辑
		if value > threshold && !isAlerting {
			// 发送告警
			alert := Alert{
				HostName:   hostName,
				IpAddr:     ip,
				AlertTime:  now,
				Message:    fmt.Sprintf("%s 超过阈值（当前: %.2f%%, 阈值: %.2f%%）", metricType, value, threshold),
				MetricType: metricType,
				Current:    value,
				Partition:  partition,
				State:      "产生告警",
			}
			enqueueAlert(alertChan, alert)

			// 更新状态
			alertStatus.Lock()
			alertStatus.m[key] = true
			alertStatus.Unlock()
		} else if value <= threshold && isAlerting {
			// 发送恢复通知
			alert := Alert{
				HostName:   hostName,
				IpAddr:     ip,
				AlertTime:  now,
				Message:    fmt.Sprintf("%s 已恢复正常（当前: %.2f%%, 阈值: %.2f%%）", metricType, value, threshold),
				MetricType: metricType,
				Current:    value,
				Partition:  partition,
				State:      "清除告警",
			}
			enqueueAlert(alertChan, alert)

			// 更新状态
			alertStatus.Lock()
			delete(alertStatus.m, key)
			alertStatus.Unlock()
		}
	}

	// 检查各指标
	checkMetric(metrics.CPUUsage, "CPU", "")
	checkMetric(metrics.MemUsage, "内存", "")
	for partition, usage := range metrics.DiskUsage {
		checkMetric(usage, "磁盘", partition)
	}
	for partition, usage := range metrics.DiskInode {
		checkMetric(usage, "Inode", partition)
	}
}
