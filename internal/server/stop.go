package server

import (
	"time"

	"github.com/hyahm/scs/internal/server/status"
)

func (svc *Server) stop() {
	for {
		select {
		case <-time.After(time.Millisecond * 10):
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
	svc.Status.Status = status.REMOVING
	for {
		select {
		case <-time.After(time.Millisecond * 10):
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
