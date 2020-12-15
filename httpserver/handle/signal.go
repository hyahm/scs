package handle

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"scs/script"

	"github.com/hyahm/golog"
)

type signal struct {
	Pname string `json:"pname"`
	Name  string `json:"name"`
	Value bool   `json:"value"`
}

func Signal(w http.ResponseWriter, r *http.Request) {
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		golog.Error(err)
		w.Write([]byte(err.Error()))
		return
	}
	sg := &signal{}
	// golog.Info(string(res))
	err = json.Unmarshal(res, sg)
	if err != nil {
		golog.Error(err)
		w.Write([]byte(err.Error()))
		return
	}

	if v, ok := script.SS.Infos[sg.Pname]; ok {
		if _, ok := v[sg.Name]; ok {
			if script.SS.Infos[sg.Pname][sg.Name].Status.Status != script.STOP {
				script.SS.Infos[sg.Pname][sg.Name].Status.CanNotStop = sg.Value
			}
		}
	}

	w.Write([]byte("update successed"))
	return
}
