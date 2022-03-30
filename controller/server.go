package controller

import (
	"encoding/json"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/scs/pkg/config/alert"
	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/scs/pkg/config/scripts/subname"
	"github.com/hyahm/scs/pkg/message"
	"github.com/hyahm/scs/status"
)

// 删除对应的server, 外部加了锁，内部调用不用加锁
func removeServer(name, subname string, update bool) {
	store.servers[subname].Logger.Close()
	delete(store.servers, subname)
	// 如果scripts的副本数为0或者1就直接删除这个scripts
	if update {
		if _, ok := store.ss[name]; ok {
			store.ss[name].Replicate--
			if store.ss[name].Replicate <= 0 {
				config.DeleteScriptToConfigFile(store.ss[name], update)
				delete(store.ss, name)
				return
			}
			// 这里修改配置文件减一
			config.UpdateScriptToConfigFile(store.ss[name], update)
		}
	}
}

func StartAllServer() {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, v := range store.servers {
		v.Start()
	}
}

func StartPermAllServer(token string) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, v := range store.servers {
		if v.Token == token {
			v.Start()
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

// 通过script 生成 server并启动服务
func makeReplicateServerAndStart(s *scripts.Script, end int) {
	// 防止删除的时候也给删掉， 直接给加上
	// ss[s.Name] = s
	// ss[s.Name].EnvLocker = &sync.RWMutex{}

	// availablePort := s.Port
	// for i := 0; i < end; i++ {
	// 	if _, ok := store.serverIndex[s.Name][i]; ok {
	// 		continue
	// 	}
	// 	// 根据副本数提取子名称
	// 	if store.serverIndex[s.Name] == nil {
	// 		store.serverIndex[s.Name] = make(map[int]struct{})
	// 	}
	// 	store.serverIndex[s.Name][i] = struct{}{}
	// 	env := make(map[string]string)
	// 	for k, v := range s.TempEnv {
	// 		env[k] = v
	// 	}
	// 	subname := subname.NewSubname(s.Name, i)
	// 	var svc *server.Server
	// 	if s.Port > 0 {
	// 		// 顺序拿到可用端口
	// 		availablePort = pkg.GetAvailablePort(availablePort)
	// 		env["PORT"] = strconv.Itoa(availablePort)
	// 		svc = s.Add(availablePort, end, i, subname)
	// 		availablePort++
	// 	} else {
	// 		env["PORT"] = "0"
	// 		svc = s.Add(0, end, i, subname)
	// 	}

	// 	env["NAME"] = subname.String()
	// 	svc.Env = env
	// 	servers[subname.String()] = svc
	// 	if !svc.Disable {
	// 		golog.Infof("start server: %s", svc.SubName)
	// 		servers[subname.String()].Start()
	// 	}

	// }
}

// 获取所有服务的状态
func All(role string) []byte {
	store.mu.RLock()
	defer store.mu.RUnlock()
	statuss := &pkg.StatusList{
		Data:    make([]status.ServiceStatus, 0),
		Version: global.VERSION,
		Role:    role,
	}
	for name := range store.servers {
		pname := subname.Subname(name).GetName()
		if v, ok := store.ss[pname]; ok {
			if !v.Disable {
				statuss.Data = append(statuss.Data, getStatus(pname, name))
			}
		}
	}
	statuss.Code = 200
	send, err := json.Marshal(statuss)
	if err != nil {
		golog.Error(err)
	}
	return send
}

// 获取所有服务的状态
func AllLook(role, token string) []byte {
	store.mu.RLock()
	defer store.mu.RUnlock()
	statuss := &pkg.StatusList{
		Data:    make([]status.ServiceStatus, 0),
		Version: global.VERSION,
		Role:    role,
	}
	for name, server := range store.servers {
		pname := subname.Subname(name).GetName()
		if v, ok := store.ss[pname]; ok {
			if server.Token == token && !v.Disable {
				statuss.Data = append(statuss.Data, getStatus(pname, name))
			}
		}
	}
	statuss.Code = 200
	send, err := json.Marshal(statuss)
	if err != nil {
		golog.Error(err)
	}
	return send
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
