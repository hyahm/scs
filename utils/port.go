package utils

import (
	"fmt"
	"net"
	"time"
)

func ProbePort(port int) int {
	// 检测端口
	index := 0
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf(":%d", port), time.Nanosecond*100)
		if err != nil {
			return index
		}
		if conn != nil {
			_ = conn.Close()
			return index
		}
		index++
	}

}
