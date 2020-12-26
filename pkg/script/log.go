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
				s.Msg <- line
				// s.appendLog(line)
			}
		}
	}
}

func read(cmd *exec.Cmd) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		golog.Error(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		golog.Error(err)
	}

	//实时循环读取输出流中的一行内容
	go appendRead(stderr, true)

	//实时循环读取输出流中的一行内容
	go appendRead(stdout, false)
}

func appendRead(stdout io.ReadCloser, iserr bool) {
	readout := bufio.NewReader(stdout)
	for {

		line, err := readout.ReadString('\n')
		if err != nil || io.EOF == err {
			stdout.Close()
			break
		}
		if iserr {
			golog.Error(line)
		} else {
			golog.Info(line)
		}

	}
}

func (s *Script) appendLog() {
	t := time.Now().Format("2006/1/2 15:04:05")
	for {
		select {
		case <-s.Ctx.Done():
			close(s.Msg)
			for line := range s.Msg {
				line = t + " -- " + line
				if len(s.Log) >= global.LogCount {
					copy(s.Log, s.Log[1:])
					s.Log[global.LogCount-1] = line
				} else {
					s.Log = append(s.Log, line)
				}
			}
			return
		case line := <-s.Msg:
			line = t + " -- " + line
			if len(s.Log) >= global.LogCount {
				copy(s.Log, s.Log[1:])
				s.Log[global.LogCount-1] = line
			} else {
				s.Log = append(s.Log, line)
			}

		}
	}

}
