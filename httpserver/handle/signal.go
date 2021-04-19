package handle

import (
	"net/http"

	"github.com/hyahm/scs"

	"github.com/hyahm/xmux"
)

func CanStop(w http.ResponseWriter, r *http.Request) {

	// golog.Info(string(res))
	name := xmux.Var(r)["name"]
	svc, err := scs.GetServerBySubname(name)
	if err != nil {
		w.Write([]byte(`{"code": 200, "msg": "not found this name"}`))
		return
	}
	svc.Status.CanNotStop = false
	// for pname := range script.SS.Infos {
	// 	if _, ok := script.SS.Infos[pname][name]; ok {
	// 		if script.SS.Infos[pname][name].Status.Status != script.STOP &&
	// 			script.SS.Infos[pname][name].Status.Status != script.INSTALL {
	// 			script.SS.Infos[pname][name].Status.CanNotStop = false
	// 		}
	// 		w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "update successed"}`)))
	// 		return
	// 	}
	// }
	w.Write([]byte(`{"code": 200, "msg": "now can stop"}`))
}

func CanNotStop(w http.ResponseWriter, r *http.Request) {

	// golog.Info(string(res))
	name := xmux.Var(r)["name"]
	svc, err := scs.GetServerBySubname(name)
	if err != nil {
		w.Write([]byte(`{"code": 200, "msg": "not found this name"}`))
		return
	}
	svc.Status.CanNotStop = true

	w.Write([]byte(`{"code": 200, "msg": "now can not stop"}`))
}
