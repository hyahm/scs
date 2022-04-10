package controller

import (
	"fmt"
	"sync/atomic"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/scs/pkg/config/scripts/subname"
)

// 更新的操作
func DisableScript(s *scripts.Script, update bool) bool {

	// 禁用 script 所在的所有server
	if _, ok := store.ss[s.Name]; !ok {
		return false
	}
	if store.ss[s.Name].Disable {
		return false
	}
	store.mu.Lock()
	store.ss[s.Name].Disable = true
	store.mu.Unlock()
	for i := range store.serverIndex[s.Name] {
		subname := fmt.Sprintf("%s_%d", s.Name, i)

		if i == 0 {
			// 如果索引时0的， 那么直接停止就好了， 并且将值修改为true
			store.mu.Lock()
			store.servers[subname].Disable = true
			store.mu.Unlock()
			go store.servers[subname].Stop()
			continue
		}
		golog.Info("add reload count")
		atomic.AddInt64(&global.CanReload, 1)
		go Remove(store.servers[subname], update)

	}
	return true
}

// enable script
func EnableScript(script *scripts.Script) bool {
	store.mu.Lock()
	defer store.mu.Unlock()
	// 禁用 script 所在的所有server
	if _, ok := store.ss[script.Name]; !ok {
		return false
	}
	if !store.ss[script.Name].Disable {
		// 如果本身就是 启用的 不做任何操作
		return false
	}
	store.ss[script.Name].Disable = false

	AddScript(store.ss[script.Name])
	replicate := script.Replicate
	if replicate == 0 {
		replicate = 1
	}
	availablePort := script.Port
	for i := 0; i < replicate; i++ {
		subname := fmt.Sprintf("%s_%d", script.Name, i)
		store.servers[subname] = &server.Server{
			Index:     i,
			Replicate: replicate,
			SubName:   subname,
			Name:      script.Name,
		}
		store.serverIndex[script.Name][i] = struct{}{}
		availablePort = store.servers[subname].MakeServer(script, availablePort)
		availablePort++
		if script.Disable {
			// 如果是禁用的 ，那么不用生成多个副本，直接执行下一个script
			return true
		}

		store.servers[subname].Start()
	}
	return true
}

// func makeAndStart(i, replicate, availablePort int, script *scripts.Script) int {
// 	subname := fmt.Sprintf("%s_%d", script.Name, i)
// 	store.servers[subname] = &server.Server{
// 		Index:     i,
// 		Replicate: replicate,
// 		SubName:   subname,
// 		Name:      script.Name,
// 	}
// 	store.serverIndex[script.Name][i] = struct{}{}
// 	availablePort = store.servers[subname].MakeServer(script, availablePort)
// 	availablePort++
// 	if script.Disable {
// 		// 如果是禁用的 ，那么不用生成多个副本，直接执行下一个script
// 		return availablePort
// 	}

// 	store.servers[subname].Start()
// 	return availablePort
// }

func UpdateAndRestart(svc *server.Server) {
	svc.UpdateServer()
	store.mu.RLock()
	store.mu.RUnlock()
	defer store.mu.RUnlock()
	_, ok := store.ss[svc.Name]
	if !ok {
		return
	}
	updateAndRestartScript(store.ss[svc.Name])
}

// 返回成功还是失败
func UpdateAndRestartScript(s *scripts.Script) {
	store.mu.RLock()

	defer store.mu.RUnlock()
	_, ok := store.ss[s.Name]
	if !ok {
		return
	}
	updateAndRestartScript(s)
}

func updateAndRestartScript(s *scripts.Script) {
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	go func() {
		for i := 0; i < replicate; i++ {
			subname := subname.NewSubname(s.Name, i)
			store.servers[subname.String()].UpdateServer()
		}
		RestartScript(s)
	}()

}

func UpdateAllServer() {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, s := range store.ss {
		updateAndRestartScript(s)
	}
}

func UpdateAllServerFromScript(names map[string]struct{}) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, s := range store.ss {
		if _, ok := names[s.Name]; ok {
			go updateAndRestartScript(s)
		}

	}
}
