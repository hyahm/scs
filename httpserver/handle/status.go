package handle

import (
	"encoding/json"
	"net/http"
	"scs/pkg/script"
	"time"

	"github.com/hyahm/xmux"
)

func Status(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	name := xmux.Var(r)["name"]
	statuss := make([]*script.ServiceStatus, 0)
	if _, pok := script.SS.Infos[pname]; pok {
		if s, ok := script.SS.Infos[pname][name]; ok {
			if s.Status.Up > 0 {
				s.Status.Start = (time.Duration(time.Now().Unix()-s.Status.Up) * time.Second).String()
			} else {
				s.Status.Start = "0s"
			}
			statuss = append(statuss, s.Status)

		}
	}

	s, _ := json.MarshalIndent(statuss, "", "\n")
	w.Write(s)
	return
}

func StatusPname(w http.ResponseWriter, r *http.Request) {
	pname := xmux.Var(r)["pname"]
	statuss := make([]*script.ServiceStatus, 0)
	for _, s := range script.SS.Infos[pname] {
		if s.Status.Up > 0 {
			s.Status.Start = (time.Duration(time.Now().Unix()-s.Status.Up) * time.Second).String()
		} else {
			s.Status.Start = "0s"
		}
		statuss = append(statuss, s.Status)
	}

	s, err := json.MarshalIndent(statuss, "", "\n")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(s)
	return
}

func AllStatus(w http.ResponseWriter, r *http.Request) {
	w.Write(script.All())
	return

}
