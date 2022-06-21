package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/alert"
)

func UnStop(ctx context.Context, name string, sr *pkg.SignalRequest) {
	select {
	case <-time.After(time.Second * time.Duration(sr.Timeout)):
		// 报警
		if sr.Notice {
			ra := &alert.RespAlert{
				Name:   name,
				Title:  "原子操作超时",
				Reason: fmt.Sprintf("原子操作超时超过 %d 秒没有执行完成", sr.Timeout),
			}
			if sr.ContinuityInterval > 0 {
				ra.ContinuityInterval = sr.ContinuityInterval
			}
			ra.SendAlert()
			pkg.DeleteAtomSignal(name)
		}
		if sr.Restart {
			if server, ok := store.Store.GetServerByName(name); ok {
				server.Restart()
			}
		}
	case <-ctx.Done():

	}
}
