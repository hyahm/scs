package handle

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
