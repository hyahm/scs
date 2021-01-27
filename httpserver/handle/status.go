package handle

import (
	"net/http"

	"github.com/hyahm/scs/script"

	"github.com/hyahm/xmux"
)

func Status(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]

	w.Write(script.ScriptName(pname, name))
	return
}

func StatusPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]

	w.Write(script.ScriptPname(pname))
	return
}

func AllStatus(w http.ResponseWriter, r *http.Request) {
	w.Write(script.All())
	return

}
