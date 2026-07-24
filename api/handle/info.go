package handle

import (
	"net/http"
)

func ServerInfo(w http.ResponseWriter, r *http.Request) {
	// name := xmux.Var(r)["name"]
	// svc, ok := store.GetStore().GetServerByName(name)
	// if !ok {
	// 	xmux.GetInstance(r).Response.(*pkg.Response).Code = 404
	// 	return
	// }
	// xmux.GetInstance(r).Response.(*pkg.Response).Data = svc
}
