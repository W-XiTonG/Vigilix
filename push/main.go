package main

import (
	"Push/config"
	"Push/src/SMS"
	"Push/src/dingding"
	"Push/src/enterpriseWeChat"
	Http "Push/src/http"
	"Push/src/mail"
	"Push/util"
	"log"
)

const Version = "ViGiLix-Push_V0.1.2"

func main() {
	Config := &config.DefaultMailConfigProvider{}
	YamlConfig := Config.GetMailConfig()
	// 创建 LogGer 实例
	LogGer := util.LogGer{}
	//log.Println(YamlConfig.LogGer.LogFile)
	LogGer.Init(YamlConfig.LogGer.LogFile, YamlConfig.LogGer.FileStatus, YamlConfig.LogGer.OutStatus, YamlConfig.LogGer.Status)
	//log.Println(YamlConfig.BeforeAndAfterCommands.BeforeStatus)
	// 版本号
	log.Printf("Version: %s\n", Version)
	if YamlConfig.BeforeAndAfterCommands.BeforeStatus {
		if len(YamlConfig.BeforeAndAfterCommands.BeforeStatusCommand) > 0 {
			out, err := SMS.CmdSms(YamlConfig.BeforeAndAfterCommands.BeforeStatusCommand)
			if err != nil {
				log.Println("命令执行错误：", err)
				return
			}
			log.Println("命令执行成功", out)
		}
	}
	push := func(Content string) {
		if YamlConfig.Status.Mail {
			mail.Mail(
				YamlConfig.ContentS.Type,
				YamlConfig.Mail.Body,
				YamlConfig.ContentS.File,
				Content,
				YamlConfig.Mail.ReceiverEmail,
				YamlConfig.Mail.SenderEmail,
				YamlConfig.Mail.Subject,
				YamlConfig.Mail.SmtpServer,
				YamlConfig.Mail.SenderPassword,
				YamlConfig.Mail.CcEmails,
			)
		}
		if YamlConfig.Status.DingDing {
			dingding.DingDing(
				YamlConfig.DingDing.Secret,
				YamlConfig.DingDing.WebhookURL,
				YamlConfig.DingDing.Message,
				YamlConfig.ContentS.File,
				Content,
				YamlConfig.ContentS.Type,
			)
		}
		if YamlConfig.Status.SMS {
			SMS.SendSms(
				YamlConfig.ContentS.Type,
				YamlConfig.SMS.Type,
				YamlConfig.SMS.Command,
				YamlConfig.SMS.Sgip.UserNumber,
				YamlConfig.SMS.Sgip.SmgIpPort,
				YamlConfig.SMS.Sgip.LoginName,
				YamlConfig.SMS.Sgip.LoginPassword,
				YamlConfig.SMS.Sgip.Message,
				YamlConfig.ContentS.File,
				Content,
				YamlConfig.SMS.Sgip.SPNumber,
				YamlConfig.SMS.Sgip.ChargeNumber,
				YamlConfig.SMS.Sgip.CorpID,
				YamlConfig.SMS.Sgip.ServiceType,
			)
		}
		if YamlConfig.Status.EnterpriseWeChat {
			enterpriseWeChat.PushEnterpriseWeChat(
				YamlConfig.EnterpriseWeChat.WebhookURL,
				YamlConfig.EnterpriseWeChat.Message,
				Content,
				YamlConfig.ContentS.File,
				YamlConfig.EnterpriseWeChat.AteSpecifyStatus,
				YamlConfig.EnterpriseWeChat.MentionedMobileList,
				YamlConfig.ContentS.Type,
			)
		}
	}
	if YamlConfig.Listening.Status {
		// 配置异步参数
		asyncConf := Http.AsyncConfig{
			MaxWorkers:    YamlConfig.Listening.MaxWorkers,
			QueueSize:     YamlConfig.Listening.QueueSize,
			WorkerTimeout: YamlConfig.Listening.WorkerTimeout,
		}
		Http.Server(
			push,
			YamlConfig.Listening.LineBreaks,
			YamlConfig.Listening.Port,
			YamlConfig.Listening.DeleteString,
			YamlConfig.Listening.AuthenticationKeyword,
			asyncConf,
			YamlConfig.Listening.LineBreaksStatus,
			YamlConfig.Listening.DeleteStringStatus,
			YamlConfig.Listening.AuthenticationStatus,
		)
	} else {
		push(YamlConfig.ContentS.Content)
	}
	if YamlConfig.BeforeAndAfterCommands.AfterStatus {
		if len(YamlConfig.BeforeAndAfterCommands.AfterCommand) > 0 {
			out, err := SMS.CmdSms(YamlConfig.BeforeAndAfterCommands.AfterCommand)
			if err != nil {
				log.Println("命令执行错误：", err)
				return
			}
			log.Println("命令执行成功", out)
		}
	}
}
