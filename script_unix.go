// +build !windows

package scs

import (
	"errors"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/hyahm/golog"
)

func (svc *Server) stop() {
	defer func() {
		if err := recover(); err != nil {
			golog.Error(err)
		}
	}()
	for {
		select {
		case <-time.After(time.Millisecond * 10):
			if !svc.Status.CanNotStop {
				err := svc.kill()
				if err != nil {
					golog.Error(err)
					return
				}
				svc.StopSigle <- true
				return
			}
		case <-svc.CancelProcess:
			// 如果收到结束的信号，直接结束停止的goroutine
			return
		}
	}

}

func (svc *Server) kill() error {
	var err error
	err = syscall.Kill(-svc.cmd.Process.Pid, syscall.SIGKILL)
	if err != nil {
		// 正常来说，不会进来的，特殊问题以后再说
		golog.Error(err)
		return err
	}
	svc.stopStatus()
	return nil
}

func (svc *Server) start() error {
	golog.Info("already start")
	svc.Status.Status = RUNNING
	for k, v := range svc.Env {
		svc.Command = strings.ReplaceAll(svc.Command, "$"+k, v)
		svc.Command = strings.ReplaceAll(svc.Command, "${"+k+"}", v)
	}
	svc.cmd = exec.Command("/bin/bash", "-c", svc.Command)
	svc.cmd.Dir = svc.Script.Dir
	if svc.cmd.Env == nil {
		svc.cmd.Env = make([]string, 0, len(svc.Env))
	}
	for k, v := range svc.Env {
		svc.cmd.Env = append(svc.cmd.Env, k+"="+v)
	}
	svc.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	golog.Warn(svc.cmd.Dir)
	golog.Warn(svc.Command)
	svc.read()
	svc.Status.Start = time.Now().Unix() // 设置启动状态是成功的
	if err := svc.cmd.Start(); err != nil {
		// 执行脚本前的错误, 改变状态
		golog.Error(err)
		return err
	}

	if svc.cmd.Process == nil {
		return errors.New("not running")
	}
	return nil
}
