// +build !windows

/*
 * @Author: cander
 * @Date: 2021-04-25 19:08:58
 * @LastEditTime: 2021-04-25 20:28:33
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /scs/script_unix.go
 */

package server

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/status"
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
	err = syscall.Kill(-svc.Cmd.Process.Pid, syscall.SIGKILL)
	if err != nil {
		// 正常来说，不会进来的，特殊问题以后再说
		golog.Error(err)
		return err
	}
	svc.stopStatus()
	return nil
}

func (svc *Server) start() error {
	svc.Status.Status = status.RUNNING
	for k, v := range svc.Env {
		svc.Command = strings.ReplaceAll(svc.Command, "$"+k, v)
		svc.Command = strings.ReplaceAll(svc.Command, "${"+k+"}", v)
	}

	svc.Cmd = exec.Command("/bin/bash", "-c", svc.Command)
	if svc.Script.Dir != "" {
		if _, err := os.Stat(svc.Script.Dir); os.IsNotExist(err) {
			golog.Error(err)
			return err
		}
		svc.Cmd.Dir = svc.Script.Dir
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
	svc.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	svc.read()
	svc.Status.Start = time.Now().Unix() // 设置启动状态是成功的
	if err := svc.Cmd.Start(); err != nil {
		// 执行脚本前的错误, 改变状态
		golog.Error(err)
		return err
	}

	if svc.Cmd.Process == nil {
		return errors.New("not running")
	}
	return nil
}
