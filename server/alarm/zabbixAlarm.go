package alarm

import (
	"log"
	"strings"
	"time"
)

// Shutdown 导出Shutdown函数供主程序调用
func Shutdown() {
	once.Do(func() {
		close(shutdownSignal)
	})
}

func ZabbixAlarm(pushUrl, zabbixURL, zabbixUser, zabbixPass string, IntervalMin, IntervalMax time.Duration, enableDebug bool) {
	var builder strings.Builder
	builder.WriteString(zabbixURL)
	builder.WriteString("/api_jsonrpc.php")
	zabbixURL = builder.String()
	// 信号处理
	go signals()

	// 初始化连接
	if err := initializeConnection(zabbixUser, zabbixPass, zabbixURL, enableDebug); err != nil {
		log.Fatalf("[ Server -> Zabbix ] 初始化失败: %v", err)
	} else {
		log.Println("[Zabbix]初始化成功")
	}

	// 启动监控循环
	startMonitoring(IntervalMin, IntervalMax, pushUrl, zabbixUser, zabbixPass, zabbixURL, enableDebug)
}
