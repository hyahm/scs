package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/xmux"
)

func ServerInfo(w http.ResponseWriter, r *http.Request) {
	name := xmux.Var(r)["name"]
	role := xmux.GetInstance(r).Get("role").(string)
	info, ok := controller.GetServerInfo(name)
	if !ok {
		w.Write(NotFoundScript(role))
		return
	}
	res := Response{
		Data: info,
	}
	w.Write(res.Sucess(""))
}
