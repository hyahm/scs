package handle

import (
	"net/http"
	"strings"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/scripts/subname"

	"github.com/hyahm/xmux"
)

func Restart(w http.ResponseWriter, r *http.Request) {
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
	go controller.RestartServer(svc)
	w.Write(pkg.Waiting("restart", role))
}

func RestartPname(w http.ResponseWriter, r *http.Request) {
	names := xmux.Var(r)["pname"]
	role := xmux.GetInstance(r).Get("role").(string)
	if global.CanReload != 0 {
		w.Write(pkg.WaitingConfigChanged(role))
		return
	}
	for _, pname := range strings.Split(names, ",") {
		s, ok := controller.GetScriptByPname(pname)
		if !ok {
			w.Write(pkg.NotFoundScript(role))
			return
		}
		controller.RestartScript(s)
	}
	w.Write(pkg.Waiting("restart", role))
}

func RestartAll(w http.ResponseWriter, r *http.Request) {
	// 删除所有的脚本
	role := xmux.GetInstance(r).Get("role").(string)
	token := xmux.GetInstance(r).Get("token").(string)
	if global.CanReload != 0 {
		w.Write(pkg.WaitingConfigChanged(role))
		return
	}
	if token != "" {
		controller.RestartPermAllServer(token)
	} else {
		controller.RestartAllServer()
	}

	w.Write(pkg.Waiting("restart", role))
}
