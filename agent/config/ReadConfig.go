package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Config() (string, string) {
	// 获取可执行文件所在的目录
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("获取可执行文件路径失败: %v\n", err)
		return "", ""
	}
	dir := filepath.Dir(exePath)
	// 构建配置文件所在的目录
	configDir := filepath.Join(dir, "config")
	relativePath := "config.yaml"
	return configDir, relativePath

	//relativePath := "config.yaml"
	//_, filename, _, _ := runtime.Caller(0)
	//absPath, err := filepath.Abs(filename)
	//if err != nil {
	//	log.Printf("转换为绝对路径失败: %v\n", err)
	//	return "", ""
	//}
	//// 获取目录路径
	//dir := filepath.Dir(absPath)
	//return dir, relativePath
}

type LogGer struct {
	Status     bool   `yaml:"Status"`
	OutStatus  bool   `yaml:"OutStatus"`
	FileStatus bool   `yaml:"FileStatus"`
	LogFile    string `yaml:"LogFile"`
}

type Alarm struct {
	Status        bool          `yaml:"Status"`
	Threshold     float64       `yaml:"Threshold"`
	CheckInterval time.Duration `yaml:"CheckInterval"`
	NetworkCard   string        `yaml:"NetworkCard"`
	Queue         int32         `yaml:"Queue"`
}

type YamlConfig struct {
	ServerAddr        string        `yaml:"ServerAddr"`
	AgentId           int           `yaml:"AgentId"`
	AuthenticationKey string        `yaml:"AuthenticationKey"`
	ReconnectTime     time.Duration `yaml:"ReconnectTime"`
	LogGer            LogGer        `yaml:"LogGer"`
	Alarm             Alarm         `yaml:"Alarm"`
}

type DefaultConfigProvider struct{}

func (d *DefaultConfigProvider) GetMailConfig() YamlConfig {
	//FilePath := filepath.Join(Config())
	dir, relativePath := Config()
	FilePath := filepath.Join(dir, relativePath)
	YamlFile, err := os.Open(FilePath)
	if err != nil {
		log.Printf("配置文件打开失败: %v\n", err)
		return YamlConfig{}
	}
	defer func() {
		if err = YamlFile.Close(); err != nil {
			log.Printf("YamlFile.Close: %v\n", err)
		}
	}()

	// 解析yaml文件内容到MailConfig
	var YamlContent YamlConfig
	err = yaml.NewDecoder(YamlFile).Decode(&YamlContent)
	if err != nil {
		log.Printf("配置文件解析失败: %v\n", err)
		return YamlConfig{}
	}
	// 这里可以根据实际情况初始化 MailConfig
	return YamlContent
}
