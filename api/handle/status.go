package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"

	"github.com/hyahm/xmux"
)

func Status(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	role := xmux.GetInstance(r).Get("role").(string)

	w.Write(controller.ScriptName(pname, name, role))
}

func StatusPname(w http.ResponseWriter, r *http.Request) {
	role := xmux.GetInstance(r).Get("role").(string)

	pname := xmux.Var(r)["pname"]
	w.Write(controller.ScriptPname(pname, role))
}

func AllStatus(w http.ResponseWriter, r *http.Request) {
	role := xmux.GetInstance(r).Get("role").(string)

	if role == "scripts" {
		token := xmux.GetInstance(r).Get("token").(string)
		w.Write(controller.AllLook(role, token))
		return

	}
	w.Write(controller.All(role))

}
