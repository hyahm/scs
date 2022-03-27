package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func ServerInfo(w http.ResponseWriter, r *http.Request) {
	name := xmux.Var(r)["name"]
	role := xmux.GetInstance(r).Get("role").(string)
	info, ok := controller.GetServerInfo(name)
	if !ok {
		w.Write(pkg.NotFoundScript(role))
		return
	}
	res := pkg.Response{
		Data: info,
	}
	w.Write(res.Sucess(""))
}
