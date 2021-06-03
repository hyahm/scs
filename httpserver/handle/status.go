package handle

import (
	"net/http"

	"github.com/hyahm/scs/server"
	"github.com/hyahm/scs/subname"

	"github.com/hyahm/xmux"
)

func Status(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]

	w.Write(server.ScriptName(pname, subname.Subname(name)))
}

func StatusPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]

	w.Write(server.ScriptPname(pname))
}

func AllStatus(w http.ResponseWriter, r *http.Request) {
	w.Write(server.All())

}
