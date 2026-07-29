package cache

import "github.com/hyahm/golog"

func (svc *Server) wait() {
	go svc.successAlert()
	// 只要是结束的，都要修改状态
	if err := svc.Cmd.Wait(); err != nil {
		// 不正常推出才重置
		golog.Warn(err, "-->", svc.SubName)
		svc.resetStatus()
	}
	if svc.IsCron {
		// 定时器执行结束不停止，状态保持 RUNNING（等待下次执行）
		svc.Status.Pid = 0
		return
	}
	// 锁内原子决策：进程退出后是否自动重启（由 Restart 设置 restartOnExit，此处消费）
	// svc.mu.Lock()
	// restart := svc.restartOnExit
	// svc.restartOnExit = false
	// svc.mu.Unlock()

	// if svc.restartOnExit {
	// 	// 重新拉起：StartAsync 仅在 Status==STOP 时启动
	// 	svc.StartAsync()
	// }
}
