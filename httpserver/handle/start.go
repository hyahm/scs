package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs/pkg/script"

	"github.com/hyahm/xmux"
)

// 只有状态为 stop的 才会启动

func Start(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]

	if v, pok := script.SS.Infos[pname]; pok {
		if _, ok := v[name]; ok {
			script.SS.Infos[pname][name].Start()
			w.Write([]byte(`{"code": 200, "msg": "already start"}`))
			return
		} else {
			w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, name)))
			return
		}
	} else {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s"}`, pname)))
		return
	}

}

func StartPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	if _, pok := script.SS.Infos[pname]; pok {
		for name := range script.SS.Infos[pname] {
			script.SS.Infos[pname][name].Start()
		}

	} else {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s"}`, pname)))
		return
	}
	w.Write([]byte(`{"code": 200, "msg": "already start"}`))
	return
}

func StartAll(w http.ResponseWriter, r *http.Request) {

	for pname := range script.SS.Infos {
		for name := range script.SS.Infos[pname] {
			script.SS.Infos[pname][name].Start()
		}
	}
	w.Write([]byte(`{"code": 200, "msg": "already start"}`))
	return
}
