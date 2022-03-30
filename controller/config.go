package controller

import "github.com/hyahm/scs/pkg/config"

var cfg *config.Config

// DelScript 删除脚本的时候才会执行
func DelScript(pname string) error {
	err := RemoveScript(pname)
	if err != nil {
		return err
	}

	for i, s := range cfg.SC {
		if s.Name == pname {
			cfg.SC = append(cfg.SC[:i], cfg.SC[i+1:]...)
			break
		}
	}
	return cfg.WriteConfig(true)
}
