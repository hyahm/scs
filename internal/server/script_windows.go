//go:build windows
// +build windows

package server

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/status"
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
				// 通知外部已经停止了
				svc.StopSigle <- true
				return
			}
		case <-svc.CancelProcess:
			// 如果收到取消结束的信号，退出之前的操作
			return
		}
	}
}

func (svc *Server) kill() error {
	if svc.Cmd == nil {
		return nil
	}
	err := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(svc.Cmd.Process.Pid)).Run()
	if err != nil {
		// 正常来说，不会进来的，特殊问题以后再说
		golog.Error(err)
	}
	svc.stopStatus()
	return nil

}

func (svc *Server) start() error {

	svc.Cmd = exec.Command("powershell", "-c", svc.Command)
	if svc.Dir != "" {
		if _, err := os.Stat(svc.Dir); os.IsNotExist(err) {
			return err
		}
		svc.Cmd.Dir = svc.Dir
	}
	if svc.Cmd.Env == nil {
		svc.Cmd.Env = make([]string, 0, len(svc.Env))
	}
	for k, v := range svc.Env {
		if k == "" || v == "" {
			continue
		}
		svc.Cmd.Env = append(svc.Cmd.Env, k+"="+v)

	}

	// 等待初始化完成完成后向后执行
	svc.read()

	if err := svc.Cmd.Start(); err != nil {
		// 执行脚本前的错误, 改变状态
		golog.Error(err)
		svc.stopStatus()
		return err
	}
	svc.Status.Status = status.RUNNING
	return nil
}
