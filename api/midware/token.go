package midware

import (
	"net/http"
	"strings"

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
	if r.Header.Get("Token") == global.GetToken() {
		// w.Write([]byte(`{"code": 203, "msg": "token error"}`))
		xmux.GetInstance(r).Set("role", "admin")
		return false
	}
	pages := xmux.GetInstance(r).Get(xmux.PAGES).(map[string]struct{})
	if _, ok := pages["look"]; ok {
		logtoken := controller.GetLogToken(name)
		if logtoken != "" && r.Header.Get("Token") == logtoken {
			xmux.GetInstance(r).Set("role", "look")
			return false
		}
	}

	w.Write([]byte(`{"code": 203, "msg": "token error"}`))
	return true
}
