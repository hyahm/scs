package handle

import (
	"fmt"
	"net/http"
	"scs/pkg/script"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

func Update(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	if _, ok := script.SS.Infos[pname]; ok {
		if _, ok := script.SS.Infos[pname][name]; ok {
			go script.SS.Infos[pname][name].UpdateAndRestart()
		} else {
			w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this script"}`)))
			return
		}

	} else {
		w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting update"}`)))
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting update"}`)))
	return
}

func UpdatePname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	if _, ok := script.SS.Infos[pname]; ok {
		for name := range script.SS.Infos[pname] {
			golog.Info("send update")
			go script.SS.Infos[pname][name].UpdateAndRestart()
		}

	} else {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s}`, pname)))
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting update"}`)))
	return
}

func UpdateAll(w http.ResponseWriter, r *http.Request) {

	for pname, v := range script.SS.Infos {
		for name := range v {
			go script.SS.Infos[pname][name].UpdateAndRestart()
		}
	}
	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting update"}`)))
	return
}
