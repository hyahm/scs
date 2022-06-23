package scripts

import (
	"os"
	"runtime"
	"strings"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/alert/to"
	"github.com/hyahm/scs/pkg/config/liveness"
	"github.com/hyahm/scs/pkg/config/scripts/cron"
	"github.com/hyahm/scs/pkg/config/scripts/prestart"
)

type Role string

func (role Role) ToString() string {
	return string(role)
}

// 3种权限
const (
	AdminRole  Role = "admin"
	ScriptRole Role = "script"
	Server     Role = "server"
)

type Script struct {
	Name         string            `yaml:"name,omitempty" json:"name"`
	Dir          string            `yaml:"dir,omitempty" json:"dir,omitempty"`
	Command      string            `yaml:"command,omitempty" json:"command"`
	Token        string            `yaml:"token,omitempty" json:"token,omitempty"` // 只用来查看的token
	Role         Role              `yaml:"role,omitempty" json:"role,omitempty"`   // 角色权限
	Replicate    int               `yaml:"replicate,omitempty" json:"replicate,omitempty"`
	Always       bool              `yaml:"always,omitempty" json:"always,omitempty"`
	DisableAlert bool              `yaml:"disableAlert,omitempty" json:"disableAlert,omitempty"`
	Env          map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
	// ContinuityInterval time.Duration        `yaml:"continuityInterval,omitempty" json:"continuityInterval,omitempty"`
	Port           int                  `yaml:"port,omitempty" json:"port,omitempty"`
	AT             *to.AlertTo          `yaml:"alert,omitempty" json:"alert,omitempty"`
	Version        string               `yaml:"version,omitempty" json:"version,omitempty"`
	PreStart       []*prestart.PreStart `yaml:"preStart,omitempty" json:"preStart,omitempty"`
	Disable        bool                 `yaml:"disable,omitempty" json:"disable,omitempty"`
	Cron           *cron.Cron           `yaml:"cron,omitempty" json:"cron,omitempty"`
	Update         string               `yaml:"update,omitempty" json:"update,omitempty"`
	DeleteWhenExit bool                 `yaml:"deleteWhenExit,omitempty" json:"deleteWhenExit,omitempty"`
	TempEnv        map[string]string    `yaml:"-" json:"-"`
	// Ready              chan bool            `yaml:"-" json:"-"`
	// 服务ready的探测器
	Liveness *liveness.Liveness `yaml:"liveness,omitempty" json:"liveness,omitempty"`
}

// 生成新的env 到 tempenv
func (s *Script) MakeTempEnv() {
	// 生成 全局脚本的 env
	tempEnv := make(map[string]string)

	pathEnvName := "PATH"
	for _, v := range os.Environ() {
		kv := strings.Split(v, "=")
		if strings.ToUpper(kv[0]) == pathEnvName {
			pathEnvName = kv[0]
		}
		tempEnv[kv[0]] = kv[1]
	}
	for k, v := range s.Env {
		// path 环境单独处理， 可以多个值， 其他环境变量多个值请以此写完
		if strings.EqualFold(k, pathEnvName) {
			if runtime.GOOS == "windows" {
				tempEnv[pathEnvName] = tempEnv[pathEnvName] + ";" + v
			} else {
				golog.Info(pathEnvName)
				tempEnv[pathEnvName] = tempEnv[pathEnvName] + ":" + v
			}
		} else {
			tempEnv[k] = v
		}
	}
	tempEnv["OS"] = runtime.GOOS
	// 增加token, 不过是随机的
	tempEnv["TOKEN"] = s.Token
	tempEnv["PNAME"] = s.Name
	tempEnv["PROJECT_HOME"] = s.Dir

	s.TempEnv = tempEnv
}

func (s *Script) GetEnv() []string {
	env := make([]string, 0, len(s.Env))
	for k, v := range s.Env {
		env = append(env, k+"="+v)
	}
	return env
}

func EqualScript(s1, s2 *Script) bool {
	if s1 == nil && s2 != nil || s1 != nil && s2 == nil {
		return false
	}
	if s1 == nil && s2 == nil {
		return true
	}
	// 这些有一个不同的。 那么就需要重启所有底下的server
	return !(s1.Name != s2.Name ||
		s1.Dir != s2.Dir ||
		s1.Command != s2.Command ||
		s1.Always != s2.Always ||
		s1.Token != s2.Token ||
		!pkg.CompareMap(s1.Env, s2.Env) ||
		!to.CompareAT(s1.AT, s2.AT) ||
		s1.DisableAlert != s2.DisableAlert ||
		s1.Disable != s2.Disable ||
		s1.Update != s2.Update ||
		!prestart.EqualPreStart(s1.PreStart, s2.PreStart) ||
		s1.Version != s2.Version ||
		!cron.CompareCron(s1.Cron, s2.Cron))
}
