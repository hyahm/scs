// +build windows

package script

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/server/script/status"
)

func (s *Server) stop() {
	defer func() {
		if err := recover(); err != nil {
			// 防止异常
			golog.Error(err)
		}
	}()
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
				// 清空失败次数
				s.stopStatus()
				s.Status.RestartCount = 0
				return
			}
		case <-s.CancelSigle:
			// 如果收到结束的信号，直接结束停止的goroutine
			return
		}
	}
}

func (s *Server) kill() error {
	err := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(s.cmd.Process.Pid)).Run()
	if err != nil {
		// 正常来说，不会进来的，特殊问题以后再说
		golog.Error(err)
		return err
	}
	s.stopStatus()
	return nil

}

func (s *Server) start() error {
	golog.Info(s.Command)
	golog.Info(s.Dir)
	s.Status.Status = status.RUNNING
	s.cmd = exec.Command("powershell", "-c", s.Command)
	if s.cmd.Env == nil {
		s.cmd.Env = make([]string, 0, len(s.Env))
	}

	for k, v := range s.Env {
		if k == "" {
			continue
		}
		// 需要单独抽出去>>
		s.cmd.Env = append(s.cmd.Env, k+"="+v)
		s.Command = strings.ReplaceAll(s.Command, "$"+k, v)
		s.Command = strings.ReplaceAll(s.Command, "${"+k+"}", v)

	}

	s.cmd.Dir = s.Dir
	// 等待初始化完成完成后向后执行
	s.read()
	s.Status.Start = time.Now().Unix() // 设置启动状态是成功的
	if err := s.cmd.Start(); err != nil {
		// 执行脚本前的错误, 改变状态
		golog.Error(err)
		s.stopStatus()
		return err
	}

	return nil
}
