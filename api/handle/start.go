package handle

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/scs/internal/config/scripts/subname"

	"github.com/hyahm/xmux"
)

// 只有状态为 stop的 才会启动

func Start(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]

	svc, ok := controller.GetServerByNameAndSubname(pname, subname.Subname(name))
	if !ok {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name: %s"}`, name)))
		return
	}
	svc.Start()

	w.Write([]byte(`{"code": 200, "msg": "already start"}`))

}

func StartPname(w http.ResponseWriter, r *http.Request) {
	names := xmux.Var(r)["pname"]
	for _, pname := range strings.Split(names, ",") {
		_, ok := controller.GetScriptByPname(pname)
		if !ok {
			w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this pname: %s"}`, pname)))
			return
		}
		controller.StartExsitScript(pname)
	}
	w.Write([]byte(`{"code": 200, "msg": "already start"}`))
}

func StartAll(w http.ResponseWriter, r *http.Request) {
	controller.StartAllServer()
	w.Write([]byte(`{"code": 200, "msg": "already start"}`))
}
