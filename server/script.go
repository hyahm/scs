package server

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/alert"
	"github.com/hyahm/scs/cron"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/lookpath"
	"github.com/hyahm/scs/status"
	"github.com/hyahm/scs/subname"
	"github.com/hyahm/scs/to"
	"github.com/hyahm/scs/utils"
)

type Script struct {
	Name               string               `yaml:"name,omitempty" json:"name"`
	Dir                string               `yaml:"dir,omitempty" json:"dir"`
	Command            string               `yaml:"command,omitempty" json:"command"`
	Replicate          int                  `yaml:"replicate,omitempty" json:"replicate,omitempty"`
	Always             bool                 `yaml:"always,omitempty" json:"always,omitempty"`
	DisableAlert       bool                 `yaml:"disableAlert,omitempty" json:"disableAlert,omitempty"`
	Env                map[string]string    `yaml:"env,omitempty" json:"env,omitempty"`
	ContinuityInterval time.Duration        `yaml:"continuityInterval,omitempty" json:"continuityInterval,omitempty"`
	Port               int                  `yaml:"port,omitempty" json:"port,omitempty"`
	AT                 *to.AlertTo          `yaml:"alert,omitempty" json:"alert,omitempty"`
	Version            string               `yaml:"version,omitempty" json:"version,omitempty"`
	Loop               int                  `yaml:"loop,omitempty" json:"loop,omitempty"`
	LookPath           []*lookpath.LoopPath `yaml:"lookPath,omitempty" json:"lookPath,omitempty"`
	Disable            bool                 `yaml:"disable,omitempty" json:"disable,omitempty"`
	Cron               *cron.Cron           `yaml:"cron,omitempty" json:"cron,omitempty"`
	Update             string               `yaml:"update,omitempty" json:"update,omitempty"`
	DeleteWhenExit     bool                 `yaml:"deleteWhenExit,omitempty" json:"deleteWhenExit,omitempty"`
	TempEnv            map[string]string    `yaml:"-" json:"-"`
}

// 生成新的env 到 tempenv
func (s *Script) MakeEnv() {
	// 生成 全局脚本的 env
	if s.TempEnv == nil {
		s.TempEnv = make(map[string]string)
	}

	pathEnvName := "PATH"
	for _, v := range os.Environ() {
		kv := strings.Split(v, "=")
		if strings.ToUpper(kv[0]) == pathEnvName {
			pathEnvName = kv[0]
		}
		s.TempEnv[kv[0]] = kv[1]
	}
	for k, v := range s.Env {
		// path 环境单独处理， 可以多个值， 其他环境变量多个值请以此写完
		if strings.EqualFold(k, pathEnvName) {
			if runtime.GOOS == "windows" {
				s.TempEnv[pathEnvName] = s.TempEnv[pathEnvName] + ";" + v
			} else {
				golog.Info(pathEnvName)
				s.TempEnv[pathEnvName] = s.TempEnv[pathEnvName] + ":" + v
			}
		} else {
			s.TempEnv[k] = v
		}
	}

	s.TempEnv["TOKEN"] = global.Token
	s.TempEnv["PNAME"] = s.Name
}

func getVersion(command string) string {
	var cmd *exec.Cmd
	golog.Info(command)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-c", command)
	} else {
		cmd = exec.Command("/bin/bash", "-c", command)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	s := new(string)
	ch := make(chan struct{})
	defer cancel()
	go func(s *string) {
		out, err := cmd.Output()
		if err != nil {
			return
		}
		*s = string(out)
		ch <- struct{}{}
	}(s)

	select {
	case <-ctx.Done():
		return ""
	case <-ch:
		output := strings.ReplaceAll(*s, "\n", "")
		output = strings.ReplaceAll(output, "\r", "")
		return output
	}

}

