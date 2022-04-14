package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/xmux"
)

var reloadKey bool

func Disable(w http.ResponseWriter, r *http.Request) {
	if global.CanReload != 0 {
		// 报警相关配置

		xmux.GetInstance(r).Response.(*pkg.Response).Code = 201
		return
	}
	pname := xmux.Var(r)["pname"]

	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	// 上面已经判断过是否存在了， 这里就忽略
	if controller.DisableScript(script, false) {
		err := config.UpdateScriptToConfigFile(script, true)
		if err != nil {
			xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
			xmux.GetInstance(r).Response.(*pkg.Response).Msg = err.Error()
			return
		}
	}

}

func Enable(w http.ResponseWriter, r *http.Request) {
	if global.CanReload != 0 {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 201
		return
	}
	pname := xmux.Var(r)["pname"]
	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	// 上面已经判断过是否存在了， 这里就忽略

	if controller.EnableScript(script) {
		err := config.UpdateScriptToConfigFile(script, true)
		if err != nil {
			xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
			xmux.GetInstance(r).Response.(*pkg.Response).Msg = err.Error()
			return
		}
	}
}
