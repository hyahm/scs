package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
	"github.com/hyahm/xmux"
)

func Disable(w http.ResponseWriter, r *http.Request) {

	pname := xmux.Var(r)["pname"]
	golog.Info("disable ", pname)
	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		golog.Warn("1111111111111111")
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	golog.Warn("1111111111111111")
	// 上面已经判断过是否存在了， 这里就忽略
	msg, ok := global.SetReLoading(fmt.Sprintf("enable %s is running", pname))
	if !ok {
		golog.Warn("1111111111111111")
		pkg.Error(r, msg)
		return
	}
	golog.Warn("1111111111111111")
	defer global.SetCanReLoad()
	// 上面已经判断过是否存在了， 这里就忽略
	if controller.DisableScript(script, false) {
		err := config.UpdateScriptToConfigFile(script, true)
		if err != nil {
			xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
			return
		}
	}

}

func Enable(w http.ResponseWriter, r *http.Request) {

	pname := xmux.Var(r)["pname"]
	script, ok := store.Store.GetScriptByName(pname)
	if !ok {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	// 上面已经判断过是否存在了， 这里就忽略
	msg, ok := global.SetReLoading(fmt.Sprintf("enable %s is running", pname))
	if !ok {
		pkg.Error(r, msg)
		return
	}
	defer global.SetCanReLoad()
	if controller.EnableScript(script) {
		err := config.UpdateScriptToConfigFile(script, true)
		if err != nil {
			golog.Error(err)
			xmux.GetInstance(r).Response.(*pkg.Response).Code = 500
			return
		}
		controller.AddScript(script)
	}

}
