package collectors

// SystemMetrics 系统指标结构体
type SystemMetrics struct {
	//Timestamp  string       `json:"timestamp"`
	Host       HostInfo     `json:"host"`
	CPU        CPUStats     `json:"cpu"`
	Memory     MemoryStats  `json:"memory"`
	Disk       DiskStats    `json:"disk"`
	Network    NetworkStats `json:"network"`
	SystemLoad LoadStats    `json:"load"`
	Processes  []Process    `json:"processes,omitempty"`
}

type HostInfo struct {
	Hostname     string `json:"hostname"`
	OS           string `json:"os"`
	Kernel       string `json:"kernel"`
	Architecture string `json:"arch"`
}

type CPUStats struct {
	PhysicalCores int     `json:"physical_cores"`
	LogicalCores  int     `json:"logical_cores"`
	TotalUsage    float64 `json:"total_usage"`
	UserMode      float64 `json:"user_mode"`
	SystemMode    float64 `json:"system_mode"`
	Idle          float64 `json:"idle"`
}

type MemoryStats struct {
	Total       uint64  `json:"total_bytes"`
	Used        uint64  `json:"used_bytes"`
	UsedPercent float64 `json:"used_percent"`
}

type DiskStats struct {
	Disk       DiskPartition   `json:"disk"`
	Partitions []DiskPartition `json:"partitions"`
}

type DiskPartition struct {
	MountPoint  string  `json:"mount"`
	Total       uint64  `json:"total_bytes"`
	Used        uint64  `json:"used_bytes"`
	UsedPercent float64 `json:"used_percent"`
}

type NetworkStats struct {
	SentTotal  uint64         `json:"sent_total"`
	RecvTotal  uint64         `json:"recv_total"`
	Interfaces []NetInterface `json:"interfaces"`
}

type NetInterface struct {
	Name    string `json:"name"`
	Sent    uint64 `json:"sent_bytes"`
	Recv    uint64 `json:"recv_bytes"`
	Packets uint64 `json:"packets"`
}

type LoadStats struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

type Process struct {
	PID        int32   `json:"pid"`
	Name       string  `json:"name"`
	CPUPercent float64 `json:"cpu_percent"`
	MemPercent float32 `json:"mem_percent"`
}
