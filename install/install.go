package install

import "github.com/hyahm/scs/server"

type InstallConfig struct {
	Depend []string          `yaml:"depend"` // 依赖其他的包， 引入环境变量， 读取script的env
	Env    map[string]string `yaml:"env"`    // 外部定义的变量， 最后合并到env
	Script *server.Script    `yaml:"script"`
}
