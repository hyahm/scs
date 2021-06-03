package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs/server"
	"github.com/hyahm/scs/subname"

	"github.com/hyahm/xmux"
)

// 只有状态为 stop的 才会启动

func Start(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]

	svc, err := server.GetServerByNameAndSubname(pname, subname.Subname(name))
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, name)))
		return
	}
	svc.Start()

	w.Write([]byte(`{"code": 200, "msg": "already start"}`))

}

func StartPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	s, err := server.GetScriptByPname(pname)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s"}`, pname)))
		return
	}
	server.StartServer(s)
	w.Write([]byte(`{"code": 200, "msg": "already start"}`))
}

func StartAll(w http.ResponseWriter, r *http.Request) {
	server.StartAllServer()
	w.Write([]byte(`{"code": 200, "msg": "already start"}`))
}
