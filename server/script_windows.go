// +build windows

/*
 * @Author: your name
 * @Date: 2021-04-25 19:08:58
 * @LastEditTime: 2021-04-25 20:29:23
 * @LastEditors: your name
 * @Description: In User Settings Edit
 * @FilePath: /scs/script_windows.go
 */

package server

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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
			// 如果收到结束的信号，直接结束停止的goroutine
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
		return err
	}
	svc.stopStatus()
	return nil

}

func (svc *Server) start() error {
	svc.Status.Status = status.RUNNING
	for k, v := range svc.Env {
		// 替换command
		svc.Command = strings.ReplaceAll(svc.Command, "$"+k, v)
		svc.Command = strings.ReplaceAll(svc.Command, "${"+k+"}", v)
	}
	if _, err := os.Stat(svc.Script.Dir); os.IsNotExist(err) {
		golog.Error(err)
		return err
	}

	svc.Cmd = exec.Command("powershell", "-c", svc.Command)
	if svc.Cmd.Env == nil {
		svc.Cmd.Env = make([]string, 0, len(svc.Env))
	}
	for k, v := range svc.Env {
		if k == "" || v == "" {
			continue
		}
		svc.Cmd.Env = append(svc.Cmd.Env, k+"="+v)

	}
	svc.Cmd.Dir = svc.Script.Dir
	// 等待初始化完成完成后向后执行
	svc.read()
	svc.Status.Start = time.Now().Unix() // 设置启动状态是成功的
	if err := svc.Cmd.Start(); err != nil {
		// 执行脚本前的错误, 改变状态
		golog.Error(err)
		svc.stopStatus()
		return err
	}

	return nil
}
