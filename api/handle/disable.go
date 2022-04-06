package handle

import (
	"fmt"
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
		w.Write(pkg.WaitingConfigChanged())
		return
	}
	pname := xmux.Var(r)["pname"]

	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		w.Write(pkg.NotFoundScript())
		return
	}
	// 上面已经判断过是否存在了， 这里就忽略
	if controller.DisableScript(script, false) {
		err := config.UpdateScriptToConfigFile(script, true)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s}`, err.Error())))
			return
		}
	}
	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))

}

func Enable(w http.ResponseWriter, r *http.Request) {
	if global.CanReload != 0 {
		w.Write(pkg.WaitingConfigChanged())
		return
	}
	pname := xmux.Var(r)["pname"]
	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		w.Write(pkg.NotFoundScript())
		return
	}
	// 上面已经判断过是否存在了， 这里就忽略

	if controller.EnableScript(script) {
		err := config.UpdateScriptToConfigFile(script, true)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s}`, err.Error())))
			return
		}
	}

	w.Write([]byte(`{"code": 200, "msg": "waiting start"}`))
}
