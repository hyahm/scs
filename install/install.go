package install

import (
	"errors"
	"scs/internal"
	"scs/script"
	"strings"
)

type InstallConfig struct {
	Depend []string          `yaml:"depend"` // 依赖其他的包， 引入环境变量， 读取script的env
	Env    map[string]string `yaml:"env"`    // 外部定义的变量， 最后合并到env
	Script *internal.Script  `yaml:"script"`
}

func (ic *InstallConfig) GetDependEnv() error {
	for _, do := range ic.Depend {
		if _, ok := script.SS.Infos[do]; ok {
			for name := range script.SS.Infos[do] {
				for _, env := range script.SS.Infos[do][name].Env {
					i := strings.Index(env, "=")
					ic.Env[env[:i]] = ic.Env[env[i+1:]]
				}

			}
		} else {
			return errors.New("depend " + do + " not found, you need install first")
		}
	}

	return nil
}
