package script

import (
	"bufio"
	"io"
	"os/exec"
	"time"

	"github.com/hyahm/scs/global"

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
	defer func() {
		if err := recover(); err != nil {
			golog.Error(err)
		}
	}()
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
			s.LogLocker.RLock()
			logCap := cap(s.Log["log"])
			s.LogLocker.RUnlock()
			if logCap == 0 {
				if iserr {
					golog.Error(line)
				} else {
					golog.Info(line)
				}
			} else {
				t := time.Now().Format("2006/1/2 15:04:05")
				line = t + " -- " + line
				golog.Info(line)
				s.Msg <- line
			}
		}
	}
}

func read(cmd *exec.Cmd, s *Script, typ string) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		golog.Error(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		golog.Error(err)
	}

	//实时循环读取输出流中的一行内容
	go appendRead(stderr, s, typ)

	//实时循环读取输出流中的一行内容
	go appendRead(stdout, s, typ)
}

func appendRead(stdout io.ReadCloser, s *Script, typ string) {
	readout := bufio.NewReader(stdout)
	for {
		line, err := readout.ReadString('\n')
		if err != nil || io.EOF == err {
			stdout.Close()
			break
		}
		golog.Info(line)
		s.LogLocker.Lock()
		if len(s.Log[typ]) >= global.LogCount {
			copy(s.Log[typ], s.Log[typ][1:])
			s.Log[typ][global.LogCount-1] = line
		} else {
			s.Log[typ] = append(s.Log[typ], line)
		}
		s.LogLocker.Unlock()
	}
}

func (s *Script) appendLog() {
	for {
		select {
		case <-s.Ctx.Done():
			for line := range s.Msg {
				s.LogLocker.Lock()
				if len(s.Log["log"]) >= global.LogCount {
					copy(s.Log["log"], s.Log["log"][1:])
					s.Log["log"][global.LogCount-1] = line
				} else {
					s.Log["log"] = append(s.Log["log"], line)
				}
				s.LogLocker.Unlock()
			}
			return
		case line := <-s.Msg:
			s.LogLocker.Lock()
			if len(s.Log["log"]) >= global.LogCount {
				copy(s.Log["log"], s.Log["log"][1:])
				s.Log["log"][global.LogCount-1] = line
			} else {
				s.Log["log"] = append(s.Log["log"], line)
			}
			s.LogLocker.Unlock()
		}
	}

}
