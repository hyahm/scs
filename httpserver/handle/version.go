package handle

import (
	"net/http"

	"github.com/hyahm/scs/script"

	"github.com/hyahm/xmux"
)

func Version(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	version := r.FormValue("version")
	if _, ok := script.SS.Infos[pname]; ok {
		if _, ok := script.SS.Infos[pname][name]; ok {
			if script.SS.Infos[pname][name].Status.Status != script.STOP {
				script.SS.Infos[pname][name].Status.Version = version
			}

		}
	}

	w.Write([]byte(`{"code": 200, "msg": "set verion"}`))
	return
}
