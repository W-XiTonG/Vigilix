package collectors

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/net"
	"runtime"
	"time"
)

func collectCPUMetrics(c *CPUStats) error {
	// 物理/逻辑核心数
	physical, _ := cpu.Counts(false)
	logical, _ := cpu.Counts(true)
	c.PhysicalCores = physical
	c.LogicalCores = logical

	// 总使用率
	percent, err := cpu.Percent(1*time.Second, false)
	if err != nil || len(percent) == 0 {
		return fmt.Errorf("获取CPU使用率失败: %v", err)
	}
	c.TotalUsage = percent[0]

	// 时间分布
	times, err := cpu.Times(false)
	if err != nil || len(times) == 0 {
		return fmt.Errorf("获取CPU时间失败: %v", err)
	}

	total := times[0].Total()
	if total == 0 {
		total = 1 // 防止除零
	}

	c.UserMode = (times[0].User / total) * 100
	c.SystemMode = (times[0].System / total) * 100
	c.Idle = (times[0].Idle / total) * 100

	// Windows特殊处理
	if runtime.GOOS == "windows" {
		c.Idle = 100 - c.TotalUsage
	}

	return nil
}

func collectDiskMetrics(d *DiskStats, diskPartition string) error {
	if usage, err := disk.Usage(diskPartition); err == nil {
		d.Disk = DiskPartition{
			MountPoint:  diskPartition,
			Total:       usage.Total,
			Used:        usage.Used,
			UsedPercent: usage.UsedPercent,
		}
	}

	// 所有分区
	if partitions, err := disk.Partitions(true); err == nil {
		for _, part := range partitions {
			if usage, err := disk.Usage(part.Mountpoint); err == nil {
				d.Partitions = append(d.Partitions, DiskPartition{
					MountPoint:  part.Mountpoint,
					Total:       usage.Total,
					Used:        usage.Used,
					UsedPercent: usage.UsedPercent,
				})
			}
		}
	}
	return nil
}

func collectNetworkMetrics(n *NetworkStats) error {
	if netIO, err := net.IOCounters(true); err == nil {
		var sentTotal, recvTotal uint64
		var interfaces []NetInterface

		for _, ifs := range netIO {
			sentTotal += ifs.BytesSent
			recvTotal += ifs.BytesRecv
			interfaces = append(interfaces, NetInterface{
				Name:    ifs.Name,
				Sent:    ifs.BytesSent,
				Recv:    ifs.BytesRecv,
				Packets: ifs.PacketsSent + ifs.PacketsRecv,
			})
		}

		n.SentTotal = sentTotal
		n.RecvTotal = recvTotal
		n.Interfaces = interfaces
	}
	return nil
}

func getLoadStats() (*LoadStats, error) {
	if runtime.GOOS == "windows" {
		// Windows使用CPU使用率模拟负载
		if percent, err := cpu.Percent(1*time.Second, false); err == nil {
			return &LoadStats{
				Load1:  percent[0],
				Load5:  percent[0],
				Load15: percent[0],
			}, nil
		}
		return nil, fmt.Errorf("无法获取Windows负载")
	}

	// Linux/Mac获取真实负载
	if avg, err := load.Avg(); err == nil {
		return &LoadStats{
			Load1:  avg.Load1,
			Load5:  avg.Load5,
			Load15: avg.Load15,
		}, nil
	}
	return nil, fmt.Errorf("无法获取系统负载")
}
