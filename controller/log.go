package controller

import (
	"github.com/hyahm/scs/pkg/config/scripts/subname"
)

// 通过名字来获取token

func GetLookToken(name string) string {
	store.mu.RLock()
	defer store.mu.RUnlock()

	if v, ok := store.ss[subname.Subname(name).GetName()]; ok {
		return v.Token
	}
	return ""
}

func GetPnameToken(pname string) string {
	store.mu.RLock()
	defer store.mu.RUnlock()

	if v, ok := store.ss[pname]; ok {
		return v.Token
	}
	return ""
}
