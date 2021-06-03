package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs/server"
	"github.com/hyahm/scs/subname"

	"github.com/hyahm/xmux"
)

func Update(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	svc, err := server.GetServerByNameAndSubname(pname, subname.Subname(name))
	if err != nil {
		w.Write([]byte(`{"code": 404, "msg": "not found this script"}`))
		return
	}
	svc.UpdateAndRestart()

	w.Write([]byte(`{"code": 200, "msg": "waiting update"}`))
}

func UpdatePname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	s, err := server.GetScriptByPname(pname)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s}`, pname)))
		return
	}
	server.UpdateAndRestartScript(s)

	w.Write([]byte(`{"code": 200, "msg": "waiting update"}`))
}

func UpdateAll(w http.ResponseWriter, r *http.Request) {

	server.UpdateAndRestartAllServer()
	w.Write([]byte(`{"code": 200, "msg": "waiting update"}`))
}
