package handle

import (
	"net/http"
	"strings"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/config/scripts/subname"

	"github.com/hyahm/xmux"
)

// 只有状态为 stop的 才会启动

func Start(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	role := xmux.GetInstance(r).Get("role").(string)
	svc, ok := controller.GetServerByNameAndSubname(pname, subname.Subname(name))
	if !ok {
		w.Write(NotFoundScript(role))
		return
	}
	svc.Start()

	w.Write(Waiting("start", role))

}

func StartPname(w http.ResponseWriter, r *http.Request) {
	names := xmux.Var(r)["pname"]
	role := xmux.GetInstance(r).Get("role").(string)
	for _, pname := range strings.Split(names, ",") {
		_, ok := controller.GetScriptByPname(pname)
		if !ok {
			w.Write(NotFoundScript(role))
			return
		}
		controller.StartExsitScript(pname)
	}
	w.Write(Waiting("start", role))
}

func StartAll(w http.ResponseWriter, r *http.Request) {
	role := xmux.GetInstance(r).Get("role").(string)
	controller.StartAllServer()
	w.Write(Waiting("start", role))
}
