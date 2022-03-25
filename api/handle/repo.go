package handle

import (
	"net/http"
)

// 获取url和DERIVATIVE

type RespRepo struct {
	Url        []string `json:"url"`
	Derivative string   `json:"derivative"`
}

func GetRepo(w http.ResponseWriter, r *http.Request) {
	// resp := &RespRepo{}
	// resp.Url = server.Cfg.Repo.Url
	// resp.Derivative = server.Cfg.Repo.Derivative
	// send, _ := json.Marshal(resp)
	// w.Write(send)
}
