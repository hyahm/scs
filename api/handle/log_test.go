package handle

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestTailf(t *testing.T) {
	f, err := os.Open("../../log/test_0.log")
	if err != nil {
		t.Skip("log file not found, skip tailf test")
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	// 只读若干行后退出，避免无限循环卡死 go test
	for i := 0; i < 10; i++ {
		line, _, err := buf.ReadLine()
		if err != nil {
			break
		}
		t.Log(string(line))
	}
}

func TestReplace(t *testing.T) {
	wsdomain := "https://aaahttp.bbb.com"
	wsdomain = strings.Replace(wsdomain, "http", "ws", 1)
	t.Log(wsdomain)
}
