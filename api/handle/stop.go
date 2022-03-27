package handle

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/config/scripts/subname"
	"github.com/hyahm/scs/pkg"

	"github.com/hyahm/xmux"
)

func Stop(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	role := xmux.GetInstance(r).Get("role").(string)
	svc, ok := controller.GetServerByNameAndSubname(pname, subname.Subname(name))

	if !ok {
		w.Write(pkg.NotFoundScript(role))
		return
	}
	go svc.Stop()
	w.Write(pkg.Waiting("stop", role))
}

func StopPname(w http.ResponseWriter, r *http.Request) {
	names := xmux.Var(r)["pname"]
	role := xmux.GetInstance(r).Get("role").(string)
	for _, pname := range strings.Split(names, ",") {
		s, ok := controller.GetScriptByPname(pname)
		if !ok {
			w.Write(pkg.NotFoundScript(role))
			return
		}
		err := controller.StopScript(s)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err.Error())))
			return
		}
	}

	w.Write(pkg.Waiting("stop", role))
}

func StopAll(w http.ResponseWriter, r *http.Request) {
	role := xmux.GetInstance(r).Get("role").(string)
	controller.StopAllServer()
	w.Write(pkg.Waiting("stop", role))
}
