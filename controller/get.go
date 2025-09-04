package controller

import (
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg/config/scripts"
)

// 获取脚本结构体
func GetServerByNameAndSubname(name string, subname string) (*server.Server, *scripts.Script, bool) {
	script, ok := store.Store.GetScriptByName(name)
	if !ok {
		return nil, nil, false
	}
	svc, ok := store.Store.GetServerByName(subname)
	if !ok {
		return nil, nil, false
	}
	return svc, script, false
}
