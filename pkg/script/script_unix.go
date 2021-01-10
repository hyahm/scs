// +build !windows

package script

import (
	"errors"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/hyahm/golog"
)

func (s *Script) stop() {

	for {
		select {
		case <-time.After(time.Millisecond * 10):
			if !s.Status.CanNotStop {
				err := syscall.Kill(-s.cmd.Process.Pid, syscall.SIGKILL)
				if err != nil {
					// 如果pid已经被杀掉了， 那么就报错
					golog.Warnf("pid already be killed, err: %v", err)
				}

				golog.Debugf("stop %s\n", s.SubName)
				s.Status.RestartCount = 0
				return
			}
		case <-s.EndStop:
			// 如果收到结束的信号，直接结束停止的goroutine
			return
		}
	}

}

func (s *Script) kill() error {
	var err error
	err = syscall.Kill(-s.cmd.Process.Pid, syscall.SIGKILL)
	if err != nil {
		// 正常来说，不会进来的，特殊问题以后再说
		golog.Error(err)
		return err
	}
	s.stopStatus()
	return nil
}

func (s *Script) start() error {
	s.cmd = exec.Command("/bin/bash", "-c", s.Command)
	s.cmd.Dir = s.Dir
	if s.cmd.Env == nil {
		s.cmd.Env = make([]string, 0, len(s.Env))
	}
	for k, v := range s.Env {
		s.cmd.Env = append(s.cmd.Env, k+":"+v)
		s.Command = strings.ReplaceAll(s.Command, "${"+k+"}", v)
	}
	s.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	s.read()
	s.Status.Start = time.Now().Unix() // 设置启动状态是成功的
	if err := s.cmd.Start(); err != nil {
		// 执行脚本前的错误, 改变状态
		golog.Error(err)
		return err
	}

	if s.cmd.Process == nil {
		return errors.New("not running")
	}
	return nil
}
