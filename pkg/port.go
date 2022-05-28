package pkg

import (
	"fmt"
	"net"
	"time"
)

// 输入一个端口，返回一个可用端口
func GetAvailablePort(port int) int {
	// 检测端口
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf(":%d", port), time.Millisecond*100)
		if err != nil {
			return port
		}
		conn.Close()
		port++
	}
}
