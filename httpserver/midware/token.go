package midware

import (
	"net/http"
	"scs/global"
	"scs/probe"
	"strings"
)

func CheckToken(w http.ResponseWriter, r *http.Request) bool {
	addr := strings.Split(r.RemoteAddr, ":")[0]
	needToken := true
	for _, v := range probe.VarAT.HWA.Monitored {
		if v == addr {
			needToken = false
			break
		}
	}
	if !needToken {
		return false
	}

	if r.Header.Get("Token") != global.Token {
		w.Write([]byte(`{"code": 203, "msg": "token error"}`))
		return true
	}

	return false
}
