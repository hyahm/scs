package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func ServerInfo(w http.ResponseWriter, r *http.Request) {
	name := xmux.Var(r)["name"]
	info, ok := controller.GetServerInfo(name)
	if !ok {
		w.Write(pkg.NotFoundScript())
		return
	}
	res := pkg.Response{
		Data: info,
	}
	w.Write(res.Sucess(""))
}
