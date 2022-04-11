package module

import (
	"net/http"

	"github.com/hyahm/xmux"
)

func Write(w http.ResponseWriter, r *http.Request, data []byte) {
	xmux.GetInstance(r).Set(xmux.RESPONSEBODY, data)
	w.Write(data)
}
