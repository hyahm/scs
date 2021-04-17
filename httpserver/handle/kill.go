package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs"
	"github.com/hyahm/xmux"
)

func Kill(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	svc, err := scs.GetServerByNameAndSubname(pname, name)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, name)))
		return
	}
	go svc.Kill()
	w.Write([]byte(`{"code": 200, "msg": "already killed"}`))
}

func KillPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	s, err := scs.GetScriptByPname(pname)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, pname)))
		return
	}
	err = s.KillScript()
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, pname)))
		return
	}
	w.Write([]byte(`{"code": 200, "msg": "already killed"}`))
}
