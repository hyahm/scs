package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs/script"

	"github.com/hyahm/xmux"
)

// 只有状态为 stop的 才会启动

func Start(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]

	svc, err := script.GetServerByNameAndSubname(pname, name)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, name)))
		return
	}
	svc.Start()

	w.Write([]byte(`{"code": 200, "msg": "already start"}`))
	return

}

func StartPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	s, err := script.GetScriptByPname(pname)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s"}`, pname)))
		return
	}
	s.StartServer()
	// if _, pok := script.SS.Infos[pname]; pok {
	// 	for name := range script.SS.Infos[pname] {
	// 		script.SS.Infos[pname][name].Start()
	// 	}

	// } else {
	// 	w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s"}`, pname)))
	// 	return
	// }
	w.Write([]byte(`{"code": 200, "msg": "already start"}`))
	return
}

func StartAll(w http.ResponseWriter, r *http.Request) {

	script.StartAllServer()
	w.Write([]byte(`{"code": 200, "msg": "already start"}`))
	return
}
