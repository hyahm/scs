//go:build windows
// +build windows

package server

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/server/status"
)

func (svc *Server) kill() {
	if svc.Cmd == nil {
		golog.Info(svc.Cmd)
		return
	}
	err := exec.Command("powershell", "/C", "taskkill", "/F", "/T", "/PID", fmt.Sprint(svc.Cmd.Process.Pid)).Run()
	if err != nil {
		// 正常来说，不会进来的，特殊问题以后再说
		golog.Error(err)
	}
}

func (svc *Server) start() error {
	golog.Info("server command: ", svc.Status.Command)
	svc.Cmd = exec.Command("powershell", "-c", svc.Status.Command)
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
		if k == "" {
			continue
		}
		// 去掉 key、value 里的 NUL 和控制字符
		v = strings.ReplaceAll(v, "\x00", "")
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
