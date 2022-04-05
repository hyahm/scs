package controller

type Auth struct {
	ScriptName string
	Role       string
}

func GetAuthScriptName(token string) []Auth {
	// 获取有权限的scripts
	store.mu.RLock()
	defer store.mu.RUnlock()
	auths := make([]Auth, 0)
	for pname, script := range store.ss {
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
	store.mu.RLock()
	defer store.mu.RUnlock()
	if v, ok := store.servers[name]; ok {
		for _, auth := range auths {
			if auth.ScriptName == v.Name {
				return auth.Role, true
			}
		}
	}

	return "", false
}
