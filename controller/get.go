package controller

import (
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg/config"
)

// 获取脚本结构体
func GetServerByNameAndSubname(name string, subname string) (*server.Server, config.Script, bool) {
	script, ok := store.GetStore().GetScriptByName(name)
	if !ok {
		return &server.Server{}, script, false
	}
	svc, ok := store.GetStore().GetServerByName(subname)
	if !ok {
		return svc, config.Script{}, false
	}
	return svc, script, false
}
