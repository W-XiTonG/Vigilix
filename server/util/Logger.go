package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type LogGer struct {
	LogFile *os.File
}

// Init Init方法用于初始化日志记录器
func (l *LogGer) Init(logFile string, FileStatus, OutStatus, disableLog bool) {

	var err error
	var multiWriter io.Writer
	if !disableLog {
		multiWriter = io.Discard
		log.SetFlags(0)
	} else {
		if FileStatus {
			// 获取当前时间
			now := time.Now()
			// 按照 "YYYY-MM-DD" 格式进行格式化
			formattedDate := now.Format("2006-01-02")
			logFileName := logFile + "server_" + formattedDate + ".log"
			l.LogFile, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Fatalf("打开日志文件(%s)失败: %v", logFile, err)
			}
		}
		switch {
		case FileStatus && !OutStatus:
			multiWriter = io.MultiWriter(l.LogFile)
		case OutStatus && !FileStatus:
			multiWriter = io.MultiWriter(os.Stdout)
		case OutStatus && FileStatus:
			multiWriter = io.MultiWriter(os.Stdout, l.LogFile)
		default:
			fmt.Println(OutStatus, FileStatus)
			multiWriter = io.Discard
		}
	}
	// 设置日志输出至multiWriter
	log.SetOutput(multiWriter)
	// 设置日志标识位，添加日期和时间
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
