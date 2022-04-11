package handle

import (
	"net/http"

	"github.com/hyahm/scs/api/module"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func GetEnvName(w http.ResponseWriter, r *http.Request) {
	// 通过pname， name 获取， 因为可能port 不一样
	name := xmux.Var(r)["name"]
	svc, ok := store.Store.GetServerByName(name)
	if !ok {
		module.Write(w, r, pkg.NotFoundScript())
		return
	}
	res := pkg.Response{
		Data: svc.Env,
	}
	module.Write(w, r, res.Sucess(""))
}
