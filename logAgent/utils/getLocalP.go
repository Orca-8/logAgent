package utils

import (
	"fmt"
	"net"
	"strings"
)

// 获取本机对外ip
func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Printf("dial 8.8.8.8:80 faild, err: %v\n", err)
		return ""
	}
	localIP := conn.LocalAddr().String()
	index := strings.Index(localIP, ":")
	return localIP[:index]
}
