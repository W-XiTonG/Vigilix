package alarm

// SystemMetrics 性能指标结构体
type SystemMetrics struct {
	HostName   string             `json:"host_name"`
	CPUUsage   float64            `json:"cpu_usage"`
	MemUsage   float64            `json:"mem_usage"`
	DiskUsage  map[string]float64 `json:"disk_usage"`
	ReportTime string             `json:"report_time"`
	DiskInode  map[string]float64 `json:"disk_inode"`
}

// Alert 告警信息结构体
type Alert struct {
	HostName   string  `json:"HostName"`
	IpAddr     string  `json:"IpAddr"`
	AlertTime  string  `json:"AlertTime"`
	Message    string  `json:"Message"`
	MetricType string  `json:"MetricType"`
	Current    float64 `json:"Current"`
	Partition  string  `json:"Partition"`
	State      string  `json:"State"`
}

// pushClientInfo 定义客户端消息结构体
type pushClientInfo struct {
	ClientID int    `json:"client_id"`
	Type     int    `json:"type"`
	Key      string `json:"key"`
	Time     string `json:"deliverTime"`
	Data     Alert  `json:"data"`
}
