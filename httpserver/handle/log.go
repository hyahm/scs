package handle

import (
	"net/http"
	"scs/pkg/script"
	"strings"

	"github.com/hyahm/xmux"
)

func Log(w http.ResponseWriter, r *http.Request) {
	name := xmux.Var(r)["name"]
	for _, v := range script.SS.Infos {
		for n, script := range v {
			if n == name {
				w.Write([]byte(strings.Join(script.Log, "")))
				return
			}
		}
	}

	w.Write([]byte("not found script"))
	return
}
