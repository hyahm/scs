package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/pkg"

	"github.com/hyahm/xmux"
)

func Status(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	status, err := controller.ScriptName(pname, name)
	if err != nil {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	xmux.GetInstance(r).Response.(*pkg.Response).Data = status
}

func StatusPname(w http.ResponseWriter, r *http.Request) {

	pname := xmux.Var(r)["pname"]
	status, err := controller.ScriptPname(pname)
	if err != nil {
		xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
		return
	}
	xmux.GetInstance(r).Response.(*pkg.Response).Data = status
}

func AllStatus(w http.ResponseWriter, r *http.Request) {
	names := xmux.GetInstance(r).Get("scriptname")
	if names == nil {
		xmux.GetInstance(r).Response.(*pkg.Response).Data = controller.AllStatus()
	} else {
		xmux.GetInstance(r).Response.(*pkg.Response).Data = controller.AllStatusFromScript(names.(map[string]struct{}))
	}
}
