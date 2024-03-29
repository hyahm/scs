package controller

import (
	"errors"
	"fmt"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg/config/scripts"
)

// 异步执行停止脚本
func StopScript(s *scripts.Script) error {
	if s.Disable {
		return nil
	}
	_, ok := store.Store.GetScriptByName(s.Name)
	if !ok {
		return errors.New("not found script: " + s.Name)
	}
	// 禁用 script 所在的所有server
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := fmt.Sprintf("%s_%d", s.Name, i)
		svc, ok := store.Store.GetServerByName(subname)
		if ok {
			golog.Debug("send all stop")
			go svc.Stop()
		}
	}
	golog.Debug("stop", s.Name)
	return nil
}
