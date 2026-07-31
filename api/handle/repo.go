package handle

import (
	"net/http"

	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

type RespRepo struct {
	Url        []string `json:"url"`
	Derivative string   `json:"derivative"`
}

func GetRepo(w http.ResponseWriter, r *http.Request) {
	xmux.GetInstance(r).Response.(*pkg.Response).Data = nil
}
