package to

import "github.com/hyahm/scs/utils"

type AlertTo struct {
	Email    []string `yaml:"email"`
	Rocket   []string `yaml:"rocket"`
	Telegram []string `yaml:"telegram"`
	WeiXin   []string `yaml:"weixin"`
	Callback []string `yaml:"callback"`
}

func CompareAT(a1, a2 *AlertTo) bool {
	if a1 == nil && a2 != nil || a1 != nil && a2 == nil {
		return false
	}
	if a1 == nil && a2 == nil {
		return true
	}
	if !utils.CompareSlice(a1.Email, a2.Email) ||
		!utils.CompareSlice(a1.Rocket, a2.Rocket) ||
		!utils.CompareSlice(a1.Telegram, a2.Telegram) ||
		!utils.CompareSlice(a1.WeiXin, a2.WeiXin) {
		return false
	}
	return true
}
