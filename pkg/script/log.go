package script

import (
	"bufio"
	"io"
	"os/exec"
	"scs/global"
	"time"

	"github.com/hyahm/golog"
)

func (s *Script) read() {
	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		golog.Error(err)
	}

	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		golog.Error(err)
	}
	s.Msg = make(chan string, 1000)
	go s.appendLog()
	//实时循环读取输出流中的一行内容
	go s.appendRead(stderr, true)
	//实时循环读取输出流中的一行内容
	go s.appendRead(stdout, false)
}

func (s *Script) appendRead(stdout io.ReadCloser, iserr bool) {
	readout := bufio.NewReader(stdout)

	for {
		select {
		case <-s.Ctx.Done():
			golog.Info("stop")
			close(s.Msg)
			return
		default:
			line, err := readout.ReadString('\n')
			if err != nil || io.EOF == err {
				stdout.Close()
				return
			}
			if cap(s.Log) == 0 {
				if iserr {
					golog.Error(line)
				} else {
					golog.Info(line)
				}
			} else {
				t := time.Now().Format("2006/1/2 15:04:05")
				golog.Info(line)
				line = t + " -- " + line
				s.Msg <- line
			}
		}
	}
}

func read(cmd *exec.Cmd, s *Script) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		golog.Error(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		golog.Error(err)
	}

	//实时循环读取输出流中的一行内容
	go appendRead(stderr, s)

	//实时循环读取输出流中的一行内容
	go appendRead(stdout, s)
}

func appendRead(stdout io.ReadCloser, s *Script) {
	readout := bufio.NewReader(stdout)
	for {
		line, err := readout.ReadString('\n')
		if err != nil || io.EOF == err {
			stdout.Close()
			break
		}
		golog.Info(line)
	}
}

func (s *Script) appendLog() {
	for {

		select {
		case <-s.Ctx.Done():
			for line := range s.Msg {
				if len(s.Log) >= global.LogCount {
					copy(s.Log, s.Log[1:])
					s.Log[global.LogCount-1] = line
				} else {
					s.Log = append(s.Log, line)
				}
			}
			return
		case line := <-s.Msg:
			if len(s.Log) >= global.LogCount {
				copy(s.Log, s.Log[1:])
				s.Log[global.LogCount-1] = line
			} else {
				s.Log = append(s.Log, line)
			}

		}
	}

}
