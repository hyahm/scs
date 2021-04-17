package handle

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hyahm/scs"
	"github.com/hyahm/xmux"
)

// func GetEnv(w http.ResponseWriter, r *http.Request) {
// 	// 通过pname， name 获取， 因为可能port 不一样
// 	pname := xmux.Var(r)["pname"]
// 	name := xmux.Var(r)["name"]
// 	for _, v := range script.SS.Infos[pname] {

// 	}
// 	if _, ok := script.SS.Infos[pname]; ok {
// 		if _, ok := script.SS.Infos[pname][name]; ok {
// 			if len(script.SS.Infos[pname][name].Env) == 0 {
// 				w.Write(nil)
// 				return
// 			}
// 			send, _ := json.Marshal(script.SS.Infos[pname][name].Env)
// 			w.Write(send)
// 			return
// 		}
// 	}

// 	w.Write([]byte("[]"))
// 	return
// }

func GetEnvName(w http.ResponseWriter, r *http.Request) {
	// 通过pname， name 获取， 因为可能port 不一样
	name := xmux.Var(r)["name"]
	svc, err := scs.GetServerBySubname(name)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 404, "msg": "not found this name %s"}`, name)))
		return
	}
	// env := make(map[string]string)
	// for k, v := range svc.Env {
	// 	golog.Info(v)
	// 	start := strings.Index(v, "=")
	// 	env[k] = v[start+1:]
	// }
	send, _ := json.Marshal(svc.Env)
	w.Write(send)

}
