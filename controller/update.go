package controller

import (
	"fmt"
	"sync/atomic"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg/config/scripts"
)

// 更新的操作
func DisableScript(s *scripts.Script, update bool) bool {

	// 禁用 script 所在的所有server
	script, ok := store.Store.GetScriptByName(s.Name)
	if !ok {
		return false
	}
	if script.Disable {
		return false
	}
	script.Disable = true
	for i := range store.Store.GetScriptIndex(s.Name) {
		subname := fmt.Sprintf("%s_%d", s.Name, i)
		svc, ok := store.Store.GetServerByName(subname)
		if !ok {
			golog.Error("严重错误， 请提交问题到https://github.com/hyahm/scs")
		}
		if i == 0 {
			// 如果索引时0的， 那么直接停止就好了， 并且将值修改为true

			svc.Disable = true

			go svc.Stop()
			continue
		}
		golog.Info("add reload count")
		atomic.AddInt64(&global.CanReload, 1)
		go Remove(svc, update)

	}
	return true
}

// enable script
func EnableScript(script *scripts.Script) bool {
	// 禁用 script 所在的所有server
	script, ok := store.Store.GetScriptByName(script.Name)
	if !ok {
		return false
	}
	if !script.Disable {
		// 如果本身就是 启用的 不做任何操作
		return false
	}

	script.Disable = false

	AddScript(script)
	replicate := script.Replicate
	if replicate == 0 {
		replicate = 1
	}
	availablePort := script.Port
	for i := 0; i < replicate; i++ {
		subname := fmt.Sprintf("%s_%d", script.Name, i)
		store.Store.InitServer(i, replicate, script.Name, subname)
		store.Store.SetScriptIndex(script.Name, i)
		svc, _ := store.Store.GetServerByName(subname)
		availablePort = svc.MakeServer(script, availablePort)
		availablePort++
		if script.Disable {
			// 如果是禁用的 ，那么不用生成多个副本，直接执行下一个script
			return true
		}

		svc.Start()
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
