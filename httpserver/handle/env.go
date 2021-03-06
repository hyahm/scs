package handle

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hyahm/scs/server"
	"github.com/hyahm/scs/subname"
	"github.com/hyahm/xmux"
)

func GetEnvName(w http.ResponseWriter, r *http.Request) {
	// 通过pname， name 获取， 因为可能port 不一样
	name := xmux.Var(r)["name"]
	svc, ok := server.GetServerBySubname(subname.Subname(name))
	if !ok {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name %s"}`, name)))
		return
	}

	send, _ := json.Marshal(svc.Env)
	w.Write(send)

}
