package handle

import (
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs"
)

// func GetConfig(w http.ResponseWriter, r *http.Request) {
// 	name := xmux.Var(r)["name"]
// 	golog.Info(name)
// 	conf := clientv3.Config{
// 		Endpoints: []string{"http://127.0.0.1:2379"},
// 	}
// 	cli, err := clientv3.New(conf)
// 	if err != nil {
// 		golog.Error(err)
// 		return
// 	}
// 	res, err := cli.Get(context.Background(), name)
// 	if err != nil {
// 		golog.Error(err)
// 		return
// 	}
// 	for _, v := range res.Kvs {
// 		golog.Info(v.Value)
// 	}

// 	return
// }

func GetServers(w http.ResponseWriter, r *http.Request) {
	golog.Info("servers")
	w.Write(scs.GetServers())
}

func GetScripts(w http.ResponseWriter, r *http.Request) {
	w.Write(scs.GetScripts())
}
