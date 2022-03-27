package controller

import (
	"encoding/json"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/config/scripts"
	"github.com/hyahm/scs/internal/config/scripts/subname"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/status"
)

func GetScripts() map[string]*scripts.Script {
	mu.RLock()
	defer mu.RUnlock()
	return ss
}

func KillScript(s *scripts.Script) {
	mu.RLock()
	defer mu.RUnlock()
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(s.Name, i)
		servers[subname.String()].Kill()
	}
}

func NeedStop(s *scripts.Script) bool {
	// 更新server
	// 判断值是否相等
	return !scripts.EqualScript(s, ss[s.Name])
}

func ScriptName(pname, subname, role string) []byte {
	mu.RLock()
	defer mu.RUnlock()
	status := &pkg.StatusList{
		Data:    make([]*status.ServiceStatus, 0),
		Version: global.VERSION,
		Role:    role,
	}
	if _, ok := ss[pname]; !ok {
		golog.Error("not found scripts")
		return nil
	}
	if _, ok := servers[subname]; !ok {
		return nil
	}
	status.Data = append(status.Data, getStatus(pname, subname))
	return status.Marshal()

}

func ScriptPname(pname, role string) []byte {
	mu.RLock()
	defer mu.RUnlock()
	statuss := &pkg.StatusList{
		Data:    make([]*status.ServiceStatus, 0),
		Version: global.VERSION,
		Role:    role,
	}
	if _, ok := ss[pname]; !ok {
		statuss.Msg = "not found " + pname
		send, err := json.MarshalIndent(statuss, "", "\n")

		if err != nil {
			golog.Error(err)
		}
		return send
	}
	replicate := ss[pname].Replicate
	if replicate == 0 {
		replicate = 1
	}
	serviceStatus := make([]*status.ServiceStatus, 0)
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(pname, i).String()
		serviceStatus = append(serviceStatus, getStatus(pname, subname))
	}
	statuss.Data = serviceStatus
	statuss.Code = 200
	send, err := json.MarshalIndent(statuss, "", "\n")

	if err != nil {
		golog.Error(err)
	}
	return send
}
