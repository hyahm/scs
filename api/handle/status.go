package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"

	"github.com/hyahm/xmux"
)

func Status(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	w.Write(controller.ScriptName(pname, name))
}

func StatusPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	w.Write(controller.ScriptPname(pname))
}

func AllStatus(w http.ResponseWriter, r *http.Request) {
	w.Write(controller.All())

}
