package handle

import (
	"encoding/json"
	"net/http"
	"scs/config"

	"github.com/hyahm/golog"
)

// 获取url和DERIVATIVE

type RespRepo struct {
	Url        []string `json:"url"`
	Derivative string   `json:"derivative"`
}

func GetRepo(w http.ResponseWriter, r *http.Request) {
	resp := &RespRepo{}
	resp.Url = config.Cfg.Repo.Url
	resp.Derivative = config.Cfg.Repo.Derivative
	send, _ := json.Marshal(resp)
	golog.Info(string(send))
	w.Write(send)
	return
}
