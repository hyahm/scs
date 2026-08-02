package handle

import (
	"net/http"

	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

type RespRepo struct {
	Url        []string `json:"url"`
	Derivative string   `json:"derivative"`
}

func GetRepo(w http.ResponseWriter, r *http.Request) {
	cfg := internal.GetConfig()
	if cfg.Repo == nil {
		xmux.GetInstance(r).Response.(*pkg.Response).Data = nil
		return
	}
	xmux.GetInstance(r).Response.(*pkg.Response).Data = &RespRepo{
		Url:        cfg.Repo.Url,
		Derivative: cfg.Repo.Derivative,
	}
}
