package module

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/hyahm/xmux"
)

func Write(w http.ResponseWriter, r *http.Request, data []byte) {
	w.Write(data)
}

func Exit(start time.Time, w http.ResponseWriter, r *http.Request) {
	var send []byte
	var err error
	if xmux.GetInstance(r).Response != nil {
		send, err = json.Marshal(xmux.GetInstance(r).Response)
		if err != nil {
			log.Println(err)
		}
		w.Write(send)
	}
	log.Printf("connect_id: %d,method: %s\turl: %s\ttime: %f\t status_code: %v, body: %v\n",
		xmux.GetInstance(r).GetConnectId(),
		r.Method,
		r.URL.Path, time.Since(start).Seconds(), xmux.GetInstance(r).StatusCode,
		string(send))
}
