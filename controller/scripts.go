package controller

import (
	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/scs/pkg/config/scripts/subname"
)

func GetScripts() map[string]*scripts.Script {
	store.mu.RLock()
	defer store.mu.RUnlock()
	return store.ss
}

func GetScriptsFromScritps(names map[string]struct{}) map[string]*scripts.Script {
	store.mu.RLock()
	defer store.mu.RUnlock()
	ss := make(map[string]*scripts.Script)
	for pname, script := range store.ss {
		if _, ok := names[pname]; ok {
			ss[pname] = script
		}
	}
	return ss
}

func GetIndexs(pname string) []int {
	indexs := make([]int, 0)
	store.mu.RLock()
	defer store.mu.RUnlock()
	for k := range store.serverIndex[pname] {
		indexs = append(indexs, k)
	}
	return indexs
}

func KillScript(s *scripts.Script) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(s.Name, i)
		store.servers[subname.String()].Kill()
	}
}

func NeedStop(s *scripts.Script) bool {
	// 更新server
	// 判断值是否相等
	return !scripts.EqualScript(s, store.ss[s.Name])
}

func GetScriptByPname(pname string) (*scripts.Script, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	v, ok := store.ss[pname]
	return v, ok
}
