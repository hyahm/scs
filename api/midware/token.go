package midware

import (
	"net/http"
	"strings"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/xmux"
)

func CheckToken(w http.ResponseWriter, r *http.Request) bool {
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
			needToken = false
			break
		}
	}
	if !needToken {
		xmux.GetInstance(r).Set("role", "admin")
		return false
	}
	token := r.Header.Get("Token")
	if token == global.GetToken() {
		// w.Write([]byte(`{"code": 203, "msg": "token error"}`))
		xmux.GetInstance(r).Set("role", "admin")
		return false
	}
	pages := xmux.GetInstance(r).Get(xmux.PAGES).(map[string]struct{})
	golog.Info(name)
	if _, ok := pages["look"]; ok {
		// 如果是查看所有状态， 那么就继续
		if xmux.GetInstance(r).Get(xmux.CURRFUNCNAME) == "AllStatus" {
			xmux.GetInstance(r).Set("token", token)
			xmux.GetInstance(r).Set("role", "look")
			return false
		}
		lookToken := controller.GetLookToken(name)
		if lookToken != "" && token == lookToken {
			xmux.GetInstance(r).Set("role", "look")
			return false
		}
	}
	w.Write([]byte(`{"code": 203, "msg": "token error"}`))
	return true
}
