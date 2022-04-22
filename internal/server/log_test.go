package server

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/hyahm/golog"
	"golang.org/x/net/context"
)

func TestByte(t *testing.T) {
	a := []byte{32}
	a = bytes.Trim(a, " ")
	t.Logf("---%v---", string(a))
}

func TestLog(t *testing.T) {
	defer golog.Sync()
	cmd := exec.Command("powershell", "-c", "ls")
	svc := Server{
		Cmd: cmd,
	}
	svc.Ctx, svc.Cancel = context.WithCancel(context.Background())
	// out, _ := cmd.StdoutPipe()
	stdout, err := svc.Cmd.StdoutPipe()
	if err != nil {
		golog.Error(err)
	}

	stderr, err := svc.Cmd.StderrPipe()
	if err != nil {
		golog.Error(err)
	}
	// svc.Msg = make(chan string, 1000)
	// 数据同步到log命令
	// go svc.appendLog()
	//实时循环读取输出流中的一行内容
	go svc.appendRead(stderr, true)
	//实时循环读取输出流中的一行内容
	go svc.appendRead(stdout, false)
	// go func() {
	// 	r := bufio.NewReader(out)
	// 	for {
	// 		line, _, err := r.ReadLine()
	// 		if err != nil {
	// 			t.Log(err)
	// 			return
	// 		}
	// 		t.Log(string(line))
	// 	}
	// }()
	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		t.Fatal(err)
	}
	svc.Cancel()
	// out, err := cmd.Output()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(string(out))
}
