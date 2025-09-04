package controller

import (
	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/scs/pkg/config/alert"
	"github.com/hyahm/scs/pkg/message"
)

// 删除对应的server, 外部加了锁，内部调用不用加锁  todo:
func removeServer(name, subname string, update bool) {
	// 如果scripts的副本数为0或者1就直接删除这个scripts

	script, ok := store.Store.GetScriptByName(name)
	if !ok {
		return
	}

	script.Replicate--
	if script.Replicate <= 0 {
		config.DeleteScriptToConfigFile(script, update)
		return
	}
	// 这里修改配置文件减一
	config.UpdateScriptToConfigFile(script, update)
}

func StartAllServer() {
	for _, svc := range store.Store.GetAllServer() {
		svc.Start()
	}
}

func StartAllServerFromScript(names map[string]struct{}) {
	for _, svc := range store.Store.GetAllServer() {
		svc.Start()
	}
}

// func HaveScript(pname string) bool {
// 	store.mu.RLock()
// 	defer store.mu.RUnlock()
// 	_, ok := store.ss[pname]
// 	return ok
// }

// func GetServers() map[string]*server.Server {
// 	store.mu.RLock()
// 	defer store.mu.RUnlock()
// 	return store.servers
// }

func GetServersFromScripts(names map[string]struct{}) map[string]*server.Server {
	servers := make(map[string]*server.Server)
	for name, svc := range store.Store.GetAllServerMap() {
		if _, ok := names[svc.Name]; ok {
			servers[name] = svc
		}
	}
	return servers
}

func GetAterts() map[string]message.SendAlerter {
	return alert.GetAlerts()
}

func StopScriptFromName(names map[string]struct{}) {
	for _, script := range store.Store.GetScriptMapFilterByName(names) {
		err := StopScript(script)
		if err != nil {
			golog.Error(err)
		}
	}
}

func StopAllServer() {
	for _, script := range store.Store.GetAllScriptMap() {
		err := StopScript(script)
		if err != nil {
			golog.Error(err)
		}
	}
}
