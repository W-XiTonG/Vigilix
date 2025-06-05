package collectors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"net/http"
	"runtime"

	"github.com/shirou/gopsutil/v3/process"
)

// Metrics 采集入口
func Metrics(diskPartition string) *SystemMetrics {
	metrics, err := collectMetrics(diskPartition)
	if err != nil {
		fmt.Printf("采集失败: %v\n", err)
		return nil
	}

	// 打印格式化JSON
	//if jsonData, err := json.MarshalIndent(metrics, "", "  "); err == nil {
	//	fmt.Println(string(jsonData))
	//}
	return metrics
	// 发送到服务端
	// sendMetrics("http://your-api.com/metrics", metrics)
}

func collectProcesses(limit int) ([]Process, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var result []Process
	count := 0

	for _, p := range processes {
		if count >= limit {
			break
		}

		name, _ := p.Name()
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()

		// Windows进程名称编码转换
		if runtime.GOOS == "windows" {
			if decoded, err := simplifiedchinese.GBK.NewDecoder().String(name); err == nil {
				name = decoded
			}
		}

		result = append(result, Process{
			PID:        p.Pid,
			Name:       name,
			CPUPercent: cpuPercent,
			MemPercent: memPercent,
		})
		count++
	}
	return result, nil
}

//func getRootPath() string {
//	if runtime.GOOS == "windows" {
//		return "C:"
//	}
//	return "/"
//}

func sendMetrics(url string, metrics *SystemMetrics) error {
	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("JSON编码失败: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Printf("resp.Body.Close error: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("服务端错误: %s", resp.Status)
	}
	return nil
}
