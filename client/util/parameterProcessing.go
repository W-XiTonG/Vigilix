package util

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func Parameter() (error, string, int, string, string) {
	// 定义命令行参数
	ipPtr := flag.String("h", "127.0.0.1", "服务器IP地址")
	portPtr := flag.Int("P", 8081, "服务器端口")
	userPtr := flag.String("u", "", "用户名")
	passwordPtr := flag.String("p", "", "密码")
	// 解析命令行参数
	flag.Parse()
	// 检查必须参数是否传递
	provided := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		provided[f.Name] = true
	})

	var missing []string
	if !provided["h"] {
		missing = append(missing, "服务器IP地址(-h)")
	}
	if !provided["P"] {
		missing = append(missing, "服务器端口(-P)")
	}
	if !provided["u"] {
		missing = append(missing, "用户名(-u)")
	}
	if !provided["p"] {
		missing = append(missing, "密码(-p)")
	}

	if len(missing) > 0 {
		_, err := fmt.Fprintln(os.Stderr, "错误：缺少以下必需参数：")
		if err != nil {
			log.Fatal(err)
			return err, "", 0, "", ""
		}
		for _, arg := range missing {
			_, err = fmt.Fprintln(os.Stderr, "  *", arg)
			if err != nil {
				log.Fatal(err)
				return err, "", 0, "", ""
			}
		}
		_, err = fmt.Fprintln(os.Stderr, "\n参数说明：")
		if err != nil {
			log.Fatal(err)
			return err, "", 0, "", ""
		}
		flag.PrintDefaults()
		os.Exit(1)
		return nil, "", 0, "", ""
	}

	return nil, *ipPtr, *portPtr, *userPtr, *passwordPtr
}
