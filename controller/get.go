package controller

import (
	"encoding/json"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/config/scripts"
	"github.com/hyahm/scs/internal/config/scripts/subname"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/status"
)

func GetScriptName(pname string, subname string) []byte {
	mu.RLock()
	mu.RUnlock()
	statuss := &StatusList{
		Data:    make([]*status.ServiceStatus, 0),
		Version: global.VERSION,
	}
	if _, ok := ss[pname]; !ok {
		statuss.Msg = "not found " + pname
		send, err := json.MarshalIndent(statuss, "", "\n")

		if err != nil {
			golog.Error(err)
		}
		return send
	}
	return nil
}

// 获取脚本结构体
func GetServerByNameAndSubname(name string, subname subname.Subname) (*server.Server, bool) {
	mu.RLock()
	defer mu.RUnlock()
	if _, ok := ss[name]; !ok {
		return nil, false
	}
	if _, ok := servers[subname.String()]; ok {
		return servers[subname.String()], true
	}
	return nil, false
}

func GetScriptByPname(name string) (*scripts.Script, bool) {
	mu.RLock()
	defer mu.RUnlock()
	v, ok := ss[name]
	return v, ok

}
