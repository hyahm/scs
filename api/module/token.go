package module

import (
	"net/http"
	"strings"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func CheckAdminToken(w http.ResponseWriter, r *http.Request) bool {
	var addr string
	if global.ProxyHeader == "" {
		addr = strings.Split(r.RemoteAddr, ":")[0]
	} else {
		addr = r.Header.Get(global.ProxyHeader)
	}

	// 先拿到这个地址需要的权限
	// 以最低的权限为当前的权限

	for _, v := range global.GetIgnoreToken() {
		// 免token的admin权限
		if v == addr {
			xmux.GetInstance(r).Set("token", "")
			xmux.GetInstance(r).Set("role", "admin")
			break
		}
	}
	token := r.Header.Get("Token")

	if token == global.GetToken() {
		// Write(w,r,([]byte(`{"code": 203, "msg": "token error"}`))
		xmux.GetInstance(r).Set("token", token)
		xmux.GetInstance(r).Set("role", "admin")
		return false
	}
	xmux.GetInstance(r).Set(xmux.STATUSCODE, 203)
	Write(w, r, []byte(`{"code": 203, "msg": "token error or no permission"}`))
	return true
}

func CheckAllScriptToken(w http.ResponseWriter, r *http.Request) bool {
	// 验证所有scripts的权限
	token := r.Header.Get("Token")
	// 接口权限
	roles := xmux.GetInstance(r).Get(xmux.PAGES).(map[string]struct{})
	// 主要是2种， 一种是 script  一种是 simple

	auths := controller.GetAuthScriptName(token)
	if len(auths) > 0 {
		// 说明是有这些脚本权限的
		pname := xmux.Var(r)["pname"]
		name := xmux.Var(r)["name"]
		// 如果都是空
		if pname == "" && name == "" {
			scriptname := make(map[string]struct{})
			for _, auth := range auths {
				if _, ok := roles[auth.Role]; ok {
					//
					scriptname[auth.ScriptName] = struct{}{}
				}
			}
			if len(scriptname) > 0 {
				// 全操作， 获取 接口权限
				xmux.GetInstance(r).Set("scriptname", scriptname)
				return false
			}

		}

		if pname != "" && name == "" {
			// 如果只根据pname来操作的话
			ok := controller.HavePname(auths, pname, token)
			if ok {
				return false
			}

		}

		if pname != "" && name != "" {
			_, ok := store.Store.GetScriptByName(pname)
			_, sok := store.Store.GetServerByName(name)
			if !ok || !sok {
				xmux.GetInstance(r).Set(xmux.STATUSCODE, 404)
				Write(w, r, ([]byte(`{"code": 404, "msg": "pname and name not match"}`)))
				return true
			}
			// 如果只根据pname来操作的话, 2种都有
			ok = controller.HavePname(auths, pname, token)
			if ok {
				return false
			}
		}

		if pname == "" && name != "" {
			// 只有name的接口
			svc, ok := store.Store.GetServerByName(name)
			if !ok {
				xmux.GetInstance(r).Set(xmux.STATUSCODE, 404)
				Write(w, r, ([]byte(`{"code": 404, "msg": "pname and name not match"}`)))
				return true
			}
			ok = controller.HavePname(auths, svc.Name, token)
			if ok {
				return false
			}
		}

	}
	xmux.GetInstance(r).Set(xmux.STATUSCODE, 203)
	// Write(w, r, []byte(`{"code": 203, "msg": "token error or no permission"}`))
	return true
}

func CheckToken(w http.ResponseWriter, r *http.Request) bool {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]

	var addr string
	if global.ProxyHeader == "" {
		addr = strings.Split(r.RemoteAddr, ":")[0]
	} else {
		addr = r.Header.Get(global.ProxyHeader)
	}

	// 先拿到这个地址需要的权限
	// 以最低的权限为当前的权限

	for _, v := range global.GetIgnoreToken() {
		// 免token的admin权限
		if v == addr {
			xmux.GetInstance(r).Set("token", "")
			xmux.GetInstance(r).Set("role", "admin")
			break
		}
	}
	token := r.Header.Get("Token")

	if token == global.GetToken() {
		// Write(w,r,([]byte(`{"code": 203, "msg": "token error"}`))
		xmux.GetInstance(r).Set("token", token)
		xmux.GetInstance(r).Set("role", "admin")
		return false
	}

	// 验证 scripts的权限
	pages := xmux.GetInstance(r).Get(xmux.PAGES).(map[string]struct{})
	if _, ok := pages["scripts"]; ok {
		// 如果是查看所有状态， 那么就继续
		if pname == "" && name == "" {
			xmux.GetInstance(r).Set("token", token)
			xmux.GetInstance(r).Set("role", "scripts")
			return false
		}
		if pname == "" && name != "" {
			svc, ok := store.Store.GetServerByName(name)
			if !ok {
				xmux.GetInstance(r).Set(xmux.STATUSCODE, 404)
				Write(w, r, (pkg.NotFoundScript()))
				return true
			}
			if svc.Token != "" && token == svc.Token {
				xmux.GetInstance(r).Set("role", "scripts")
				xmux.GetInstance(r).Set("token", "")
				return false
			}
		}

		if pname != "" {
			script, ok := store.Store.GetScriptByName(pname)
			if !ok {
				xmux.GetInstance(r).Set(xmux.STATUSCODE, 404)
				Write(w, r, (pkg.NotFoundScript()))
				return true
			}
			if script.Token != "" && token == script.Token {
				xmux.GetInstance(r).Set("role", "scripts")
				xmux.GetInstance(r).Set("token", "")
				return false
			}
		}

	}
	Write(w, r, ([]byte(`{"code": 203, "msg": "token error"}`)))
	return true
}
