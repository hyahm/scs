package server

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"

	"github.com/hyahm/golog"
)

func (svc *Server) read() {
	if svc.Cmd == nil {
		return
	}
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
}

func (svc *Server) appendRead(stdout io.ReadCloser, iserr bool) {
	if stdout == nil {
		return
	}
	readout := bufio.NewReader(stdout)
	for {
		select {
		case <-svc.Ctx.Done():
			return
		default:
			if readout == nil {
				return
			}
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

	go appendErrRead(stderr, svc)

	go appendRead(stdout, svc)
}

func appendErrRead(stdout io.ReadCloser, svc *Server) {
	if stdout == nil {
		return
	}
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
	if stdout == nil {
		return
	}
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
