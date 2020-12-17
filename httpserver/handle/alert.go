package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"scs/alert"
)

func Alert(w http.ResponseWriter, r *http.Request) {

	ra := &alert.RespAlert{}
	err := json.NewDecoder(r.Body).Decode(ra)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code":1, "msg": "%s"}`, err.Error())))
		return
	}
	ra.SendAlert()
	w.WriteHeader(http.StatusOK)
	return
}

func GetAlert(w http.ResponseWriter, r *http.Request) {
	w.Write(alert.GetDispatcher())
	return
}

func Probe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}
