package module

import (
	"net/http"

	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/pkg"
)

func UpdateConfig(w http.ResponseWriter, r *http.Request) bool {
	msg, ok := global.IsCanReload()
	if !ok {
		pkg.Error(r, msg)
	}
	return false
}
