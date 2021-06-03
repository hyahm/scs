package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs/server"
	"github.com/hyahm/scs/subname"

	"github.com/hyahm/xmux"
)

func Stop(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]

	svc, ok := server.GetServerByNameAndSubname(pname, subname.Subname(name))
	if !ok {
		w.Write([]byte(`{"code": 404, "msg": "not found this script"}`))
		return
	}
	go svc.Stop()
	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))
}

func StopPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	s, ok := server.GetScriptByPname(pname)
	if !ok {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s}`, pname)))
		return
	}
	err := server.StopScript(s)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err.Error())))
		return
	}
	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))
}

func StopAll(w http.ResponseWriter, r *http.Request) {

	server.StopAllServer()
	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))
}
