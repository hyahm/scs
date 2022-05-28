package server

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"

	"github.com/hyahm/golog"
)

func (svc *Server) read() {
	stdout, err := svc.Cmd.StdoutPipe()
	if err != nil {
		golog.Error(err)
	}
	golog.Info(stdout)
	stderr, err := svc.Cmd.StderrPipe()
	if err != nil {
		golog.Error(err)
	}
	golog.Info(stderr)
	// svc.Msg = make(chan string, 1000)
	// 数据同步到log命令
	// go svc.appendLog()
	//实时循环读取输出流中的一行内容
	go svc.appendRead(stderr, true)
	//实时循环读取输出流中的一行内容
	go svc.appendRead(stdout, false)
}

func (svc *Server) appendRead(stdout io.ReadCloser, iserr bool) {
	golog.Info(stdout)
	readout := bufio.NewReader(stdout)
	for {
		select {
		case <-svc.Ctx.Done():
			// close(svc.Msg)
			return
		default:
			golog.Info(readout)
			line, _, err := readout.ReadLine()
			if err != nil {
				stdout.Close()
				return
			}
			line = bytes.Trim(line, " ")
			if string(line) == "" {
				continue
			}
			if iserr {
				svc.Logger.Error(string(line))
			} else {
				svc.Logger.Info(string(line))
			}
			// t := time.Now().Format("2006/1/2 15:04:05")
			// msg := t + " -- " + string(line)
			// svc.Msg <- msg
		}
	}
}

func read(cmd *exec.Cmd, svc *Server) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		golog.Error(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		golog.Error(err)
	}

	//实时循环读取输出流中的一行内容
	go appendErrRead(stderr, svc)

	//实时循环读取输出流中的一行内容
	go appendRead(stdout, svc)
}

func appendErrRead(stdout io.ReadCloser, svc *Server) {
	readout := bufio.NewReader(stdout)
	for {
		line, err := readout.ReadString('\n')
		if err != nil {
			stdout.Close()
			break
		}
		svc.Logger.Error(line[:len(line)-1])
	}
}

func appendRead(stdout io.ReadCloser, svc *Server) {
	readout := bufio.NewReader(stdout)
	for {
		line, err := readout.ReadString('\n')
		if err != nil {
			stdout.Close()
			break
		}
		svc.Logger.Info(line[:len(line)-1])
	}
}
