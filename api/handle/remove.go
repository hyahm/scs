package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/scripts/subname"

	"github.com/hyahm/xmux"
)

func Remove(w http.ResponseWriter, r *http.Request) {

	// 读取配置文件
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	role := xmux.GetInstance(r).Get("role").(string)
	if global.CanReload != 0 {
		w.Write(pkg.WaitingConfigChanged(role))
		return
	}
	svc, ok := controller.GetServerByNameAndSubname(pname, subname.Subname(name))
	if !ok {
		w.Write(pkg.NotFoundScript(role))
		return
	}

	go controller.Remove(svc, true)
	w.Write(pkg.Waiting("stop", role))
}

func RemovePname(w http.ResponseWriter, r *http.Request) {
	role := xmux.GetInstance(r).Get("role").(string)
	if global.CanReload != 0 {
		w.Write(pkg.WaitingConfigChanged(role))
		return
	}
	pname := xmux.Var(r)["pname"]
	s, ok := controller.GetScriptByPname(pname)
	if !ok {
		w.Write(pkg.NotFoundScript(role))
		return
	}

	controller.RemoveScript(s.Name)
	w.Write(pkg.Waiting("stop", role))
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
