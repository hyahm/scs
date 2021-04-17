package handle

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hyahm/scs/script"

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
			w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this log key: %s, only support log | update | lookPath"}`, ns[1])))
			return
		}
	}

	svc, err := script.GetServerBySubname(name)
	if err != nil {
		w.Write([]byte(`{"code": 404, "msg":"not found script"}`))
		return
	}
	w.Write([]byte(strings.Join(svc.Log[key], "")))

}
