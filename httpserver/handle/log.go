package handle

import (
	"net/http"
	"strings"

	"github.com/hyahm/scs/server"
	"github.com/hyahm/scs/subname"
	"github.com/hyahm/xmux"
)

func Log(w http.ResponseWriter, r *http.Request) {
	name := xmux.Var(r)["name"]

	svc, ok := server.GetServerBySubname(subname.Subname(name))
	if !ok {
		w.Write([]byte(`{"code": 404, "msg":"not found script"}`))
		return
	}
	w.Write([]byte(strings.Join(svc.Log, "")))

}
