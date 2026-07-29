package cache

import (
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/pkg"
)

func (svc *Server) cron() {
	ticker := time.NewTicker(time.Duration(svc.Cron.Loop) * time.Second)
	defer ticker.Stop()
	if svc.Cron.Start != "" {
		start, err := time.ParseInLocation("2006-01-02 15:04:05", svc.Cron.Start, time.Local)
		if err != nil {
			golog.Error(err)
			return
		}
		reset := time.Until(start)
		if reset > 0 {
			ticker.Reset(reset)
		}
	}
	svc.doTicker()
	// 计算下次启动时间

	svc.Times = svc.Cron.Times
	for {
		select {
		case <-svc.Ctx.Done():
			golog.Info("name:" + svc.SubName + " end cron")
			svc.resetStatus()
			// svc.StopSignal <- true
			return
		case <-ticker.C:
			svc.doTicker()
			ticker.Reset(time.Duration(svc.Cron.Loop) * time.Second)
		}
	}
}

func (svc *Server) doTicker() {
	svc.Status.Status = pkg.RUNNING
	svc.Status.Start = time.Now().Unix()
	if err := svc.start(); err != nil {
		svc.Logger.Error("cron start error: ", err)
		svc.resetStatus()
		return
	}
	if svc.Cmd != nil && svc.Cmd.Process != nil {
		svc.Status.Pid = svc.Cmd.Process.Pid
	}
	svc.wait()
	svc.Times--
	if svc.Cron.Times > 0 && svc.Times <= 0 {
		svc.Logger.Infof("循环器%s执行次数结束", svc.SubName)
		svc.resetStatus()
		return
	}
	svc.Status.Pid = 0
}
