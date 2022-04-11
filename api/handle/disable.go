package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs/api/module"
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
		module.Write(w, r, pkg.WaitingConfigChanged())
		return
	}
	pname := xmux.Var(r)["pname"]

	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		module.Write(w, r, pkg.NotFoundScript())
		return
	}
	// 上面已经判断过是否存在了， 这里就忽略
	if controller.DisableScript(script, false) {
		err := config.UpdateScriptToConfigFile(script, true)
		if err != nil {
			module.Write(w, r, []byte(fmt.Sprintf(`{"code": 500, "msg": "%s}`, err.Error())))
			return
		}
	}
	module.Write(w, r, []byte(`{"code": 200, "msg": "waiting stop"}`))

}

func Enable(w http.ResponseWriter, r *http.Request) {
	if global.CanReload != 0 {
		module.Write(w, r, pkg.WaitingConfigChanged())
		return
	}
	pname := xmux.Var(r)["pname"]
	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		module.Write(w, r, pkg.NotFoundScript())
		return
	}
	// 上面已经判断过是否存在了， 这里就忽略

	if controller.EnableScript(script) {
		err := config.UpdateScriptToConfigFile(script, true)
		if err != nil {
			module.Write(w, r, []byte(fmt.Sprintf(`{"code": 500, "msg": "%s}`, err.Error())))
			return
		}
	}
	module.Write(w, r, []byte(`{"code": 200, "msg": "waiting start"}`))
}
