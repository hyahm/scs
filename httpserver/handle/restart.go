package handle

import (
	"net/http"
	"scs/script"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

func Restart(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	if v, ok := script.SS.Infos[pname]; ok {
		if _, ok := v[name]; ok {
			if script.SS.Infos[pname][name].Status.Status == script.RUNNING {
				go script.SS.Infos[pname][name].Restart()
			} else if script.SS.Infos[pname][name].Status.Status == script.STOP {
				script.SS.Infos[pname][name].Start()
			}

		} else {
			w.Write([]byte("not found this name:" + name))
			return
		}

	} else {
		w.Write([]byte("not found this pname:" + pname))
		return
	}

	w.Write([]byte("waiting restart"))
	return
}

func RestartPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	if _, ok := script.SS.Infos[pname]; ok {
		for name := range script.SS.Infos[pname] {
			golog.Info("loop ", pname, " ", name)
			if script.SS.Infos[pname][name].Status.Status == script.RUNNING {
				golog.Info("restart ", pname, " ", name)
				go script.SS.Infos[pname][name].Restart()
			} else if script.SS.Infos[pname][name].Status.Status == script.STOP {
				script.SS.Infos[pname][name].Start()
			}
		}

	} else {
		w.Write([]byte("not found this pname:" + pname))
		return
	}

	w.Write([]byte("waiting restart"))
	return
}

func RestartAll(w http.ResponseWriter, r *http.Request) {
	// 删除所有的脚本
	for pname := range script.SS.Infos {
		for name := range script.SS.Infos[pname] {
			if script.SS.Infos[pname][name].Status.Status == script.RUNNING {
				go script.SS.Infos[pname][name].Restart()
			} else if script.SS.Infos[pname][name].Status.Status == script.STOP {
				script.SS.Infos[pname][name].Start()
			}
		}

	}

	w.Write([]byte("waiting restart"))
	return
}
