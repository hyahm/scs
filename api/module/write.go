package module

import (
	"net/http"
)

func Write(w http.ResponseWriter, r *http.Request, data []byte) {
	w.Write(data)
}
