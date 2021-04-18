package scs

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/hyahm/scs/global"

	"github.com/hyahm/golog"
)

func (svc *Server) read() {
	stdout, err := svc.cmd.StdoutPipe()
	if err != nil {
		golog.Error(err)
	}

	stderr, err := svc.cmd.StderrPipe()
	if err != nil {
		golog.Error(err)
	}
	svc.Msg = make(chan string, 1000)
	go svc.appendLog()
	//实时循环读取输出流中的一行内容
	go svc.appendRead(stderr, true)
	//实时循环读取输出流中的一行内容
	go svc.appendRead(stdout, false)
}

func (svc *Server) appendRead(stdout io.ReadCloser, iserr bool) {
	readout := bufio.NewReader(stdout)
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		golog.Error(err)
	// 	}
	// }()
	for {
		select {
		case <-svc.Ctx.Done():
			golog.Info("stop")
			close(svc.Msg)
			return
		default:
			line, err := readout.ReadString('\n')
			if err != nil || io.EOF == err {
				stdout.Close()
				return
			}
			if strings.Trim(line[:len(line)-1], " ") == "" {
				return
			}

			svc.LogLocker.RLock()
			svc.LogLocker.RUnlock()
			if iserr {
				golog.Error(line[:len(line)-1])
			} else {
				golog.Info(line[:len(line)-1])
			}
			t := time.Now().Format("2006/1/2 15:04:05")
			line = t + " -- " + line
			svc.Msg <- line
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
		if err != nil || io.EOF == err {
			stdout.Close()
			break
		}
		golog.Info(line[:len(line)-1])
		svc.LogLocker.Lock()
		if len(svc.Log) >= global.LogCount {
			copy(svc.Log, svc.Log[1:])
			svc.Log[global.LogCount-1] = line
		} else {
			svc.Log = append(svc.Log, line)
		}
		svc.LogLocker.Unlock()
	}
}

func appendRead(stdout io.ReadCloser, svc *Server) {
	readout := bufio.NewReader(stdout)
	for {
		line, err := readout.ReadString('\n')
		if err != nil || io.EOF == err {
			stdout.Close()
			break
		}
		golog.Info(line[:len(line)-1])
	}
}

func (svc *Server) appendLog() {
	for {
		select {
		case <-svc.Ctx.Done():
			for line := range svc.Msg {
				svc.LogLocker.Lock()
				if len(svc.Log) >= global.LogCount {
					copy(svc.Log, svc.Log[1:])
					svc.Log[global.LogCount-1] = line
				} else {
					svc.Log = append(svc.Log, line)
				}
				svc.LogLocker.Unlock()
			}
			return
		case line := <-svc.Msg:
			svc.LogLocker.Lock()
			if len(svc.Log) >= global.LogCount {
				copy(svc.Log, svc.Log[1:])
				svc.Log[global.LogCount-1] = line
			} else {
				svc.Log = append(svc.Log, line)
			}
			svc.LogLocker.Unlock()
		}
	}

}
