package handle

import (
	"net/http"
)

func GetEnvName(w http.ResponseWriter, r *http.Request) {
	// 通过pname， name 获取， 因为可能port 不一样
	// name := xmux.Var(r)["name"]
	// svc, ok := store.GetStore().GetServerByName(name)
	// if !ok {
	// 	xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
	// 	return
	// }
	// xmux.GetInstance(r).Response.(*pkg.Response).Data = svc.Env
}
