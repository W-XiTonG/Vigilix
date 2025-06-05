package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"server/config"
	Http "server/interactive"
	"server/util"
)

const Version = "ViGiLix_Server_V0.1.2"

func main() {
	Config := &config.DefaultConfigProvider{}
	YamlConfig := Config.GetMailConfig()
	// 创建 LogGer 实例
	LogGer := util.LogGer{}
	//log.Println(YamlConfig.LogGer.LogFile)
	LogGer.Init(YamlConfig.LogGer.LogFile, YamlConfig.LogGer.FileStatus, YamlConfig.LogGer.OutStatus, YamlConfig.LogGer.Status)
	// 版本号
	log.Printf("Version: %s\n", Version)
	var upGrader = websocket.Upgrader{
		ReadBufferSize:  YamlConfig.WebSocket.ReadBufferSize,
		WriteBufferSize: YamlConfig.WebSocket.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// 判断agents是否从数据库获取数据
	var agents map[int][]string
	if YamlConfig.Agents.DatabaseOrConfigurationFile {
		agents = util.ReadDatabaseAgentData(YamlConfig.Database.DatabaseParameter, YamlConfig.Database.DatabaseIP, YamlConfig.Database.DatabaseUser,
			YamlConfig.Database.DatabasePass, YamlConfig.Database.DatabaseName, YamlConfig.Database.DatabasePort)
	} else {
		agents = YamlConfig.Agents.Id
	}
	// 判断client是否从数据库读取
	var clientTool map[string]string
	if YamlConfig.ClientTool.DatabaseOrConfigurationFile {
		clientTool = util.ReadDatabaseClientData(YamlConfig.Database.DatabaseParameter, YamlConfig.Database.DatabaseIP, YamlConfig.Database.DatabaseUser,
			YamlConfig.Database.DatabasePass, YamlConfig.Database.DatabaseName, YamlConfig.Database.DatabasePort)
	} else {
		clientTool = YamlConfig.ClientTool.ClientConfig
	}
	//agentsId, agentMountPoint, agentAuthenticationKey := util.SplitMapValues(YamlConfig.Agents.Id)
	//agentsIdAndAgentMountPoint := util.NewCustomMap(agentsId, agentMountPoint, agentAuthenticationKey)

	// 打开客户端连接
	Http.Server(
		YamlConfig.WebSocket.ServerPort,
		YamlConfig.Agents.DetectionFrequency,
		YamlConfig.Agents.DetectionTime,
		YamlConfig.ZabbixAlarm.GetIntervalMin,
		YamlConfig.ZabbixAlarm.GetIntervalMax,
		upGrader,
		YamlConfig.Table.Status,
		YamlConfig.Push.PushStatus,
		YamlConfig.Agents.AgentAuthenticationStatus,
		YamlConfig.Agents.DetectionStatus,
		YamlConfig.ZabbixAlarm.Status,
		YamlConfig.ZabbixAlarm.EnableDebug,
		YamlConfig.Table.TablePaths,
		YamlConfig.Push.PushAddr,
		YamlConfig.ZabbixAlarm.ZabbixURL,
		YamlConfig.ZabbixAlarm.ZabbixUser,
		YamlConfig.ZabbixAlarm.ZabbixPass,
		agents,
		YamlConfig.ClientTool.Status,
		clientTool,
	)
}
