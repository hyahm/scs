package midware

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/hyahm/xmux"
)

func Unmarshal(w http.ResponseWriter, r *http.Request) bool {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return true
	}
	err = json.Unmarshal(b, xmux.GetInstance(r).Data)
	if err != nil {
		w.Write([]byte(err.Error()))
		return true
	}
	return false
}
