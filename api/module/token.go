package module

import (
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

// 验证script权限及一下
func CheckToken(w http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get("Token")
	// if token == "" {
	// 	xmux.GetInstance(r).Response.(*pkg.Response).Code = 203
	// 	return true
	// }
	golog.Warn("check token ", token)
	if token == global.CS.Token {
		golog.Info("check token success ")
		auths := controller.GetAllAuth()
		xmux.GetInstance(r).Set("validAuths", auths)
		xmux.GetInstance(r).Set("role", "admin")
		return false
	}
	golog.Info("check token ", token)
	roles := xmux.GetInstance(r).GetPageKeys()
	// if _, roles[scripts.SimpleRole.ToString()]
	// 验证所有scripts的权限
	// 接口权限

	// 主要是2种， 一种是 script  一种是 simple
	auths := controller.GetAuthByToken(token)
	validAuths := make([]controller.Auth, 0, len(auths))
	// 根据权限过滤出有用的
	for _, auth := range auths {
		if _, ok := roles[auth.Role]; ok {
			validAuths = append(validAuths, auth)
		}
	}
	if len(validAuths) > 0 {
		// 	// 说明是有这些脚本权限的
		pname := xmux.Var(r)["pname"]
		name := xmux.Var(r)["name"]
		// 	// 如果都是空
		if pname == "" && name == "" {
			// 如果都是空的，那么后面接口基本需要这个来操作
			xmux.GetInstance(r).Set("validAuths", validAuths)
			return false
		}

		if pname != "" && name == "" {
			// 如果只根据pname来操作的话
			for _, auth := range validAuths {
				if pname == auth.ScriptName {
					xmux.GetInstance(r).Set("role", auth.Role)
					return false
				}
			}
		}

		if name != "" && pname == "" {
			for _, auth := range validAuths {
				if name == auth.ServerName {
					xmux.GetInstance(r).Set("role", auth.Role)
					return false
				}
			}
		}

		if pname != "" && name != "" {
			// 只有name的接口
			for _, auth := range validAuths {
				if name == auth.ServerName && pname == auth.ScriptName {
					xmux.GetInstance(r).Set("role", auth.Role)
					return false
				}
			}
		}

	}
	// todo: 后期为了防止密码报错，增开ip黑名单功能
	// var addr string
	// if global.ProxyHeader == "" {
	// 	addr = strings.Split(r.RemoteAddr, ":")[0]
	// } else {
	// 	addr = r.Header.Get(global.ProxyHeader)
	// }
	xmux.GetInstance(r).Response.(*pkg.Response).Code = 203
	return true
}
