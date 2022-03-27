package scripts

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/internal/config/alert"
	"github.com/hyahm/scs/internal/config/alert/to"
	"github.com/hyahm/scs/internal/config/liveness"
	"github.com/hyahm/scs/internal/config/scripts/cron"
	"github.com/hyahm/scs/internal/config/scripts/prestart"
	"github.com/hyahm/scs/internal/config/scripts/subname"
	"github.com/hyahm/scs/internal/server"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/status"
)

type Script struct {
	Token              string               `yaml:"token,omitempty" json:"token,omitempty"` // 只用来查看的token
	Name               string               `yaml:"name,omitempty" json:"name"`
	Dir                string               `yaml:"dir,omitempty" json:"dir,omitempty"`
	Command            string               `yaml:"command,omitempty" json:"command"`
	Replicate          int                  `yaml:"replicate,omitempty" json:"replicate,omitempty"`
	Always             bool                 `yaml:"always,omitempty" json:"always,omitempty"`
	DisableAlert       bool                 `yaml:"disableAlert,omitempty" json:"disableAlert,omitempty"`
	Env                map[string]string    `yaml:"env,omitempty" json:"env,omitempty"`
	ContinuityInterval time.Duration        `yaml:"continuityInterval,omitempty" json:"continuityInterval,omitempty"`
	Port               int                  `yaml:"port,omitempty" json:"port,omitempty"`
	AT                 *to.AlertTo          `yaml:"alert,omitempty" json:"alert,omitempty"`
	Version            string               `yaml:"version,omitempty" json:"version,omitempty"`
	PreStart           []*prestart.PreStart `yaml:"preStart,omitempty" json:"preStart,omitempty"`
	Disable            bool                 `yaml:"disable,omitempty" json:"disable,omitempty"`
	Cron               *cron.Cron           `yaml:"cron,omitempty" json:"cron,omitempty"`
	Update             string               `yaml:"update,omitempty" json:"update,omitempty"`
	DeleteWhenExit     bool                 `yaml:"deleteWhenExit,omitempty" json:"deleteWhenExit,omitempty"`
	TempEnv            map[string]string    `yaml:"-" json:"-"`
	EnvLocker          *sync.RWMutex        `yaml:"-" json:"-"`
	// Ready              chan bool            `yaml:"-" json:"-"`
	// 服务ready的探测器
	Liveness *liveness.Liveness `yaml:"liveness,omitempty" json:"liveness,omitempty"`
}

// 生成新的env 到 tempenv
func (s *Script) MakeEnv() {
	// 生成 全局脚本的 env
	tempEnv := make(map[string]string)
	if s.EnvLocker == nil {
		s.EnvLocker = &sync.RWMutex{}
	}

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
	tempEnv["TOKEN"] = global.GetToken()
	tempEnv["PNAME"] = s.Name
	tempEnv["PROJECT_HOME"] = s.Dir
	for k := range tempEnv {
		if len(k) > 8 && k[:7] == "SCS_TPL" {
			tempEnv[k] = internal.Format(tempEnv[k], tempEnv)
		}
	}

	s.Command = internal.Format(s.Command, tempEnv)
	s.EnvLocker.Lock()
	defer s.EnvLocker.Unlock()
	s.TempEnv = tempEnv
}

// 生成 server
func (s *Script) Add(port, replicate, id int, subname subname.Subname) *server.Server {
	continuityInterval := s.ContinuityInterval
	if continuityInterval == 0 {
		continuityInterval = global.GeContinuityInterval()
	}

	svc := &server.Server{
		// Script:  s,
		Name:    s.Name,
		Index:   id,
		Token:   s.Token,
		Command: s.Command,
		// Log:       make([]string, 0, global.GetLogCount()),
		SubName: subname,
		Dir:     s.Dir,
		Status: &status.ServiceStatus{
			Name:    subname.GetName(),
			PName:   s.Name,
			Status:  status.STOP,
			Command: s.Command,
			Path:    s.Dir,
			OS:      runtime.GOOS,
		},
		Replicate: replicate,
		Logger: golog.NewLog(
			filepath.Join(global.LogDir, subname.String()+".log"), 10<<10, false, global.CleanLog),
		Update:             s.Update,
		ContinuityInterval: continuityInterval,
		AI:                 &alert.AlertInfo{},
		Port:               port,
		AT:                 s.AT,
		StopSigle:          make(chan bool, 1),

		Liveness:     s.Liveness,
		Ready:        make(chan bool, 1),
		Always:       s.Always,
		Disable:      s.Disable,
		DisableAlert: s.DisableAlert,
		PreStart:     s.PreStart,
	}
	svc.Logger.Format = global.FORMAT
	if s.Cron != nil {
		svc.Cron = &cron.Cron{
			Start:   s.Cron.Start,
			Loop:    s.Cron.Loop,
			IsMonth: s.Cron.IsMonth,
			Times:   s.Cron.Times,
		}
	}
	return svc
}

func (s *Script) GetEnv() []string {
	env := make([]string, 0, len(s.Env))
	for k, v := range s.Env {
		env = append(env, k+"="+v)
	}
	return env
}

// 通过script 生成 server
// func (s *Script) MakeServer(ss map[string]*server.Server) {
// 	s.MakeEnv()
// 	replicate := s.Replicate
// 	if replicate == 0 {
// 		replicate = 1
// 	}
// 	availablePort := s.Port

// 	for i := 0; i < replicate; i++ {
// 		// 根据副本数提取子名称
// 		env := make(map[string]string)
// 		for k, v := range s.TempEnv {
// 			env[k] = v
// 		}
// 		svc := &server.Server{}
// 		subname := subname.NewSubname(s.Name, i)
// 		if availablePort > 0 {
// 			// 检测端口是否被占用， 如果占用了
// 			availablePort = pkg.GetAvailablePort(s.Port)

// 		}
// 		env["PORT"] = strconv.Itoa(availablePort)
// 		svc = s.Add(availablePort, replicate, subname)

// 		env["NAME"] = subname.String()
// 		svc.Env = env
// 		ss[subname.String()] = svc
// 	}
// }

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
		s1.ContinuityInterval != s2.ContinuityInterval ||
		!to.CompareAT(s1.AT, s2.AT) ||
		s1.DisableAlert != s2.DisableAlert ||
		s1.Disable != s2.Disable ||
		s1.Update != s2.Update ||
		!prestart.EqualPreStart(s1.PreStart, s2.PreStart) ||
		s1.Version != s2.Version ||
		!cron.CompareCron(s1.Cron, s2.Cron))
}
