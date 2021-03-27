package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/script"

	"github.com/hyahm/xmux"
)

func Remove(w http.ResponseWriter, r *http.Request) {
	if reloadKey {
		w.Write([]byte(`{"code": 201, "msg": "config file is reloading, waiting completed first"}`))
		return
	}
	reloadKey = true
	defer func() {
		reloadKey = false
	}()
	// 读取配置文件
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	if err := script.ReadConfig(); err != nil {
		w.Write([]byte(fmt.Sprint(`{"code": 500, "msg": "config file error"}`)))
		return
	}
	if script.SS.PnameLen(pname) == 1 {

		script.DeleteName(pname)
		if err := script.WriteConfig(); err != nil {
			w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%v"}`, err)))
			return
		}
	}

	if _, ok := script.SS.Infos[pname]; ok {
		if _, ok := script.SS.Infos[pname][name]; ok {
			go script.SS.Infos[pname][name].Remove()
		} else {
			w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this script"}`)))
			return
		}

	} else {
		w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting stop"}`)))
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting stop"}`)))
	return
}

func RemovePname(w http.ResponseWriter, r *http.Request) {
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
	script.DeleteName(pname)
	if err := script.WriteConfig(); err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%v"}`, err)))
		return
	}
	if _, ok := script.SS.Infos[pname]; ok {
		for name := range script.SS.Infos[pname] {
			go script.SS.Infos[pname][name].Remove()
		}

	} else {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s}`, pname)))
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting stop"}`)))
	return
}

func RemoveAll(w http.ResponseWriter, r *http.Request) {
	if reloadKey {
		w.Write([]byte(`{"code": 201, "msg": "config file is reloading, waiting completed first"}`))
		return
	}
	reloadKey = true
	defer func() {
		reloadKey = false
	}()
	if err := script.ReadConfig(); err != nil {
		w.Write([]byte(fmt.Sprint(`{"code": 500, "msg": "config file error"}`)))
		return
	}
	script.DeleteAll()
	if err := script.WriteConfig(); err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%v"}`, err)))
		return
	}
	for pname, v := range script.SS.Infos {
		for name := range v {
			go script.SS.Infos[pname][name].Remove()
		}
	}
	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting stop"}`)))
	return
}
