package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs"

	"github.com/hyahm/xmux"
)

func Stop(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]

	svc, err := scs.GetServerByNameAndSubname(pname, name)
	if err != nil {
		w.Write([]byte(`{"code": 404, "msg": "not found this script"}`))
		return
	}
	go svc.Stop()
	// if _, ok := script.SS.Infos[pname]; ok {
	// 	if _, ok := script.SS.Infos[pname][name]; ok {
	// 		go script.SS.Infos[pname][name].Stop()
	// 	} else {
	// 		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this script"}`)))
	// 		return
	// 	}

	// } else {
	// 	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting stop"}`)))
	// 	return
	// }

	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting stop"}`)))
	return
}

func StopPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	s, err := scs.GetScriptByPname(pname)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s}`, pname)))
		return
	}
	err = s.StopScript()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err.Error())))
		return
	}
	// if _, ok := script.SS.Infos[pname]; ok {
	// 	for name := range script.SS.Infos[pname] {
	// 		golog.Info("send stop")
	// 		go script.SS.Infos[pname][name].Stop()
	// 	}

	// } else {
	// 	w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s}`, pname)))
	// 	return
	// }

	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))
}

func StopAll(w http.ResponseWriter, r *http.Request) {

	scs.StartAllServer()
	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))
}
