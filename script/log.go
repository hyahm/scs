package script

import (
	"bufio"
	"io"
	"scs/global"
	"time"

	"github.com/hyahm/golog"
)

func (s *Script) appendLog(line string) {
	t := time.Now().Format("2006/1/2 15:04:05")
	line = t + " -- " + line
	if len(s.Log) >= global.LogCount {
		copy(s.Log, s.Log[1:])
		s.Log[global.LogCount-1] = line
	} else {
		s.Log = append(s.Log, line)
	}
}

func (s *Script) read() {
	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		golog.Error(err)
	}

	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		golog.Error(err)
	}

	//实时循环读取输出流中的一行内容
	go s.appendRead(stderr, true)

	//实时循环读取输出流中的一行内容
	go s.appendRead(stdout, false)
}

func (s *Script) appendRead(stdout io.ReadCloser, iserr bool) {
	readout := bufio.NewReader(stdout)
	for {
		line, err := readout.ReadString('\n')
		if err != nil || io.EOF == err {
			stdout.Close()
			break
		}

		if cap(s.Log) == 0 {
			if iserr {
				golog.Error(line)
			} else {
				golog.Info(line)
			}
		} else {
			s.appendLog(line)
		}

	}
}
