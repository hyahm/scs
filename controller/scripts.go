package controller

import (
	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/scripts"
)

func NeedStop(s *scripts.Script) bool {
	// 更新server
	// 判断值是否相等
	script, ok := store.Store.GetScriptByName(s.Name)
	if !ok {
		golog.Error(pkg.ErrBugMsg)
		return false
	}
	return !scripts.EqualScript(s, script)
}
