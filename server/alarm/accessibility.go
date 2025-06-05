package alarm

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// 生成随机间隔
func generateInterval(maxInterval, minInterval time.Duration) time.Duration {
	// 1. 计算有效区间范围
	intervalRange := maxInterval - minInterval

	// 2. 处理无效区间（当max < min时）
	if intervalRange <= 0 {
		return minInterval
	}

	// 3. 生成随机增量（包含纳秒级精度）
	randomOffset := time.Duration(rand.Int63n(int64(intervalRange)))

	// 4. 组合最终结果
	return minInterval + randomOffset
}

// 带锁的Token获取
func (zc *ZabbixClient) getToken() string {
	zc.mutex.Lock()
	defer zc.mutex.Unlock()
	return zc.authToken
}

// 登录并更新Token
func loginAndUpdateToken(user, pass, zabbixURL string, enableDebug bool) (string, error) {
	zabbixClient.mutex.Lock()
	defer zabbixClient.mutex.Unlock()

	token, err := login(user, pass, zabbixURL, enableDebug)
	if err != nil {
		return "", fmt.Errorf("[Zabbix]登录失败: %w", err)
	} else {
		log.Println("[Zabbix]登录成功")
	}

	zabbixClient.authToken = token
	log.Println("[Zabbix]成功更新认证Token")
	return token, nil
}
