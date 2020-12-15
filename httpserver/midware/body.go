package midware

import (
	"encoding/json"
	"net/http"

	"github.com/hyahm/xmux"
)

func Unmarshal(w http.ResponseWriter, r *http.Request) bool {

	err := json.NewDecoder(r.Body).Decode(xmux.GetData(r).Data)
	if err != nil {
		w.Write([]byte(err.Error()))
		return true
	}
	return false
}
