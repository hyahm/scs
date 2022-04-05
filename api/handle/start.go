package handle

import (
	"net/http"
	"strings"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/scripts/subname"

	"github.com/hyahm/xmux"
)

// 只有状态为 stop的 才会启动

func Start(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	svc, _, ok := controller.GetServerByNameAndSubname(pname, subname.Subname(name))
	if !ok {
		w.Write(pkg.NotFoundScript())
		return
	}
	svc.Start()

	w.Write(pkg.Waiting("start"))

}

func StartPname(w http.ResponseWriter, r *http.Request) {
	names := xmux.Var(r)["pname"]
	for _, pname := range strings.Split(names, ",") {
		_, ok := controller.GetScriptByPname(pname)
		if !ok {
			w.Write(pkg.NotFoundScript())
			return
		}
		controller.StartExsitScript(pname)
	}
	w.Write(pkg.Waiting("start"))
}

func StartAll(w http.ResponseWriter, r *http.Request) {
	names := xmux.GetInstance(r).Get("scriptname")
	if names == nil {
		controller.StartAllServer()
	} else {
		controller.StartAllServerFromScript(names.(map[string]struct{}))
	}

	w.Write(pkg.Waiting("start"))
}
