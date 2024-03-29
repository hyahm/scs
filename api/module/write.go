package module

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/xmux"
)

func Write(w http.ResponseWriter, r *http.Request, data []byte) {
	w.Write(data)
}

func Exit(start time.Time, w http.ResponseWriter, r *http.Request) {
	var send []byte
	var err error
	if xmux.GetInstance(r).Response != nil {
		response := xmux.GetInstance(r).Response.(*pkg.Response)
		response.Msg = pkg.ResponseMsg[response.Code]
		send, err = json.Marshal(response)
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
