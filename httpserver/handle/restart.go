package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs"

	"github.com/hyahm/xmux"
)

func Restart(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	golog.Info(1111)
	svc, err := scs.GetServerByNameAndSubname(pname, name)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, name)))
		return
	}
	golog.Info(1111)
	go svc.Restart()
	// if v, ok := script.SS.Infos[pname]; ok {
	// 	if _, ok := v[name]; ok {
	// 		go script.SS.Infos[pname][name].Restart()

	// 	} else {
	// 		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, name)))
	// 		return
	// 	}

	// } else {
	// 	w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s"}`, pname)))
	// 	return
	// }
	w.Write([]byte(`{"code": 200, "msg": "waiting restart"}`))
}

func RestartPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]

	s, err := scs.GetScriptByPname(pname)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s"}`, pname)))
		return
	}
	s.RestartScript()
	// if _, ok := script.SS.Infos[pname]; ok {
	// 	for name := range script.SS.Infos[pname] {
	// 		go script.SS.Infos[pname][name].Restart()
	// 	}

	// } else {
	// 	w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s"}`, pname)))
	// 	return
	// }

	w.Write([]byte(`{"code": 200, "msg": "waiting restart"}`))
	return
}

func RestartAll(w http.ResponseWriter, r *http.Request) {
	// 删除所有的脚本

	scs.RestartAllServer()
	// for pname := range script.SS.Infos {
	// 	for name := range script.SS.Infos[pname] {
	// 		go script.SS.Infos[pname][name].Restart()
	// 	}

	// }
	w.Write([]byte(`{"code": 200, "msg": "waiting restart"}`))
	return
}
