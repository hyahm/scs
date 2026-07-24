package config

import (
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/pkg"
)

type Role string

func (role Role) ToString() string {
	return string(role)
}

// 3种权限
const (
	AdminRole  Role = "admin"
	ScriptRole Role = "script"
	SimpleRole Role = "simple"
)

type Script struct {
	Name    string `yaml:"name,omitempty" json:"name"`
	Dir     string `yaml:"dir,omitempty" json:"dir,omitempty"`
	Command string `yaml:"command,omitempty" json:"command"`
	// ScriptToken  string            `yaml:"scriptToken,omitempty" json:"scriptToken,omitempty"` // 只用来查看的token
	// SimpleToken  string            `yaml:"simpleToken,omitempty" json:"simpleToken,omitempty"` // 角色权限
	Replicate    int               `yaml:"replicate,omitempty" json:"replicate,omitempty"`
	Always       bool              `yaml:"always,omitempty" json:"always,omitempty"`
	DisableAlert bool              `yaml:"disableAlert,omitempty" json:"disableAlert,omitempty"`
	Env          map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
	// ContinuityInterval time.Duration        `yaml:"continuityInterval,omitempty" json:"continuityInterval,omitempty"`
	Port int `yaml:"port,omitempty" json:"port,omitempty"`
	// AT             *AlertTo          `yaml:"alert,omitempty" json:"alert,omitempty"`
	Version        string            `yaml:"version,omitempty" json:"version,omitempty"`
	PreStart       []PreStart        `yaml:"preStart,omitempty" json:"preStart,omitempty"`
	Disable        bool              `yaml:"disable,omitempty" json:"disable,omitempty"`
	Cron           Cron              `yaml:"cron,omitempty" json:"cron,omitempty"`
	Update         string            `yaml:"update,omitempty" json:"update,omitempty"`
	DeleteWhenExit bool              `yaml:"deleteWhenExit,omitempty" json:"deleteWhenExit,omitempty"`
	TempEnv        map[string]string `yaml:"-" json:"-"`
	User           string            `yaml:"user,omitempty" json:"user,omitempty"`
	Group          string            `yaml:"group,omitempty" json:"group,omitempty"`
	StartTime      string            `yaml:"startTime,omitempty" json:"startTime,omitempty"` // 启动时间
	StopTime       string            `yaml:"stopTime,omitempty" json:"stopTime,omitempty"`   // 停止时间
	// Ready              chan bool            `yaml:"-" json:"-"`
	// 服务ready的探测器
	// Liveness *Liveness `yaml:"liveness,omitempty" json:"liveness,omitempty"`
}

// func (s *Script) Start() {
// 	// 这里转到 store
// 	// parameter := ""
// 	switch s.Status.Status {
// 	case status.STOP:
// 		// 开始启动的时候，需要将遍历变量值的模板渲染
// 		go s.asyncStart()
// 	}
// }

func (s *Script) MakeServer() {

	env := make(map[string]string)
	for k, v := range s.TempEnv {
		env[k] = v
	}
	if s.Port > 0 {
		// 顺序拿到可用端口
		s.Port = pkg.GetAvailablePort(s.Port)
		env["PORT"] = strconv.Itoa(s.Port)
	} else {
		env["PORT"] = "0"
	}
	// 填充server到store
	// s.FillServer()
	// env["NAME"] = s.SubName

	s.Env = env
}

// func (s *Script) SetToken(token string) {

// }

// func (s *Script) FillServer(subname string) {
// store := server.GetStore()
// svc := store.GetServerBySubName(subname)
// svc.ScriptToken = s.ScriptToken
// svc.SimpleToken = s.SimpleToken
// if svc.SimpleToken == "" {
// 	svc.SimpleToken = pkg.RandomToken()
// }
// svc.Command = script.Command
// svc.User = script.User
// svc.Group = script.Group
// svc.Disable = script.Disable
// // Log:       make([]string, 0, global.GetLogCount()),
// svc.Dir = script.Dir
// if svc.Status == nil {
// 	svc.Status = &status.Status{
// 		Status: status.STOP,
// 	}
// }
// svc.StartTime = script.StartTime
// svc.StopTime = script.StopTime

// svc.Update = script.Update
// svc.AI = &config.AlertInfo{}
// // svc.AT = script.AT
// svc.StopSignal = make(chan bool, 1)

// // svc.Liveness = script.Liveness
// svc.Ready = make(chan bool, 1)
// svc.Always = script.Always
// svc.AlwaysSign = script.Always
// svc.DeleteWhenExit = script.DeleteWhenExit

// // svc.DisableAlert = script.DisableAlert
// svc.PreStart = script.PreStart
// svc.Cron = script.Cron
// if svc.Cron != nil {
// 	svc.Cron = &config.Cron{
// 		Start:   script.Cron.Start,
// 		Loop:    script.Cron.Loop,
// 		IsMonth: script.Cron.IsMonth,
// 		Times:   script.Cron.Times,
// 	}
// }
// }

// 生成新的env 到 tempenv
func (s *Script) MakeTempEnv() map[string]string {
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
	// tempEnv["TOKEN"] = s.ScriptToken
	tempEnv["PNAME"] = s.Name
	tempEnv["PROJECT_HOME"] = s.Dir

	return tempEnv
}

func (s *Script) GetEnv() []string {
	env := make([]string, 0, len(s.Env))
	for k, v := range s.Env {
		env = append(env, k+"="+v)
	}
	return env
}

func EqualScript(s1, s2 Script) bool {
	// if s1 == nil && s2 != nil || s1 != nil && s2 == nil {
	// 	return false
	// }
	// if s1 == nil && s2 == nil {
	// 	return true
	// }
	// 这些有一个不同的。 那么就需要重启所有底下的server
	return !(s1.Name != s2.Name ||
		s1.Dir != s2.Dir ||
		s1.Command != s2.Command ||
		s1.Always != s2.Always ||
		// s1.ScriptToken != s2.ScriptToken ||
		!pkg.CompareMap(s1.Env, s2.Env) ||
		// !CompareAT(s1.AT, s2.AT) ||
		s1.DisableAlert != s2.DisableAlert ||
		s1.Disable != s2.Disable ||
		s1.Update != s2.Update ||
		s1.User != s2.User ||
		s1.Group != s2.Group ||
		!EqualPreStart(s1.PreStart, s2.PreStart) ||
		s1.Version != s2.Version ||
		!CompareCron(s1.Cron, s2.Cron))
}
