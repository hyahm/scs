package server

import (
	"context"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/server/status"
	"github.com/hyahm/scs/pkg/config/alert"
	"github.com/hyahm/scs/pkg/message"
)

func (svc *Server) wait() {
	go svc.successAlert()
	ctx, cancel := context.WithCancel(context.Background())
	go svc.CheckReady(ctx)
	defer cancel()

	if err := svc.Cmd.Wait(); err != nil {
		// 脚本退出后才会执行这里的代码
		select {
		case ec := <-svc.Exit:
			switch ec {
			case 9:
				// stop操作
				// 返回一个停止的信号

				// 主动退出, kill， stop
				svc.stopStatus()
				return
			case 10:
				// 重启 restart
				svc.stopStatus()
				svc.Start()
				return
			}
		default:
			// 意外退出
			if !svc.DisableAlert {
				golog.Error(svc.SubName+": ", err.Error())
				if alert.HaveAlert() {
					am := &message.Message{
						Title:  "service error stop",
						Pname:  svc.Name,
						Name:   svc.SubName,
						Reason: err.Error(),
					}
					if !svc.AI.Broken {
						// 第一次
						svc.AI.Start = time.Now()
						am.BrokenTime = svc.AI.Start.String()
						svc.AI.AlertTime = time.Now()
						svc.AI.Broken = true
						alert.AlertMessage(am, svc.AT)
					} else {
						// 间隔时间内才发送报警
						ci := svc.ContinuityInterval
						if ci == 0 {
							ci = defaultContinuityInterval
						}
						if time.Since(svc.AI.AlertTime) >= ci {
							svc.AI.AlertTime = time.Now()
							alert.AlertMessage(am, svc.AT)
						}
					}
				}

			}
			// 如果是定时器的话， 直接结束
			if svc.Cron != nil && svc.Cron.Loop > 0 {
				svc.Status.Pid = 0
				svc.Logger.Close()
				return
			}
			if svc.Always {
				svc.Status.Status = status.STOP
				svc.Status.Pid = 0
				svc.Status.Start = 0
				// 失败了， 每秒启动一次
				svc.Logger.Close()
				svc.Status.RestartCount++
				time.Sleep(time.Second)
				svc.Start()
				return
			}
		}
		golog.Errorf("serviceName: %s, subScript: %s, error: %v \n", svc.Name, svc.SubName, err)
		svc.stopStatus()
		return
	}
	if svc.IsCron {
		// 如果是个定时器， 执行结束不停止
		return
	}
	if svc.Removed {
		svc.StopSigle <- true
	}
	// todo: if svc.DeleteWhenExit {
	// svc.DeleteWhenExitSingle <- true
	// }
	svc.stopStatus()

}
