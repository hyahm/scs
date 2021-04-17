package handle

import (
	"encoding/json"
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs"
)

// 获取url和DERIVATIVE

type RespRepo struct {
	Url        []string `json:"url"`
	Derivative string   `json:"derivative"`
}

func GetRepo(w http.ResponseWriter, r *http.Request) {
	resp := &RespRepo{}
	resp.Url = scs.Cfg.Repo.Url
	resp.Derivative = scs.Cfg.Repo.Derivative
	send, _ := json.Marshal(resp)
	golog.Info(string(send))
	w.Write(send)
	return
}
