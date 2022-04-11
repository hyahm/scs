package controller

import (
	"github.com/hyahm/scs/internal/store"
)

type Auth struct {
	ScriptName string
	Role       string
}

func GetAuthScriptName(token string) []Auth {
	// 获取有权限的scripts
	auths := make([]Auth, 0)
	for pname, script := range store.Store.GetAllScriptMap() {
		if script.Token == token {
			auths = append(auths, Auth{
				ScriptName: pname,
				Role:       script.Role.ToString(),
			})
		}
	}
	return auths
}

func HavePname(auths []Auth, pname, token string) bool {
	script, ok := store.Store.GetScriptByName(pname)
	if ok && script.Token == token {
		for _, auth := range auths {
			if auth.Role == script.Role.ToString() {
				return true
			}

		}
	}
	return false
}

func HaveName(auths []Auth, name, token string) bool {
	svc, ok := store.Store.GetServerByName(name)
	if ok {
		script, ok := store.Store.GetScriptByName(svc.Name)
		if ok && script.Token == token {
			for _, auth := range auths {
				if auth.Role == script.Role.ToString() {
					return true
				}

			}
		}

	}

	return false
}
