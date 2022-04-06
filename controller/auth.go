package controller

import "github.com/hyahm/scs/internal/store"

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

func HavePname(auths []Auth, pname string) (string, bool) {
	for _, auth := range auths {
		if auth.ScriptName == pname {
			return auth.Role, true
		}
	}
	return "", false
}

func HaveName(auths []Auth, name string) (string, bool) {
	svc, ok := store.Store.GetServerByName(name)
	if ok {
		for _, auth := range auths {
			if auth.ScriptName == svc.Name {
				return auth.Role, true
			}
		}
	}

	return "", false
}
