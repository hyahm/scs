package handle

import (
	"fmt"
	"net/http"
	"scs/pkg/script"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

func Stop(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	if _, ok := script.SS.Infos[pname]; ok {
		if _, ok := script.SS.Infos[pname][name]; ok {
			if script.SS.Infos[pname][name].Status.Status == script.RUNNING {
				go script.SS.Infos[pname][name].Stop()

			} else {
				w.Write([]byte(fmt.Sprintf(`{"code": 201, "msg": "this script not running"}`)))
				return
			}
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
			if script.SS.Infos[pname][name].Status.Status == script.RUNNING {
				golog.Info("send stop")
				go script.SS.Infos[pname][name].Stop()
			}
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
			if script.SS.Infos[pname][name].Status.Status == script.RUNNING {
				go script.SS.Infos[pname][name].Stop()
			}
		}
	}
	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting stop"}`)))
	return
}
