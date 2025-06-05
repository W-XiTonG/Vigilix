package SMS

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os/exec"
)

func CmdSms(command []string) (string, error) {
	log.Println("正在执行命令：", command)
	if len(command) == 0 {
		return "", fmt.Errorf("命令为空")
	}
	cmd := exec.Command(command[0], command[1:]...)
	output, err := cmd.Output()
	log.Println(cmd)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return string(exitErr.Stderr), err
		}
		return string(output), fmt.Errorf("执行命令出错: %v，输出信息: %s", err, string(output))
	}
	return string(output), nil
}

func SendSms(ContentSType int, Type int8, Command, UserNumber []string, SmgIpPort, LoginName, LoginPassword,
	Message, File, Content, SPNumber, ChargeNumber, CorpID, ServiceType string) {
	log.Print("正在发送短信：")
	// Type 0：使用自身参数发送，1：调用文件发送
	if Type == 0 {
		log.Println("类型0：调用Sgip板块发送……")
		conn, err := net.Dial("tcp", SmgIpPort)
		if err != nil {
			log.Printf("连接失败: %v\n", err)
			return
		}
		SGIPBindSend(LoginName, LoginPassword, Message, File, Content, SPNumber, ChargeNumber, CorpID, ServiceType, UserNumber, ContentSType, conn)

		if err = conn.Close(); err != nil {
			log.Printf("conn.Close: %v\n", err)
		}
	} else if Type == 1 {
		log.Println("类型1：调用命令发送（", Command, "）")
		command := Command
		out, err := CmdSms(command)
		if err != nil {
			log.Println("命令执行错误：", err)
			return
		}
		log.Println("命令执行成功", out)
	}
}
