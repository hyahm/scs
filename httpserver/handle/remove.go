package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs"

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
	svc, err := scs.GetServerByNameAndSubname(pname, name)
	if err != nil {
		w.Write([]byte(`{"code": 500, "msg": "config file error"}`))
		return
	}

	s, err := scs.GetScriptByPname(pname)
	if err != nil {
		w.Write([]byte(`{"code": 404, "msg": "not found this script"}`))
		return
	}
	golog.Info(s.Replicate)
	s.Replicate -= 1
	err = scs.UpdateScriptToConfigFile(s)
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

	s, err := scs.GetScriptByPname(pname)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%v"}`, err)))
		return
	}
	err = scs.DeleteScriptToConfigFile(s)
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
	scs.DeleteAllScriptToConfigFile()
	scs.RemoveAllScripts()

	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "waiting stop"}`)))
	return
}
