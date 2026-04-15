package controller

import (
	"github.com/hyahm/scs/internal/store"
)

// 通过名字来获取token

func GetPnameToken(pname string) string {

	script, ok := store.GetStore().GetScriptByName(pname)
	if !ok {
		return ""
	}
	return script.ScriptToken
}
