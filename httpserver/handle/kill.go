package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs/script"

	"github.com/hyahm/xmux"
)

func Kill(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	if _, ok := script.SS.Infos[pname]; ok {
		if _, ok := script.SS.Infos[pname][name]; ok {
			script.SS.Infos[pname][name].Kill()
		}
	} else {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, name)))
		return
	}

	w.Write([]byte(`{"code": 200, "msg": "already killed"}`))
	return
}

func KillPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	if _, ok := script.SS.Infos[pname]; ok {
		for name := range script.SS.Infos[pname] {
			script.SS.Infos[pname][name].Kill()
		}
	} else {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s"}`, pname)))
		return
	}
	w.Write([]byte(`{"code": 200, "msg": "already killed"}`))
	return
}
