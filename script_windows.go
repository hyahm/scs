// +build windows

package scs

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/hyahm/golog"
)

func (svc *Server) stop() {
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
	if svc.cmd == nil {
		return nil
	}
	err := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(svc.cmd.Process.Pid)).Run()
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
		// 替换command
		svc.Command = strings.ReplaceAll(svc.Command, "$"+k, v)
		svc.Command = strings.ReplaceAll(svc.Command, "${"+k+"}", v)
	}
	svc.cmd = exec.Command("powershell", "-c", svc.Command)
	if svc.cmd.Env == nil {
		svc.cmd.Env = make([]string, 0, len(svc.Env))
	}
	for k, v := range svc.Env {
		// 需要单独抽出去>>
		svc.cmd.Env = append(svc.cmd.Env, k+"="+v)

	}
	svc.cmd.Dir = svc.Script.Dir
	golog.Warn(svc.cmd.Dir)
	golog.Warn(svc.Command)
	// 等待初始化完成完成后向后执行
	svc.read()
	svc.Status.Start = time.Now().Unix() // 设置启动状态是成功的
	if err := svc.cmd.Start(); err != nil {
		// 执行脚本前的错误, 改变状态
		golog.Error(err)
		svc.stopStatus()
		return err
	}

	return nil
}
