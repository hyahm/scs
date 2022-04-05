package controller

import (
	"encoding/json"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/scs/pkg/config/scripts/subname"
)

func GetScriptName(pname string, subname string) []byte {
	store.mu.RLock()
	store.mu.RUnlock()
	statuss := &pkg.StatusList{
		Data:    make([]pkg.ServiceStatus, 0),
		Version: global.VERSION,
	}
	if _, ok := store.ss[pname]; !ok {
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
func GetServerByNameAndSubname(name string, subname subname.Subname) (*server.Server, *scripts.Script, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	if _, ok := store.ss[name]; !ok {
		return nil, nil, false
	}
	if _, ok := store.servers[subname.String()]; ok {
		return store.servers[subname.String()], store.ss[name], true
	}
	return nil, nil, false
}
