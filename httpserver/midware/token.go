package midware

import (
	"net/http"
	"scs/global"
)

func CheckToken(w http.ResponseWriter, r *http.Request) bool {
	if global.Token == "" {
		return false
	}
	if r.Header.Get("Token") != global.Token {
		w.Write([]byte(`{"code": 203, "msg": "token error"}`))
		return true
	}

	return false
}
