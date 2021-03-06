package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs/server"
	"github.com/hyahm/scs/subname"
	"github.com/hyahm/xmux"
)

func Kill(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	svc, ok := server.GetServerByNameAndSubname(pname, subname.Subname(name))
	if !ok {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, name)))
		return
	}
	go svc.Kill()
	w.Write([]byte(`{"code": 200, "msg": "already killed"}`))
}

func KillPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	s, ok := server.GetScriptByPname(pname)
	if !ok {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, pname)))
		return
	}
	server.KillScript(s)

	w.Write([]byte(`{"code": 200, "msg": "already killed"}`))
}
