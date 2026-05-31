package server

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/internal/server/status"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
)

// Start  启动服务 异步的
func (svc *Server) Start() {
	// parameter := ""
	switch svc.Status.Status {
	case status.STOP:
		// 开始启动的时候，需要将遍历变量值的模板渲染
		go svc.asyncStart()
	}
}

// 当是停止状态的时候异步启动
func (svc *Server) asyncStart() {
	svc.Logger = golog.NewLog(
		filepath.Join(config.Cfg.Log.Path, svc.SubName+".log"), 0, true)

	// svc.Env["PARAMETER"] = param
	// 格式化 SCS_TPL 开头的环境变量
	for k := range svc.Env {
		if len(k) > 8 && k[:7] == "SCS_TPL" {
			svc.Env[k] = internal.Format(svc.Env[k], svc.Env)
		}
	}
	svc.Always = svc.AlwaysSign
	svc.Version = pkg.GetVersion(svc.Version)
	err := svc.Install()
	if err != nil {
		golog.Error(err)
		svc.stopStatus()
		return
	}
	svc.Exit = make(chan int, 2)
	svc.CancelProcess = make(chan bool, 2)

	svc.Status.Command = internal.Format(svc.Command, svc.Env)
	svc.Ctx, svc.Cancel = context.WithCancel(context.Background())
	go func() {
		if svc.StopTime != "" {
			stopTime, err := time.ParseInLocation(time.DateTime, svc.StopTime, time.Local)
			if err != nil {
				golog.Warnf("parse stop time failed: %v", err)
				return
			}
			if time.Since(stopTime).Seconds() < 0 {
				duration := -time.Since(stopTime)
				select {
				case <-time.After(duration):
					svc.Stop()
					svc.Cancel()
					return
				case <-svc.Ctx.Done():
					return
				}
			}
		}

	}()
	if svc.Cron != nil && svc.Cron.Loop > 0 {
		golog.Info("name:" + svc.SubName + " start cron")
		svc.IsCron = true
		svc.Status.Status = status.RUNNING
		// 循环的起止时间可以只设置时分秒， 自动补齐今天的日期
		svc.Cron.Start = strings.Trim(svc.Cron.Start, " ")

		go svc.cron()
		return
	}
	if svc.StartTime != "" {
		startTime, err := time.ParseInLocation("2006-01-02 15:04:05", svc.StartTime, time.Local)
		if err == nil && time.Since(startTime).Seconds() < 0 {
			for {
				select {
				case <-time.After(-1 * time.Since(startTime)):
					svc.Stop()
					svc.Cancel()
					return
				case <-svc.Ctx.Done():
					return
				}
			}
		}
	}
	svc.Status.Start = time.Now().Unix() // 设置启动状态是成功的
	if err := svc.start(); err != nil {
		golog.Info(err)
		svc.stopStatus()
		return
	}

	go svc.wait()

	if svc.Cmd.Process != nil {
		svc.Status.Pid = svc.Cmd.Process.Pid
	}
}
