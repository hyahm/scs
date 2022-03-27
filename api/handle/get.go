package handle

import (
	"net/http"

	"github.com/hyahm/scs/controller"
	"github.com/hyahm/xmux"
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

func GetAlarms(w http.ResponseWriter, r *http.Request) {
	role := xmux.GetInstance(r).Get("role").(string)
	res := &Response{
		Data: controller.GetAterts(),
		Role: role,
	}
	w.Write(res.Sucess(""))
}

func GetServers(w http.ResponseWriter, r *http.Request) {
	role := xmux.GetInstance(r).Get("role").(string)
	res := &Response{
		Data: controller.GetServers(),
		Role: role,
	}
	w.Write(res.Sucess(""))
}

func GetScripts(w http.ResponseWriter, r *http.Request) {
	role := xmux.GetInstance(r).Get("role").(string)
	res := &Response{
		Data: controller.GetScripts(),
		Role: role,
	}
	w.Write(res.Sucess(""))
}
