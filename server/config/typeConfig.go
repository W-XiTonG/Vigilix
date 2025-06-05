package config

import "time"

type Push struct {
	PushStatus bool   `yaml:"PushStatus"`
	PushAddr   string `yaml:"PushAddr"`
}
type LogGer struct {
	Status     bool   `yaml:"Status"`
	OutStatus  bool   `yaml:"OutStatus"`
	FileStatus bool   `yaml:"FileStatus"`
	LogFile    string `yaml:"LogFile"`
}

type Database struct {
	DatabaseIP        string `yaml:"DatabaseIP"`
	DatabasePort      int    `yaml:"DatabasePort"`
	DatabaseName      string `yaml:"DatabaseName"`
	DatabaseUser      string `yaml:"DatabaseUser"`
	DatabasePass      string `yaml:"DatabasePass"`
	DatabaseParameter string `yaml:"DatabaseParameter"`
}

type Agents struct {
	AgentAuthenticationStatus   bool             `yaml:"AgentAuthenticationStatus"`
	DatabaseOrConfigurationFile bool             `yaml:"DatabaseOrConfigurationFile"`
	Id                          map[int][]string `yaml:"Id"`
	DetectionStatus             bool             `yaml:"DetectionStatus"`
	DetectionTime               time.Duration    `yaml:"DetectionTime"`
	DetectionFrequency          int              `yaml:"DetectionFrequency"`
}
type WebSocket struct {
	ServerPort      int `yaml:"ServerPort"`
	ReadBufferSize  int `yaml:"ReadBufferSize"`
	WriteBufferSize int `yaml:"WriteBufferSize"`
}

type Table struct {
	Status     bool   `yaml:"Status"`
	TablePaths string `yaml:"TablePaths"`
}

// ClientTool 客户端连接工具
type ClientTool struct {
	Status                      bool              `yaml:"Status"`
	DatabaseOrConfigurationFile bool              `yaml:"DatabaseOrConfigurationFile"`
	ClientConfig                map[string]string `yaml:"ClientConfig"`
}

// ZabbixAlarm zabbix告警
type ZabbixAlarm struct {
	Status         bool          `yaml:"Status"`
	EnableDebug    bool          `yaml:"EnableDebug"`
	ZabbixURL      string        `yaml:"ZabbixURL"`
	ZabbixUser     string        `yaml:"ZabbixUser"`
	ZabbixPass     string        `yaml:"ZabbixPass"`
	GetIntervalMin time.Duration `yaml:"GetIntervalMin"`
	GetIntervalMax time.Duration `yaml:"GetIntervalMax"`
}

type YamlConfig struct {
	WebSocket   WebSocket   `yaml:"WebSocket"`
	Table       Table       `yaml:"Table"`
	Agents      Agents      `yaml:"Agents"`
	LogGer      LogGer      `yaml:"LogGer"`
	Push        Push        `yaml:"Push"`
	ClientTool  ClientTool  `yaml:"ClientTool"`
	ZabbixAlarm ZabbixAlarm `yaml:"ZabbixAlarm"`
	Database    Database    `yaml:"Database"`
}
