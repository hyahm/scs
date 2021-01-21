package script

import (
	"time"

	"github.com/hyahm/scs/alert"
)

func (s *Script) successAlert() {
	// 启动成功后恢复的通知
	if !s.AI.Broken {
		return
	}
	for {
		select {
		// 每3秒一次操作
		case <-time.After(time.Second * 3):
			am := &alert.Message{
				Title:      "service recover",
				Pname:      s.Name,
				Name:       s.SubName,
				BrokenTime: s.AI.Start.String(),
				FixTime:    time.Now().String(),
			}
			alert.AlertMessage(am, s.AT)
			s.AI.Broken = false
			return
		case <-s.Ctx.Done():
			return
		}
	}

}
