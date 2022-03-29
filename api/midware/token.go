package midware

import (
	"net/http"
	"strings"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/xmux"
)

func CheckToken(w http.ResponseWriter, r *http.Request) bool {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]

	var addr string
	if global.ProxyHeader == "" {
		addr = strings.Split(r.RemoteAddr, ":")[0]
	} else {
		addr = r.Header.Get(global.ProxyHeader)
	}

	needToken := true
	for _, v := range global.GetIgnoreToken() {
		if v == addr {
			xmux.GetInstance(r).Set("token", "")
			needToken = false
			break
		}
	}
	if !needToken {
		xmux.GetInstance(r).Set("token", "")
		xmux.GetInstance(r).Set("role", "admin")
		return false
	}
	token := r.Header.Get("Token")
	if token == global.GetToken() {
		// w.Write([]byte(`{"code": 203, "msg": "token error"}`))
		xmux.GetInstance(r).Set("token", "")
		xmux.GetInstance(r).Set("role", "admin")
		return false
	}
	pages := xmux.GetInstance(r).Get(xmux.PAGES).(map[string]struct{})
	if _, ok := pages["scripts"]; ok {
		// 如果是查看所有状态， 那么就继续
		if pname == "" && name == "" {
			xmux.GetInstance(r).Set("token", token)
			xmux.GetInstance(r).Set("role", "scripts")
			return false
		}
		if pname == "" && name != "" {
			lookToken := controller.GetLookToken(name)
			if lookToken != "" && token == lookToken {
				xmux.GetInstance(r).Set("role", "scripts")
				xmux.GetInstance(r).Set("token", "")
				return false
			}
		}

		if pname != "" {
			lookToken := controller.GetPnameToken(pname)
			if lookToken != "" && token == lookToken {
				xmux.GetInstance(r).Set("role", "scripts")
				xmux.GetInstance(r).Set("token", "")
				return false
			}
		}

	}
	w.Write([]byte(`{"code": 203, "msg": "token error"}`))
	return true
}
