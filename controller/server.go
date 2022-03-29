package controller

import (
	"encoding/json"
	"strconv"
	"sync"

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
func removeServer(name subname.Subname, update bool) {
	servers[name.String()].Logger.Close()
	delete(servers, name.String())
	// 如果scripts的副本数为0或者1就直接删除这个scripts
	pname := name.GetName()
	if update {
		if _, ok := ss[pname]; ok {
			ss[pname].Replicate--
			if ss[pname].Replicate <= 0 {
				config.DeleteScriptToConfigFile(ss[pname], update)
				delete(ss, pname)
				return
			}
			// 这里修改配置文件减一
			config.UpdateScriptToConfigFile(ss[pname], update)
		}
	}

}

func StartAllServer() {
	mu.RLock()
	defer mu.RUnlock()
	for _, v := range servers {
		v.Start()
	}
}

func StartPermAllServer(token string) {
	mu.RLock()
	defer mu.RUnlock()
	for _, v := range servers {
		if v.Token == token {
			v.Start()
		}

	}
}

func HaveScript(pname string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := ss[pname]
	return ok
}

func GetServers() map[string]*server.Server {
	mu.RLock()
	defer mu.RUnlock()
	return servers
}

func GetPremServers(token string) map[string]*server.Server {
	mu.RLock()
	defer mu.RUnlock()
	tempServers := make(map[string]*server.Server)
	for name, v := range servers {
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
	ss[s.Name] = s
	ss[s.Name].EnvLocker = &sync.RWMutex{}
	s.MakeEnv()
	availablePort := s.Port
	for i := 0; i < end; i++ {
		if _, ok := serverIndex[s.Name][i]; ok {
			continue
		}
		// 根据副本数提取子名称
		if serverIndex[s.Name] == nil {
			serverIndex[s.Name] = make(map[int]struct{})
		}
		serverIndex[s.Name][i] = struct{}{}
		env := make(map[string]string)
		for k, v := range s.TempEnv {
			env[k] = v
		}
		subname := subname.NewSubname(s.Name, i)
		var svc *server.Server
		if s.Port > 0 {
			// 顺序拿到可用端口
			availablePort = pkg.GetAvailablePort(availablePort)
			env["PORT"] = strconv.Itoa(availablePort)
			svc = s.Add(availablePort, end, i, subname)
			availablePort++
		} else {
			env["PORT"] = "0"
			svc = s.Add(0, end, i, subname)
		}

		env["NAME"] = subname.String()
		svc.Env = env
		servers[subname.String()] = svc
		if !svc.Disable {
			golog.Infof("start server: %s", svc.SubName)
			servers[subname.String()].Start()
		}

	}
}

// 获取所有服务的状态
func All(role string) []byte {
	mu.TryRLock()
	defer mu.RUnlock()
	statuss := &pkg.StatusList{
		Data:    make([]*status.ServiceStatus, 0),
		Version: global.VERSION,
		Role:    role,
	}
	serviceStatus := make([]*status.ServiceStatus, 0)
	for name := range servers {
		serviceStatus = append(serviceStatus, getStatus(subname.Subname(name).GetName(), name))
	}
	statuss.Code = 200
	statuss.Data = serviceStatus
	send, err := json.MarshalIndent(statuss, "", "\t")
	if err != nil {
		golog.Error(err)
	}
	return send
}

// 获取所有服务的状态
func AllLook(role, token string) []byte {
	mu.TryRLock()
	defer mu.RUnlock()
	statuss := &pkg.StatusList{
		Data:    make([]*status.ServiceStatus, 0),
		Version: global.VERSION,
		Role:    role,
	}
	serviceStatus := make([]*status.ServiceStatus, 0)
	for name, server := range servers {
		if server.Token == token {
			serviceStatus = append(serviceStatus, getStatus(subname.Subname(name).GetName(), name))
		}
	}
	statuss.Code = 200
	statuss.Data = serviceStatus
	send, err := json.MarshalIndent(statuss, "", "\t")
	if err != nil {
		golog.Error(err)
	}
	return send
}

func StopAllServer() {
	mu.RLock()
	defer mu.RUnlock()
	for _, s := range ss {
		err := StopScript(s)
		if err != nil {
			golog.Error(err)
		}
	}
}

func StopPermAllServer(token string) {
	mu.RLock()
	defer mu.RUnlock()
	for _, s := range ss {
		if s.Token == token {
			err := StopScript(s)
			if err != nil {
				golog.Error(err)
			}
		}

	}
}

func GetServerInfo(name string) (*server.Server, bool) {
	mu.RLock()
	defer mu.RUnlock()
	v, ok := servers[name]
	return v, ok
}
