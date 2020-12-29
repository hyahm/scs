// +build !windows

package script

import (
	"errors"
	"os/exec"
	"syscall"
	"time"

	"github.com/hyahm/golog"
)

func (s *Script) shell(command string) error {
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Env = s.Env

	read(cmd, s)

	err := cmd.Start()
	if err != nil {

		golog.Error(err)
		return err
	}
	defer func() {
		golog.Error(cmd.ProcessState.ExitCode())
	}()
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

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

func (s *Script) kill() {
	var err error
	err = syscall.Kill(-s.cmd.Process.Pid, syscall.SIGKILL)
	if err != nil {
		// 正常来说，不会进来的，特殊问题以后再说
		golog.Error(err)
	}
	s.stopStatus()
	return
}

func (s *Script) start() error {
	s.cmd = exec.Command("/bin/bash", "-c", s.Command)
	s.cmd.Dir = s.Dir
	s.cmd.Env = s.Env
	s.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	s.read()
	s.Status.Up = time.Now() // 设置启动状态是成功的
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
