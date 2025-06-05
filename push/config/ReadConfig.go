package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

func Config() (string, string) {
	// 获取可执行文件所在的目录
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("获取可执行文件路径失败: %v\n", err)
		return "", ""
	}
	dir := filepath.Dir(exePath)
	config := "config"
	relativePath := "config.yaml"
	relativePath = filepath.Join(config, relativePath)
	return dir, relativePath

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

// DefaultMailConfigProvider 实现接口的结构体
type DefaultMailConfigProvider struct{}

func (d *DefaultMailConfigProvider) GetMailConfig() PushConfig {
	FilePath := filepath.Join(Config())
	YamlFile, err := os.Open(FilePath)
	if err != nil {
		log.Printf("配置文件打开失败: %v\n", err)
		return PushConfig{}
	}
	defer func() {
		if err = YamlFile.Close(); err != nil {
			log.Printf("YamlFile.Close: %v\n", err)
		}
	}()

	// 解析yaml文件内容到MailConfig
	var YamlContent PushConfig
	err = yaml.NewDecoder(YamlFile).Decode(&YamlContent)
	if err != nil {
		log.Printf("配置文件解析失败: %v\n", err)
		return PushConfig{}
	}
	// 这里可以根据实际情况初始化 MailConfig
	return YamlContent
}
