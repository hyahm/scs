package controller

import (
	"github.com/hyahm/scs/internal/config/scripts"
	"github.com/hyahm/scs/internal/config/scripts/subname"
)

func WaitKillAllServer() {
	// ss.ScriptLocker.RLock()
	// defer ss.ScriptLocker.RUnlock()
	for _, s := range ss {
		WaitKillScript(s)
	}
}

// 同步杀掉
func WaitKillScript(s *scripts.Script) {
	// ss.ServerLocker.RLock()
	// defer ss.ServerLocker.RUnlock()
	// 禁用 script 所在的所有server
	// mu.RLock()
	// defer mu.RUnlock()
	// 禁用 script 所在的所有server
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(s.Name, i)
		servers[subname.String()].Kill()
	}
}
