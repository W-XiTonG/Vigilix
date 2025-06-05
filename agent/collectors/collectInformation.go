package collectors

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// 采集信息
func collectMetrics(diskPartition string) (*SystemMetrics, error) {
	metrics := &SystemMetrics{
		//Timestamp: time.Now().Format(time.RFC3339),
	}

	// 采集主机信息
	if hostInfo, err := host.Info(); err == nil {
		metrics.Host = HostInfo{
			Hostname:     hostInfo.Hostname,
			OS:           hostInfo.OS,
			Kernel:       hostInfo.KernelVersion,
			Architecture: hostInfo.KernelArch,
		}
	}

	// 采集CPU信息
	if err := collectCPUMetrics(&metrics.CPU); err != nil {
		return nil, fmt.Errorf("CPU采集失败: %v", err)
	}

	// 采集内存信息
	if memInfo, err := mem.VirtualMemory(); err == nil {
		metrics.Memory = MemoryStats{
			Total:       memInfo.Total,
			Used:        memInfo.Used,
			UsedPercent: memInfo.UsedPercent,
		}
	}

	// 采集磁盘信息
	if err := collectDiskMetrics(&metrics.Disk, diskPartition); err != nil {
		return nil, fmt.Errorf("磁盘采集失败: %v", err)
	}

	// 采集网络信息
	if err := collectNetworkMetrics(&metrics.Network); err != nil {
		return nil, fmt.Errorf("网络采集失败: %v", err)
	}

	// 采集系统负载
	if loadStats, err := getLoadStats(); err == nil {
		metrics.SystemLoad = *loadStats
	}

	// 采集进程信息（可选）
	//if processes, err := collectProcesses(5); err == nil {
	//	metrics.Processes = processes
	//}

	return metrics, nil
}
