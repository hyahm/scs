package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/config"
	"github.com/hyahm/xmux"
)

var reloadKey bool

func Disable(w http.ResponseWriter, r *http.Request) {
	role := xmux.GetInstance(r).Get("role").(string)
	if global.CanReload != 0 {
		w.Write(WaitingConfigChanged(role))
		return
	}
	pname := xmux.Var(r)["pname"]

	s, ok := controller.GetScriptByPname(pname)
	if !ok {
		w.Write(NotFoundScript(role))
		return
	}
	// 上面已经判断过是否存在了， 这里就忽略
	s.Disable = true
	err := config.UpdateScriptToConfigFile(s)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s}`, err.Error())))
		return
	}
	controller.DisableScript(s)
	w.Write([]byte(`{"code": 200, "msg": "waiting stop"}`))

}

func Enable(w http.ResponseWriter, r *http.Request) {
	role := xmux.GetInstance(r).Get("role").(string)
	if global.CanReload != 0 {
		w.Write(WaitingConfigChanged(role))
		return
	}
	pname := xmux.Var(r)["pname"]
	s, ok := controller.GetScriptByPname(pname)
	if !ok {
		w.Write(NotFoundScript(role))
		return
	}
	// 上面已经判断过是否存在了， 这里就忽略
	s.Disable = false
	err := config.UpdateScriptToConfigFile(s)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s}`, err.Error())))
		return
	}
	controller.EnableScript(s)
	w.Write([]byte(`{"code": 200, "msg": "waiting start"}`))
}
