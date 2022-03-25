package controller

import "github.com/hyahm/scs/internal/config/scripts/subname"

// 通过名字来获取token

func GetLogToken(name string) string {
	mu.RLock()
	defer mu.RUnlock()

	if v, ok := ss[subname.Subname(name).GetName()]; ok {
		return v.Token
	}
	return ""
}
