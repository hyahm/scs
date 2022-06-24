//go:build !windows
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
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/server/status"
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
				return
			}
		case <-svc.CancelProcess:
			// 如果收到取消结束的信号，退出之前的操作
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
	return nil
}

func (svc *Server) start() error {
	svc.Cmd = exec.Command("/bin/bash", "-c", svc.Command)
	if svc.Dir != "" {
		if _, err := os.Stat(svc.Dir); os.IsNotExist(err) {
			golog.Error(err)
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
	svc.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true,
		Credential: &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		}}
	svc.read()
	svc.Status.Start = time.Now().Unix() // 设置启动状态是成功的
	if err := svc.Cmd.Start(); err != nil {
		// 执行脚本前的错误, 改变状态
		golog.Error(err)
		return err
	}
	svc.Status.Status = status.RUNNING
	return nil
}
