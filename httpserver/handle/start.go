package handle

import (
	"net/http"
	"scs/script"

	"github.com/hyahm/xmux"
)

// 只有状态为 stop的 才会启动

func Start(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]

	if v, pok := script.SS.Infos[pname]; pok {
		if _, ok := v[name]; ok {
			if script.SS.Infos[pname][name].Status.Status == script.RUNNING {
				w.Write([]byte("is running"))
				return
			}
			script.SS.Infos[pname][name].Status.Status = script.RUNNING
			script.SS.Infos[pname][name].Start()
			w.Write([]byte("already start"))
			return
		} else {
			w.Write([]byte("not found this name:" + name))
			return
		}
	} else {
		w.Write([]byte("not found this pname:" + pname))
		return
	}

}

func StartPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	if _, pok := script.SS.Infos[pname]; pok {
		for name := range script.SS.Infos[pname] {
			if script.SS.Infos[pname][name].Status.Status == script.STOP {
				script.SS.Infos[pname][name].Start()
			}

		}

	} else {
		w.Write([]byte("not found this pname:" + pname))
		return
	}
	w.Write([]byte("waiting start"))
	return
}

func StartAll(w http.ResponseWriter, r *http.Request) {

	for pname := range script.SS.Infos {
		for name := range script.SS.Infos[pname] {
			if script.SS.Infos[pname][name].Status.Status == script.STOP {
				script.SS.Infos[pname][name].Start()
			}
		}
	}

	w.Write([]byte("already start"))
	return
}
