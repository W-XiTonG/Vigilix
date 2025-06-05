package alarm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 登录逻辑
func login(user, pass string, zabbixURL string, enableDebug bool) (string, error) {
	loginReq := Request{
		Jsonrpc: "2.0",
		Method:  "user.login",
		Params: map[string]string{
			"user":     user,
			"password": pass,
		},
		ID: 1,
	}

	body, err := sendRequest(zabbixURL, enableDebug, loginReq)
	if err != nil {
		return "", err
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return "", fmt.Errorf("[Zabbix]响应解析失败: %w", err)
	}

	if loginResp.Error.Code != 0 {
		return "", fmt.Errorf("[Zabbix]API错误 %d: %s", loginResp.Error.Code, loginResp.Error.Data)
	}

	return loginResp.Result, nil
}

// 获取触发器
func getTriggers(zabbixURL, authToken string, enableDebug bool) ([]Trigger, error) {
	triggerReq := Request{
		Jsonrpc: "2.0",
		Method:  "trigger.get",
		Params: map[string]interface{}{
			"output":          []string{"triggerid", "description", "priority", "lastchange"},
			"filter":          map[string]interface{}{"value": 1, "status": 0},
			"selectHosts":     []string{"hostid", "name"},
			"selectLastEvent": "extend",
		},
		Auth: authToken,
		ID:   2,
	}

	body, err := sendRequest(zabbixURL, enableDebug, triggerReq)
	if err != nil {
		return nil, err
	}

	var triggerResp TriggerResponse
	if err := json.Unmarshal(body, &triggerResp); err != nil {
		return nil, fmt.Errorf("[Zabbix]触发器解析失败: %w", err)
	}

	if triggerResp.Error.Code != 0 {
		return nil, fmt.Errorf("[Zabbix]API错误 %d: %s", triggerResp.Error.Code, triggerResp.Error.Data)
	}

	return triggerResp.Result, nil
}

// 发送请求
func sendRequest(zabbixURL string, enableDebug bool, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("[Zabbix]JSON编码失败: %w", err)
	}

	if enableDebug {
		log.Printf("[Zabbix]请求体: %s", jsonData)
	}

	req, err := http.NewRequest("POST", zabbixURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("[Zabbix]创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json-rpc")

	resp, err := zabbixClient.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[Zabbix]请求失败: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("%v", err)
		}
	}()

	if enableDebug {
		log.Printf("[Zabbix]响应状态: %s", resp.Status)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[Zabbix]HTTP错误状态码: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[Zabbix]读取响应失败: %w", err)
	}

	if enableDebug {
		log.Printf("[Zabbix]原始响应: %s", body)
	}

	return body, nil
}

// 判断认证错误
func isAuthError(err error) bool {
	return err.Error() == "API错误 -32602: Not authorised" ||
		err.Error() == "API错误 -32603: Not authenticated"
}

// 打印触发器详情
func printTriggers(triggers []Trigger) {
	for _, t := range triggers {
		lastChange := time.Unix(parseTimestamp(t.LastChange), 0)
		fmt.Printf("[Zabbix - 告警ID] %s\n描述:    %s\n严重性:  %s\n触发时间: %s\n关联主机: %v\n确认状态: %s\n-----------------------------------\n",
			t.TriggerID,
			t.Description,
			getPriority(t.Priority),
			lastChange.Format("2006-01-02 15:04:05"),
			getHostNames(t.Hosts),
			getAckStatus(t.LastEvent.Acknowledged),
		)
	}
}

// 处理信号
func handleSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("[Zabbix]接收到系统信号: %v", sig)
	close(shutdownSignal)
}
