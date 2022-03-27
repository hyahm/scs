package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/config/scripts/subname"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func GetEnvName(w http.ResponseWriter, r *http.Request) {
	// 通过pname， name 获取， 因为可能port 不一样
	name := xmux.Var(r)["name"]
	role := xmux.GetInstance(r).Get("role").(string)
	svc, ok := controller.GetServerBySubname(subname.Subname(name).String())
	if !ok {
		w.Write(pkg.NotFoundScript(role))
		return
	}
	res := pkg.Response{
		Data: svc.Env,
		Role: role,
	}
	w.Write(res.Sucess(""))

}
