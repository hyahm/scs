package controller

import (
	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/scs/pkg/config/alert"
	"github.com/hyahm/scs/pkg/message"
)

// 删除对应的server, 外部加了锁，内部调用不用加锁  todo:
func removeServer(name, subname string, update bool) {
	// 如果scripts的副本数为0或者1就直接删除这个scripts
	if _, ok := store.ss[name]; ok {
		store.ss[name].Replicate--
		if store.ss[name].Replicate <= 0 {
			config.DeleteScriptToConfigFile(store.ss[name], update)
			return
		}
		// 这里修改配置文件减一
		config.UpdateScriptToConfigFile(store.ss[name], update)
	}
}

func StartAllServer() {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for name := range store.servers {
		store.servers[name].Start()

	}
}

func StartAllServerFromScript(names map[string]struct{}) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for name, svc := range store.servers {
		if _, ok := names[svc.Name]; ok {
			store.servers[name].Start()
		}
	}
}

func HaveScript(pname string) bool {
	store.mu.RLock()
	defer store.mu.RUnlock()
	_, ok := store.ss[pname]
	return ok
}

func GetServers() map[string]*server.Server {
	store.mu.RLock()
	defer store.mu.RUnlock()
	return store.servers
}

func GetServersFromScripts(names map[string]struct{}) map[string]*server.Server {
	store.mu.RLock()
	defer store.mu.RUnlock()
	servers := make(map[string]*server.Server)
	for name, svc := range store.servers {
		if _, ok := names[svc.Name]; ok {
			servers[name] = svc
		}
	}
	return servers
}

func GetPremServers(token string) map[string]*server.Server {
	store.mu.RLock()
	defer store.mu.RUnlock()
	tempServers := make(map[string]*server.Server)
	for name, v := range store.servers {
		if v.Token == token {
			tempServers[name] = v
		}
	}
	return tempServers
}

func GetAterts() map[string]message.SendAlerter {
	return alert.GetAlerts()
}

func StopScriptFromName(names map[string]struct{}) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, script := range store.ss {
		if _, ok := names[script.Name]; ok {
			err := StopScript(script)
			if err != nil {
				golog.Error(err)
			}
		}
	}
}

func StopAllServer() {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, s := range store.ss {
		err := StopScript(s)
		if err != nil {
			golog.Error(err)
		}
	}
}

func StopPermAllServer(token string) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, s := range store.ss {
		if s.Token == token {
			err := StopScript(s)
			if err != nil {
				golog.Error(err)
			}
		}

	}
}

func GetServerInfo(name string) (*server.Server, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	v, ok := store.servers[name]
	return v, ok
}
