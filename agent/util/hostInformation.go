package util

import (
	"fmt"
	"net"
)

func GetIPByInterfaceName(interfaceName string) (string, error) {
	ifAce, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return "", err
	}
	address, err := ifAce.Addrs()
	if err != nil {
		return "", err
	}
	for _, addr := range address {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip != nil && !ip.IsLoopback() {
			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("未找到 %s 接口的有效 IP 地址", interfaceName)
}
