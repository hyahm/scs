package handle

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/scripts/subname"

	"github.com/hyahm/xmux"
)

func Stop(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	svc, _, ok := controller.GetServerByNameAndSubname(pname, subname.Subname(name))

	if !ok {
		w.Write(pkg.NotFoundScript())
		return
	}
	go svc.Stop()
	w.Write(pkg.Waiting("stop"))
}

func StopPname(w http.ResponseWriter, r *http.Request) {
	names := xmux.Var(r)["pname"]
	for _, pname := range strings.Split(names, ",") {
		s, ok := controller.GetScriptByPname(pname)
		if !ok {
			w.Write(pkg.NotFoundScript())
			return
		}
		err := controller.StopScript(s)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err.Error())))
			return
		}
	}

	w.Write(pkg.Waiting("stop"))
}

func StopAll(w http.ResponseWriter, r *http.Request) {
	scriptname := xmux.GetInstance(r).Get("scriptname").(map[string]struct{})
	if scriptname == nil {
		controller.StopAllServer()
	} else {
		controller.StopScriptFromName(scriptname)
	}

	w.Write(pkg.Waiting("stop"))
}
