package handle

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hyahm/scs/script"

	"github.com/hyahm/golog"
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
	for pname := range script.SS.Infos {
		if _, ok := script.SS.Infos[pname][name]; ok {
			env := make(map[string]string)
			for _, v := range script.SS.Infos[pname][name].GetEnv() {
				golog.Info(v)
				start := strings.Index(v, "=")
				env[v[:start]] = v[start+1:]
			}
			send, _ := json.Marshal(env)
			w.Write(send)
			return
		}
		// if len(script.SS.Infos[pname][name].Env) == 0 {
		// 	w.Write(nil)
		// 	return
		// }
		// send, _ := json.Marshal(script.SS.Infos[pname][name].Env)
		// w.Write(send)
		// return
	}

	w.Write([]byte(`{"code": 404, "msg": "not found this name"}`))
	return
}
