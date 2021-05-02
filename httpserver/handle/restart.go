package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs"

	"github.com/hyahm/xmux"
)

func Restart(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	svc, err := scs.GetServerByNameAndSubname(pname, scs.Subname(name))
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, name)))
		return
	}
	go svc.Restart()
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
	w.Write([]byte(`{"code": 200, "msg": "waiting restart"}`))
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
}
