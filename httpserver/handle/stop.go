package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs/pkg/script"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

func Stop(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	if _, ok := script.SS.Infos[pname]; ok {
		if _, ok := script.SS.Infos[pname][name]; ok {
			go script.SS.Infos[pname][name].Stop()
		} else {
			w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this script"}`)))
			return
		}

	} else {
		w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting stop"}`)))
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting stop"}`)))
	return
}

func StopPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	if _, ok := script.SS.Infos[pname]; ok {
		for name := range script.SS.Infos[pname] {
			golog.Info("send stop")
			go script.SS.Infos[pname][name].Stop()
		}

	} else {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s}`, pname)))
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting stop"}`)))
	return
}

func StopAll(w http.ResponseWriter, r *http.Request) {

	for pname, v := range script.SS.Infos {
		for name := range v {
			go script.SS.Infos[pname][name].Stop()
		}
	}
	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting stop"}`)))
	return
}
