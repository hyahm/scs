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
	svc, err := script.GetServerByNameAndSubname(pname, name)
	if err != nil {
		w.Write([]byte(`{"code": 500, "msg": "config file error"}`))
		return
	}

	s, err := script.GetScriptByPname(pname)
	if err != nil {
		w.Write([]byte(`{"code": 404, "msg": "not found this script"}`))
		return
	}
	golog.Info(s.Replicate)
	s.Replicate -= 1
	err = script.UpdateScriptToConfigFile(s)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%v"}`, err)))
		return
	}
	svc.Remove()
	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))
}

func RemovePname(w http.ResponseWriter, r *http.Request) {
	golog.Info("1111")
	if reloadKey {
		w.Write([]byte(`{"code": 201, "msg": "config file is reloading, waiting completed first"}`))
		return
	}
	reloadKey = true
	defer func() {
		reloadKey = false
	}()
	pname := xmux.Var(r)["pname"]

	s, err := script.GetScriptByPname(pname)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%v"}`, err)))
		return
	}
	err = script.DeleteScriptToConfigFile(s)
	if err != nil {
		w.Write([]byte(`{"code": 404, "msg": "not found this script"}`))
		return
	}

	s.RemoveScript()

	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))
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
	script.DeleteAllScriptToConfigFile()
	script.RemoveAllScripts()

	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting stop"}`)))
	return
}
