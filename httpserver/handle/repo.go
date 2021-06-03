package handle

import (
	"encoding/json"
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/server"
)

// 获取url和DERIVATIVE

type RespRepo struct {
	Url        []string `json:"url"`
	Derivative string   `json:"derivative"`
}

func GetRepo(w http.ResponseWriter, r *http.Request) {
	resp := &RespRepo{}
	resp.Url = server.Cfg.Repo.Url
	resp.Derivative = server.Cfg.Repo.Derivative
	send, _ := json.Marshal(resp)
	golog.Info(string(send))
	w.Write(send)
}
