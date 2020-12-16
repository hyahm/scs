package handle

import (
	"encoding/json"
	"fmt"
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
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%v"}`, err)))
		return
	}

	if v, ok := script.SS.Infos[sg.Pname]; ok {
		if _, ok := v[sg.Name]; ok {
			if script.SS.Infos[sg.Pname][sg.Name].Status.Status != script.STOP {
				script.SS.Infos[sg.Pname][sg.Name].Status.CanNotStop = sg.Value
			}
		}
	}
	w.Write([]byte(fmt.Sprintf(`{"code": 200, "msg": "update successed"}`)))
	return
}
