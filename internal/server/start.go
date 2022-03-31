package server

import (
	"context"
	"strings"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/internal/server/status"
	"github.com/hyahm/scs/pkg"
)

// Start  启动服务 异步的
func (svc *Server) Start() {

	// svc.make(s * scripts.Script)
	switch svc.Status.Status {
	case status.WAITSTOP:
		// 如果之前是等待停止的状态， 更改为重启状态
		<-svc.Exit
		svc.Exit <- 10
		svc.Status.Status = status.WAITRESTART
	case status.STOP:
		// 开始启动的时候，需要将遍历变量值的模板渲染
		go svc.asyncStart()
	}
}

// 当是停止状态的时候异步启动
func (svc *Server) asyncStart() {
	svc.Version = pkg.GetVersion(svc.Version)
	err := svc.Install()
	if err != nil {
		golog.Error(err)
		svc.stopStatus()
		return
	}
	svc.Exit = make(chan int, 2)
	svc.CancelProcess = make(chan bool, 2)
	svc.Ctx, svc.Cancel = context.WithCancel(context.Background())
	if svc.Cron != nil && svc.Cron.Loop > 0 {
		svc.IsCron = true

		// 循环的起止时间可以只设置时分秒， 自动补齐今天的日期
		svc.Cron.Start = strings.Trim(svc.Cron.Start, " ")
		if svc.Cron.Start != "" {
			// 计算下次启动的时间, 不等于空就按照上面的时间来计算
			index := strings.Index(svc.Cron.Start, " ")
			if index < 0 {
				// 如果只有时间， 自动获取今天的年月日
				svc.Cron.Start = strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")[0] + " " + svc.Cron.Start
			}
			svc.Cron.StartTime, err = time.ParseInLocation("2006-01-02 15:04:05", svc.Cron.Start, time.Local)
			if err != nil {
				golog.Error(err)
			}
			// 计算下次启动时间

			svc.Cron.ComputerStartTime()
			// 比较是否过了时间点， 如果过了就重新计算， 否则就是给定的时间
		} else {
			svc.Cron.StartTime = time.Now()
			// 如果没设置， 设置下此启动的时间为当前时间
		}
		svc.Status.Start = time.Now().Unix() // 设置启动状态是成功的
		go svc.cron()
		return
	}
	svc.Command = internal.Format(svc.Command, svc.Env)
	golog.Info(svc.Command)
	if svc.Disable {
		svc.stopStatus()
		return
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
