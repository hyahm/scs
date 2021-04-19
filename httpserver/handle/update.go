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
		w.Write([]byte(`{"code": 404, "msg": "not found this script"}`))
		return
	}
	svc.UpdateAndRestart()

	w.Write([]byte(`{"code": 200, "msg": "waiting update"}`))
}

func UpdatePname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	s, err := scs.GetScriptByPname(pname)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s}`, pname)))
		return
	}
	s.UpdateAndRestartScript()

	w.Write([]byte(`{"code": 200, "msg": "waiting update"}`))
}

func UpdateAll(w http.ResponseWriter, r *http.Request) {

	scs.UpdateAndRestartAllServer()
	w.Write([]byte(`{"code": 200, "msg": "waiting update"}`))
}
