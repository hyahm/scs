package midware

import (
	"fmt"
	"net/http"
	"scs/global"
	"strings"
)

func CheckToken(w http.ResponseWriter, r *http.Request) bool {
	addr := strings.Split(r.RemoteAddr, ":")[0]
	needToken := true
	for _, v := range global.IgnoreToken {
		if v == addr {
			needToken = false
			break
		}
	}
	if !needToken {
		return false
	}
	fmt.Println("header token:", r.Header.Get("Token"))
	fmt.Println("global token:", global.Token)
	if r.Header.Get("Token") != global.Token {
		w.Write([]byte(`{"code": 203, "msg": "token error"}`))
		return true
	}

	return false
}
