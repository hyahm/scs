package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/script"

	"github.com/hyahm/xmux"
)

func Disable(w http.ResponseWriter, r *http.Request) {
	if reloadKey {
		w.Write([]byte(`{"code": 201, "msg": "config file is reloading, waiting completed first"}`))
		return
	}
	reloadKey = true
	defer func() {
		reloadKey = false
	}()
	pname := xmux.Var(r)["pname"]
	if err := script.ReadConfig(); err != nil {
		w.Write([]byte(fmt.Sprint(`{"code": 500, "msg": "config file error"}`)))
		return
	}
	script.Disable(pname)
	golog.Info(pname)
	if err := script.WriteConfig(); err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%v"}`, err)))
		return
	}
	if _, ok := script.SS.Infos[pname]; ok {
		for name := range script.SS.Infos[pname] {
			script.SS.Infos[pname][name].Disable = true
			script.SS.Infos[pname][name].Status.Disable = true
			go script.SS.Infos[pname][name].Stop()
		}
	} else {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s}`, pname)))
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting stop"}`)))
	return
}

func Enable(w http.ResponseWriter, r *http.Request) {
	if reloadKey {
		w.Write([]byte(`{"code": 201, "msg": "config file is reloading, waiting completed first"}`))
		return
	}
	reloadKey = true
	defer func() {
		reloadKey = false
	}()
	pname := xmux.Var(r)["pname"]
	if err := script.ReadConfig(); err != nil {
		w.Write([]byte(fmt.Sprint(`{"code": 500, "msg": "config file error"}`)))
		return
	}
	golog.Info("have pname")
	script.Enable(pname)
	if err := script.WriteConfig(); err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%v"}`, err)))
		return
	}
	if _, ok := script.SS.Infos[pname]; ok {
		for name := range script.SS.Infos[pname] {
			script.SS.Infos[pname][name].Disable = false
			script.SS.Infos[pname][name].Status.Disable = false
			go script.SS.Infos[pname][name].Start()
		}

	} else {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s}`, pname)))
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting stop"}`)))
	return
}
