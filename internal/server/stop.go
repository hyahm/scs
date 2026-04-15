package server

import (
	"time"

	"github.com/hyahm/scs/internal/server/status"
)

func (svc *Server) stop() {
	ticker := time.NewTicker(time.Millisecond * 10)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if !svc.Status.CanNotStop {
				svc.kill()
				// 通知外部已经停止了
				return
			}
		case <-svc.CancelProcess:
			// 如果收到取消结束的信号，退出停止信号
			return
		}
	}
}

func (svc *Server) remove() {
	ticker := time.NewTicker(time.Millisecond * 10)
	defer ticker.Stop()
	svc.Status.Status = status.REMOVING
	for {
		select {
		case <-ticker.C:
			if !svc.Status.CanNotStop {
				svc.kill()
				// 通知外部已经停止了
				return
			}
		case <-svc.CancelProcess:
			// 如果收到取消结束的信号，退出停止信号
			return
		}
	}
}
