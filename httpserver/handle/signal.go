package handle

import (
	"fmt"
	"net/http"
	"scs/pkg/script"

	"github.com/hyahm/xmux"
)

func CanStop(w http.ResponseWriter, r *http.Request) {

	// golog.Info(string(res))
	name := xmux.Var(r)["name"]
	for pname := range script.SS.Infos {
		if _, ok := script.SS.Infos[pname][name]; ok {
			if script.SS.Infos[pname][name].Status.Status != script.STOP &&
				script.SS.Infos[pname][name].Status.Status != script.INSTALL {
				script.SS.Infos[pname][name].Status.CanNotStop = false
			}
		}
	}

	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "update successed"}`)))
	return
}

func CanNotStop(w http.ResponseWriter, r *http.Request) {

	// golog.Info(string(res))
	name := xmux.Var(r)["name"]
	for pname := range script.SS.Infos {
		if _, ok := script.SS.Infos[pname][name]; ok {
			if script.SS.Infos[pname][name].Status.Status != script.STOP &&
				script.SS.Infos[pname][name].Status.Status != script.INSTALL {
				script.SS.Infos[pname][name].Status.CanNotStop = true
			}
		}
	}

	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "update successed"}`)))
	return
}
