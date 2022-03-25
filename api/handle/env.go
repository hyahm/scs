package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/config/scripts/subname"
	"github.com/hyahm/xmux"
)

func GetEnvName(w http.ResponseWriter, r *http.Request) {
	// 通过pname， name 获取， 因为可能port 不一样
	name := xmux.Var(r)["name"]
	svc, ok := controller.GetServerBySubname(subname.Subname(name).String())
	if !ok {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name %s"}`, name)))
		return
	}
	res := Response{
		Data: svc.Env,
	}
	w.Write(res.Sucess())

}
