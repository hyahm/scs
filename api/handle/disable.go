package handle

import (
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/xmux"
)

func Disable(w http.ResponseWriter, r *http.Request) {
	if global.CanReload != 0 {
		// 报警相关配置
		xmux.GetInstance(r).Set(xmux.STATUSCODE, 201)
		return
	}
	pname := xmux.Var(r)["pname"]

	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		xmux.GetInstance(r).Set(xmux.STATUSCODE, 404)
		return
	}
	// 上面已经判断过是否存在了， 这里就忽略
	if controller.DisableScript(script, false) {
		err := config.UpdateScriptToConfigFile(script, true)
		if err != nil {
			xmux.GetInstance(r).Set(xmux.STATUSCODE, 500)
			return
		}
	}

}

func Enable(w http.ResponseWriter, r *http.Request) {
	if global.CanReload != 0 {
		xmux.GetInstance(r).Set(xmux.STATUSCODE, 201)
		return
	}
	pname := xmux.Var(r)["pname"]
	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		xmux.GetInstance(r).Set(xmux.STATUSCODE, 404)
		return
	}
	// 上面已经判断过是否存在了， 这里就忽略

	if controller.EnableScript(script) {
		err := config.UpdateScriptToConfigFile(script, true)
		if err != nil {
			golog.Error(err)
			xmux.GetInstance(r).Set(xmux.STATUSCODE, 500)
			return
		}
	}
}
