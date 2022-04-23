package handle

import (
	"net/http"
	"sync/atomic"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"

	"github.com/hyahm/xmux"
)

func Remove(w http.ResponseWriter, r *http.Request) {

	// 读取配置文件
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	if global.CanReload != 0 {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 201
		return
	}
	_, ok := store.Store.GetScriptByName(pname)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	svc, ok := store.Store.GetServerByName(name)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	atomic.AddInt64(&global.CanReload, 1)
	go controller.Remove(svc, true)
}

func RemovePname(w http.ResponseWriter, r *http.Request) {
	if global.CanReload != 0 {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 201
		return
	}
	pname := xmux.Var(r)["pname"]
	_, ok := store.Store.GetScriptByName(pname)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404

		return
	}

	controller.RemoveScript(pname)
}

// func RemoveAll(w http.ResponseWriter, r *http.Request) {
// 	if reloadKey {
// 		Write(w, r,[]byte(`{"code": 201, "msg": "config file is reloading, waiting completed first"}`))
// 		return
// 	}
// 	reloadKey = true
// 	defer func() {
// 		reloadKey = false
// 	}()
// 	config.DeleteAllScriptToConfigFile()
// 	controller.RemoveAllScripts()

// 	Write(w, r,[]byte(`{"code": 200, "msg": "waiting stop"}`))
// }
