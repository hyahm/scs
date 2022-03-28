package controller

import (
	"errors"

	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/scs/pkg/config/scripts/subname"
)

// 异步执行停止脚本
func StopScript(s *scripts.Script) error {
	mu.RLock()
	defer mu.RUnlock()
	if _, ok := ss[s.Name]; !ok {
		return errors.New("")
	}
	// 禁用 script 所在的所有server
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(s.Name, i)
		go servers[subname.String()].Stop()
	}
	return nil
}
