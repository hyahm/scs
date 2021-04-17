package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs"

	"github.com/hyahm/xmux"
)

func Update(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	svc, err := scs.GetServerByNameAndSubname(pname, name)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this script"}`)))
		return
	}
	svc.UpdateAndRestart()

	// if _, ok := script.SS.Infos[pname]; ok {
	// 	if _, ok := script.SS.Infos[pname][name]; ok {
	// 		go script.SS.Infos[pname][name].UpdateAndRestart()
	// 	}

	// } else {
	// w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this script"}`)))
	// return
	// }

	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting update"}`)))
	return
}

func UpdatePname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	s, err := scs.GetScriptByPname(pname)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s}`, pname)))
		return
	}
	s.UpdateAndRestartScript()
	// if _, ok := script.SS.Infos[pname]; ok {
	// 	for name := range script.SS.Infos[pname] {
	// 		golog.Info("send update")
	// 		go script.SS.Infos[pname][name].UpdateAndRestart()
	// 	}

	// } else {
	// 	w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s}`, pname)))
	// 	return
	// }

	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting update"}`)))
	return
}

func UpdateAll(w http.ResponseWriter, r *http.Request) {

	scs.UpdateAndRestartAllServer()
	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting update"}`)))
	return
}
