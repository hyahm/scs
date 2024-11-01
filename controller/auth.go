package controller

import (
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg/config/scripts"
)

type Auth struct {
	ServerName string
	ScriptName string
	Role       string
}

// 获取有权限的server
func GetAuthByToken(token string) []Auth {

	auths := make([]Auth, 0)
	for name, server := range store.Store.GetAllServerMap() {
		if server.ScriptToken == token {
			auths = append(auths, Auth{
				ServerName: name,
				ScriptName: server.Name,
				Role:       string(scripts.ScriptRole),
			})
		}
		if server.SimpleToken == token {
			auths = append(auths, Auth{
				ServerName: name,
				ScriptName: server.Name,
				Role:       string(scripts.SimpleRole),
			})
		}
	}
	return auths
}

func GetAllAuth() []Auth {
	auths := make([]Auth, 0)
	for name, server := range store.Store.GetAllServerMap() {
		auths = append(auths, Auth{
			ServerName: name,
			ScriptName: server.Name,
			Role:       string(scripts.ScriptRole),
		})
	}
	return auths
}

// func HavePname1(auths []Auth, pname, token string) bool {
// 	script, ok := store.Store.GetScriptByName(pname)
// 	if ok {
// 		if script.ScriptToken == token {
// 			for _, auth := range auths {
// 				if auth.Role == script.Role.ToString() {
// 					return true
// 				}

// 			}
// 		}
// 	}
// 	return false
// }

// func HaveName1(auths []Auth, name, token string) bool {
// 	svc, ok := store.Store.GetServerByName(name)
// 	if ok {
// 		script, ok := store.Store.GetScriptByName(svc.Name)
// 		if ok && script.Token == token {
// 			for _, auth := range auths {
// 				if auth.Role == script.Role.ToString() {
// 					return true
// 				}

// 			}
// 		}

// 	}

// 	return false
// }
