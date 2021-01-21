package handle

import (
	"net/http"
	"strings"

	"github.com/hyahm/scs/pkg/script"

	"github.com/hyahm/xmux"
)

func Log(w http.ResponseWriter, r *http.Request) {
	name := xmux.Var(r)["name"]
	key := "log"
	ns := strings.Split(name, ":")
	if len(ns) > 1 {
		if strings.ToLower(ns[1]) == "lookpath" {
			key = "lookPath"
		} else if strings.ToLower(ns[1]) == "update" {
			key = "update"
		} else if strings.ToLower(ns[1]) == "log" {
			key = "log"
		} else {
			w.Write([]byte("not found this log key: " + ns[1] + ", only support log | update | lookPath"))
			return
		}
	}
	for _, v := range script.SS.Infos {
		for n, script := range v {
			if n == ns[0] {
				script.LogLocker.RLock()
				w.Write([]byte(strings.Join(script.Log[key], "")))
				script.LogLocker.RUnlock()
				return
			}
		}
	}

	w.Write([]byte("not found script"))
	return
}
