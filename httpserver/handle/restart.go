package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs/server"
	"github.com/hyahm/scs/subname"

	"github.com/hyahm/xmux"
)

func Restart(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	svc, ok := server.GetServerByNameAndSubname(pname, subname.Subname(name))
	if !ok {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, name)))
		return
	}
	go svc.Restart()
	w.Write([]byte(`{"code": 200, "msg": "waiting restart"}`))
}

func RestartPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	s, ok := server.GetScriptByPname(pname)
	if !ok {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s"}`, pname)))
		return
	}
	server.RestartScript(s)
	w.Write([]byte(`{"code": 200, "msg": "waiting restart"}`))
}

func RestartAll(w http.ResponseWriter, r *http.Request) {
	// 删除所有的脚本
	server.RestartAllServer()

	w.Write([]byte(`{"code": 200, "msg": "waiting restart"}`))
}
