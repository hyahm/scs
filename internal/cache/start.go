package cache

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/pkg"
)

// Start  启动服务 异步的
func (svc *Server) StartAsync() {

	// parameter := ""
	switch svc.Status.Status {
	case pkg.STOP:
		// 开始启动的时候，需要将遍历变量值的模板渲染
		go svc.Start()
	}
}

func FirstStartAllScript() {
	for _, script := range internal.GetAllScript() {
		// 如果没设置token， 默认生成一个脚本的token
		// 默认修改script
		AddScript(script)
	}
}

// 当是停止状态的时候异步启动
func (svc *Server) Start() {
	svc.Logger = golog.NewLog(
		filepath.Join(internal.GetLogPath(), svc.SubName+".log"), 0, true)
	// golog 的 (*Log).Sync() 会无条件 close 内部 channel，重复调用会 panic，
	// 这里用 syncOnce 保证同一个 Logger 实例只 Sync 一次
	defer svc.syncLogger()
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
	// svc.CancelProcess = make(chan bool, 2)

	svc.Status.Command = internal.Format(svc.Command, svc.Env)
	svc.Ctx, svc.Cancel = context.WithCancel(context.Background())
	defer svc.Cancel()
	// go func() {
	// 	if svc.StopTime != "" {
	// 		stopTime, err := time.ParseInLocation(time.DateTime, svc.StopTime, time.Local)
	// 		if err != nil {
	// 			golog.Warnf("parse stop time failed: %v", err)
	// 			return
	// 		}
	// 		if time.Since(stopTime).Seconds() < 0 {
	// 			duration := -time.Since(stopTime)
	// 			select {
	// 			case <-time.After(duration):
	// 				svc.Stop()
	// 				svc.Cancel()
	// 				return
	// 			case <-svc.Ctx.Done():
	// 				return
	// 			}
	// 		}
	// 	}
	// }()
	if svc.Cron.Loop > 0 {
		// 启动定时
		golog.Info("name:" + svc.SubName + " start cron")
		svc.IsCron = true
		svc.Status.Status = pkg.RUNNING
		// 循环的起止时间可以只设置时分秒， 自动补齐今天的日期
		svc.Cron.Start = strings.Trim(svc.Cron.Start, " ")
		svc.cron()
		return
	}
	// if svc.StartTime != "" {
	// 	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", svc.StartTime, time.Local)
	// 	if err == nil && time.Since(startTime).Seconds() < 0 {
	// 		for {
	// 			select {
	// 			case <-time.After(-1 * time.Since(startTime)):
	// 				svc.Stop()
	// 				svc.Cancel()
	// 				return
	// 			case <-svc.Ctx.Done():
	// 				return
	// 			}
	// 		}
	// 	}
	// }
	svc.Status.Start = time.Now().Unix() // 设置启动状态是成功的
	if err := svc.start(); err != nil {
		golog.Info(err)
		svc.stopStatus()
		return
	}
	if svc.Cmd.Process != nil {
		svc.Status.Pid = svc.Cmd.Process.Pid
	}
	svc.wait()

}

// syncLogger 安全地刷新日志。
// golog 的 (*Log).Sync() 会无条件 close 内部 channel，重复调用会触发
// "close of closed channel" panic。这里用 recover 兜底，避免在
// 定时任务/重启等路径下对同一 Logger 实例重复 Sync 导致整个进程崩溃。
func (svc *Server) syncLogger() {
	if svc.Logger == nil {
		return
	}
	defer func() {
		_ = recover()
	}()
	svc.Logger.Sync()
}
