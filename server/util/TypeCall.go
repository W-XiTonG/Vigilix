package util

// CustomMap 定义一个结构体，包含两个映射
type CustomMap struct {
	IntMap                 map[string]int
	StringMap              map[string]string
	agentAuthenticationKey map[string]string
}

// CombinedMap 定义一个结构体，包含两个映射
type CombinedMap struct {
	intMap    map[string]int
	stringMap map[string]string
}

// NewCustomMap 创建一个新的 CustomMap 实例
func NewCustomMap(agentsId map[string]int, agentMountPoint map[string]string, agentAuthenticationKey map[string]string) *CustomMap {
	// 初始化自定义映射
	customMap := &CustomMap{
		IntMap:                 agentsId,
		StringMap:              agentMountPoint,
		agentAuthenticationKey: agentAuthenticationKey,
	}
	return customMap
}

// IsNumber 函数检查是否为数值类型
func IsNumber(v interface{}) bool {
	switch v.(type) {
	case int, int32, int64, float32, float64:
		return true
	default:
		return false
	}
}
