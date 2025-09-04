package handle

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestTailf(t *testing.T) {
	f, err := os.Open("../../log/test_0.log")
	if err != nil {
		log.Fatal(err)
	}
	buf := bufio.NewReader(f)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		t.Log(string(line))
	}
}

func TestReplace(t *testing.T) {
	wsdomain := "https://aaahttp.bbb.com"
	wsdomain = strings.Replace(wsdomain, "http", "ws", 1)
	t.Log(wsdomain)
}
