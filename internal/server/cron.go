package server

import (
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/server/status"
)

func (svc *Server) getNextDuration() time.Duration {
	now := time.Now()
	if svc.Cron.StartTime.After(now) {
		// 下次执行时间在未来
		return svc.Cron.StartTime.Sub(now)
	}
	// 下次执行时间已过，计算下次
	return -now.Sub(svc.Cron.StartTime)
}

func (svc *Server) cron() {
	svc.Times = svc.Cron.Times
	for {
		select {
		case <-svc.Ctx.Done():
			golog.Info("name:" + svc.SubName + " end cron")
			golog.Debug("ready send stop single")
			svc.StopSignal <- true
			return
		case <-time.After(svc.getNextDuration()):
			svc.Status.Status = status.RUNNING
			svc.Status.Start = time.Now().Unix()
			golog.Infof("cron start: %s time: %v", svc.SubName, svc.Cron.StartTime)
			if err := svc.start(); err != nil {
				golog.Error("cron start error: ", err)
				golog.Error(err)
				svc.stopStatus()
				return
			}
			if svc.Cmd != nil && svc.Cmd.Process != nil {
				svc.Status.Pid = svc.Cmd.Process.Pid
			}
			svc.wait()
			svc.Times--
			if svc.Cron.Times > 0 && svc.Times <= 0 {
				golog.Infof("循环器%s执行次数结束", svc.SubName)
				svc.stopStatus()
				return
			}
			svc.Status.Pid = 0
			svc.Cron.ComputerStartTime()
		}
	}
}
