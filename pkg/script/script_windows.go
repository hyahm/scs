// +build windows

package script

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/hyahm/golog"
)

func (s *Script) shell(command string) error {
	cmd := exec.Command("cmd", "/c", command)
	cmd.Env = s.Env
	read(cmd, s)
	err := cmd.Start()
	if err != nil {
		golog.Error(err)
		return err
	}
	return cmd.Wait()
}

func (s *Script) stop() {
	for {
		select {
		case <-time.After(time.Millisecond * 10):
			if !s.Status.CanNotStop {
				err := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(s.cmd.Process.Pid)).Run()
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
	err := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(s.cmd.Process.Pid)).Run()
	if err != nil {
		// 正常来说，不会进来的，特殊问题以后再说
		golog.Error(err)
		return err
	}
	s.stopStatus()
	return nil

}

func (s *Script) start() error {
	golog.Info(s.Command)
	s.cmd = exec.Command("cmd", "/C", s.Command)
	// 需要单独抽出去>>
	s.cmd.Env = s.Env
	s.cmd.Dir = s.Dir
	// 等待初始化完成完成后向后执行
	s.read()
	s.Status.Up = time.Now() // 设置启动状态是成功的
	if err := s.cmd.Start(); err != nil {
		// 执行脚本前的错误, 改变状态
		golog.Error(err)
		s.stopStatus()
		return err
	}

	return nil
}
