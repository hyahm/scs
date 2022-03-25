package controller

import (
	"github.com/hyahm/scs/internal/config/scripts"
	"github.com/hyahm/scs/internal/config/scripts/subname"
)

// 更新的操作
func DisableScript(s *scripts.Script) bool {
	mu.Lock()
	defer mu.Unlock()
	// 禁用 script 所在的所有server
	if _, ok := ss[s.Name]; !ok {
		return false
	}

	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}

	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(s.Name, i)
		if _, ok := servers[subname.String()]; ok {
			go Remove(servers[subname.String()])
		}

	}
	return true
}

func EnableScript(s *scripts.Script) bool {
	mu.Lock()
	defer mu.Unlock()
	// 禁用 script 所在的所有server
	if _, ok := ss[s.Name]; !ok {
		return false
	}

	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		// subname := subname.NewSubname(s.Name, i)
		makeReplicateServerAndStart(ss[s.Name], replicate)
		// go servers[subname.String()].Start()
	}
	return true
}
