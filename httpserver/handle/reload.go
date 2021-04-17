package handle

import (
	"fmt"
	"net/http"

	"github.com/hyahm/scs"
)

var reloadKey bool

func Reload(w http.ResponseWriter, r *http.Request) {
	// 关闭上次监控的goroutine
	if reloadKey {
		w.Write([]byte(`{"code": 201, "msg": "config file is reloading, waiting completed first"}`))
		return
	}

	// script.Reloadlocker.Lock()
	// defer script.Reloadlocker.Unlock()
	reloadKey = true
	defer func() {
		reloadKey = false
	}()
	// 拷贝一份到当前运行的脚本列表
	if err := scs.ReLoad(); err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code": 500, "msg": "%s"}`, err.Error())))
		return
	}
	w.Write([]byte(`{"code": 200, "msg": "config file reloaded"}`))
}
