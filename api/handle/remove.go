package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/config/scripts/subname"

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
	svc, ok := controller.GetServerByNameAndSubname(pname, subname.Subname(name))
	if !ok {
		w.Write([]byte(`{"code": 404, "msg": "not found this script"}`))
		return
	}

	go controller.Remove(svc)
	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))
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
	s, ok := controller.GetScriptByPname(pname)
	if !ok {
		w.Write([]byte(`{"code": 404, "msg": "not found this pname"}`))
		return
	}

	// err := config.DeleteScriptToConfigFile(s)
	// if err != nil {
	// 	w.Write([]byte(`{"code": 404, "msg": "not found this script"}`))
	// 	return
	// }

	controller.RemoveScript(s.Name)
	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))
}

// func RemoveAll(w http.ResponseWriter, r *http.Request) {
// 	if reloadKey {
// 		w.Write([]byte(`{"code": 201, "msg": "config file is reloading, waiting completed first"}`))
// 		return
// 	}
// 	reloadKey = true
// 	defer func() {
// 		reloadKey = false
// 	}()
// 	config.DeleteAllScriptToConfigFile()
// 	controller.RemoveAllScripts()

// 	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))
// }
