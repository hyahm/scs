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

	svc, err := scs.GetServerByNameAndSubname(pname, scs.Subname(name))
	if err != nil {
		w.Write([]byte(`{"code": 404, "msg": "not found this script"}`))
		return
	}
	go svc.Stop()
	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))
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
	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))
}

func StopAll(w http.ResponseWriter, r *http.Request) {

	scs.StopAllServer()
	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))
}
