package server

import (
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/server/status"
	"github.com/hyahm/scs/pkg/config/alert"
	"github.com/hyahm/scs/pkg/message"
)

func (svc *Server) wait() {
	go svc.successAlert()
	// 这三行无视，
	// ctx, cancel := context.WithCancel(context.Background())
	// go svc.CheckReady(ctx)
	// defer cancel()
	// 只要是结束的。 都要修改状态， 关闭日志
	if err := svc.Cmd.Wait(); err != nil {
		// 脚本退出后才会执行这里的代码
		select {
		case ec := <-svc.Exit:
			switch ec {
			case 9:
				// stop操作
				// 返回一个停止的信号
				// 主动退出, kill
			case 10:
				// 重启 restart 感觉应该在外部重新makeserver
				svc.stopStatus()
				golog.Debug("ready send stop single")
				svc.StopSignal <- true
				return
			case 11:
				// 停止信号 stop
			case 12:
				// remove的信号
				svc.stopStatus()
				golog.Debug("ready send stop single")
				svc.StopSignal <- true
				return
			}

		default:
			// 意外退出的报警
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
						if time.Since(svc.AI.AlertTime) >= global.GeContinuityInterval() {
							svc.AI.AlertTime = time.Now()
							alert.AlertMessage(am, svc.AT)
						}
					}
				}

			}
			// 如果是定时器的话， 直接结束
			// if svc.Cron != nil && svc.Cron.Loop > 0 {
			// 	svc.Status.Pid = 0
			// 	svc.Logger.Close()
			// 	return
			// }
			if svc.Always {
				golog.Info("is always restart")
				// 如果总是启动的话， 再次拉起服务
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
			golog.Errorf("serviceName: %s, subScript: %s, error: %v \n", svc.Name, svc.SubName, err)
		}

	}
	// if svc.IsCron {
	// 	// 如果是个定时器， 执行结束不停止
	// 	return
	// }
	// if svc.Removed {

	// }
	// todo: if svc.DeleteWhenExit {
	// svc.DeleteWhenExitSingle <- true
	// }
	svc.stopStatus()

}
