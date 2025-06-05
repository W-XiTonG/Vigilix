package main

import (
	"agent/config"
	"agent/interactive"
	"agent/util"
	"context"
	"log"
)

const Version = "ViGiLix_Agent_V0.1.2"

func main() {
	Config := &config.DefaultConfigProvider{}
	YamlConfig := Config.GetMailConfig()
	// 创建 LogGer 实例
	LogGer := util.LogGer{}
	//log.Println(YamlConfig.LogGer.LogFile)
	LogGer.Init(YamlConfig.LogGer.LogFile, YamlConfig.LogGer.FileStatus, YamlConfig.LogGer.OutStatus, YamlConfig.LogGer.Status)
	// 版本号
	log.Printf("Version: %s\n", Version)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 确保资源被释放
	interactive.Client(
		ctx,
		cancel,
		YamlConfig.AgentId,
		YamlConfig.ServerAddr,
		YamlConfig.AuthenticationKey,
		YamlConfig.Alarm.NetworkCard,
		YamlConfig.ReconnectTime,
		YamlConfig.Alarm.Status,
		YamlConfig.Alarm.CheckInterval,
		YamlConfig.Alarm.Queue,
		YamlConfig.Alarm.Threshold,
	)
}