func (s *Script) add(port int, subname subname.Subname) *Server {
	continuityInterval := s.ContinuityInterval
	if continuityInterval == 0 {
		continuityInterval = global.ContinuityInterval
	}
	svc := &Server{
		// LookPath:  s.LookPath,
		Script:    s,
		Command:   s.Command,
		Log:       make([]string, 0, global.LogCount),
		LogLocker: &sync.RWMutex{},
		SubName:   subname,
		Version:   getVersion(s.Version),
		Status: &status.ServiceStatus{
			Name:   subname.GetName(),
			PName:  s.Name,
			Status: status.STOP,
		},
		Update:             s.Update,
		ContinuityInterval: continuityInterval,
		AI:                 &alert.AlertInfo{},
		Port:               port,
		AT:                 s.AT,
		StopSigle:          make(chan bool, 1),
	}
	if s.Cron != nil {
		svc.Cron = &cron.Cron{
			Start:   s.Cron.Start,
			Loop:    s.Cron.Loop,
			IsMonth: s.Cron.IsMonth,
		}
	}
	return svc
}

func (s *Script) GetEnv() []string {
	env := make([]string, 0, len(s.Env))
	for k, v := range s.Env {
		env = append(env, k+"="+v)
	}
	golog.Info(env)
	return env
}

// 通过script 生成 server
func (s *Script) MakeServer() {
	ss.Mu.Lock()
	defer ss.Mu.Unlock()
	ss.Scripts[s.Name] = s
	s.MakeEnv()
	replica := s.Replicate
	if replica == 0 {
		replica = 1
	}
	portIndex := 0
	for i := 0; i < replica; i++ {
		// 根据副本数提取子名称
		env := make(map[string]string)
		for k, v := range s.TempEnv {
			env[k] = v
		}
		subname := subname.NewSubname(s.Name, i)
		var svc *Server
		if s.Port > 0 {
			// 检测端口是否被占用， 如果占用了
			portIndex += utils.ProbePort(s.Port)
			env["PORT"] = strconv.Itoa(s.Port + i + portIndex)
			svc = s.add(s.Port+i+portIndex, subname)
		} else {
			env["PORT"] = "0"
			svc = s.add(0, subname)
		}

		env["NAME"] = subname.String()
		svc.Env = env
		golog.Debug("start " + subname)
		ss.Infos[subname] = svc
		ss.Infos[subname].Start()
	}
}

// 通过script 生成 server
func (s *Script) MakeReplicateServer(start, end int) {
	ss.Mu.Lock()
	defer ss.Mu.Unlock()
	ss.Scripts[s.Name].Replicate = s.Replicate
	s.MakeEnv()
	// replica := s.Replicate
	// if replica == 0 {
	// 	replica = 1
	// }
	portIndex := 0
	for i := start; i < end; i++ {
		// 根据副本数提取子名称
		env := make(map[string]string)
		for k, v := range s.TempEnv {
			env[k] = v
		}
		subname := subname.NewSubname(s.Name, i)
		var svc *Server
		if s.Port > 0 {
			// 检测端口是否被占用， 如果占用了
			portIndex += utils.ProbePort(s.Port)
			env["PORT"] = strconv.Itoa(s.Port + i + portIndex)
			svc = s.add(s.Port+i+portIndex, subname)
		} else {
			env["PORT"] = "0"
			svc = s.add(0, subname)
		}

		env["NAME"] = subname.String()
		svc.Env = env
		golog.Debug("start " + subname)
		ss.Infos[subname] = svc
		ss.Infos[subname].Start()
	}
}

func CompareScript(s1, s2 *Script) bool {
	if s1 == nil && s2 != nil || s1 != nil && s2 == nil {
		return false
	}
	if s1 == nil && s2 == nil {
		return true
	}
	// 这些有一个不同的。 那么就需要重启所有底下的server
	if s1.Name != s2.Name ||
		s1.Dir != s2.Dir ||
		s1.Command != s2.Command ||
		s1.Always != s2.Always ||
		!utils.CompareMap(s1.Env, s2.Env) ||
		s1.ContinuityInterval != s2.ContinuityInterval ||
		!to.CompareAT(s1.AT, s2.AT) ||
		s1.DisableAlert != s2.DisableAlert ||
		s1.Disable != s2.Disable ||
		s1.Version != s2.Version ||
		!cron.CompareCron(s1.Cron, s2.Cron) ||
		s1.Port != s2.Port {
		return false
	}

	return true
}
