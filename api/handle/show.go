package handle

import (
	"net/http"

	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func ShowConfig(w http.ResponseWriter, r *http.Request) {
	cfg := internal.GetConfig()
	xmux.GetInstance(r).Response.(*pkg.Response).Data = cfg
}
