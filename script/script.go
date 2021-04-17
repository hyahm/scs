package script

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
)

type Script struct {
	Name               string            `yaml:"name,omitempty" json:"name"`
	Dir                string            `yaml:"dir,omitempty" json:"dir"`
	Command            string            `yaml:"command,omitempty" json:"command"`
	Replicate          int               `yaml:"replicate,omitempty" json:"replicate,omitempty"`
	Always             bool              `yaml:"always,omitempty" json:"always,omitempty"`
	DisableAlert       bool              `yaml:"disableAlert,omitempty" json:"disableAlert,omitempty"`
	Env                map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
	ContinuityInterval time.Duration     `yaml:"continuityInterval,omitempty" json:"continuityInterval,omitempty"`
	Port               int               `yaml:"port,omitempty" json:"port,omitempty"`
	AT                 *AlertTo          `yaml:"alert,omitempty" json:"alert,omitempty"`
	Version            string            `yaml:"version,omitempty" json:"version,omitempty"`
	Loop               int               `yaml:"loop,omitempty" json:"loop,omitempty"`
	LookPath           []*LoopPath       `yaml:"lookPath,omitempty" json:"lookPath,omitempty"`
	Disable            bool              `yaml:"disable,omitempty" json:"disable,omitempty"`
	Cron               *Cron             `yaml:"cron,omitempty" json:"cron,omitempty"`
	Update             string            `yaml:"update,omitempty" json:"update,omitempty"`
	DeleteWhenExit     bool              `yaml:"deleteWhenExit,omitempty" json:"deleteWhenExit,omitempty"`
	TempEnv            map[string]string `yaml:"-" json:"-"`
}

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

func (s *Script) add(port, replacate int, subname string) *Server {
	continuityInterval := s.ContinuityInterval
	if continuityInterval == 0 {
		continuityInterval = global.ContinuityInterval
	}

	svc := &Server{
		Name:      s.Name,
		LookPath:  s.LookPath,
		Command:   s.Command,
		Dir:       s.Dir,
		Replicate: replacate,
		Log:       make(map[string][]string),
		LogLocker: &sync.RWMutex{},
		SubName:   subname,
		Status: &ServiceStatus{
			Name:    subname,
			PName:   s.Name,
			Status:  STOP,
			Path:    s.Dir,
			Version: getVersion(s.Version),
			Disable: s.Disable,
		},
		DeleteWhenExit:     s.DeleteWhenExit,
		Update:             s.Update,
		DisableAlert:       s.DisableAlert,
		ContinuityInterval: continuityInterval,
		Always:             s.Always,
		Disable:            s.Disable,
		AI:                 &AlertInfo{},
		Port:               port,
		AT:                 s.AT,
		StopSigle:          make(chan bool, 1),
		Cron:               s.Cron,
	}
	// 生成对应的文件类型
	svc.Log["log"] = make([]string, 0, global.LogCount)
	svc.Log["lookPath"] = make([]string, 0, global.LogCount)
	svc.Log["update"] = make([]string, 0, global.LogCount)

	return svc
}

func (s *Script) UpdateServer() {
	// 更新server
	// 判断值是否相等

	if s.Dir != ss.Scripts[s.Name].Dir ||
		s.Command != ss.Scripts[s.Name].Command ||
		s.Replicate != ss.Scripts[s.Name].Replicate ||
		s.Always != ss.Scripts[s.Name].Always ||
		s.DisableAlert != ss.Scripts[s.Name].DisableAlert ||
		!EqualMap(s.Env, ss.Scripts[s.Name].Env) ||
		s.Port != ss.Scripts[s.Name].Port ||
		s.Version != ss.Scripts[s.Name].Version ||
		s.Disable != ss.Scripts[s.Name].Disable ||
		s.Update != ss.Scripts[s.Name].Update ||
		s.DeleteWhenExit != ss.Scripts[s.Name].DeleteWhenExit ||
		s.Cron.Start != ss.Scripts[s.Name].Cron.Start ||
		s.Cron.Loop != ss.Scripts[s.Name].Cron.Loop ||
		s.Cron.IsMonth != ss.Scripts[s.Name].Cron.IsMonth ||
		!EqualStringArray(s.AT.Email, ss.Scripts[s.Name].AT.Email) ||
		!EqualStringArray(s.AT.Rocket, ss.Scripts[s.Name].AT.Rocket) ||
		!EqualStringArray(s.AT.Telegram, ss.Scripts[s.Name].AT.Telegram) ||
		!EqualStringArray(s.AT.WeiXin, ss.Scripts[s.Name].AT.WeiXin) {
		// 如果值有变动， 那么需要重新server
		// 先同步停止之前的server， 然后启动新的server
		// server 是单独的， 在通知后需要同步更新server

	}

}

// 启动方法， 异步执行
func (s *Script) StartServer() {

	for _, svc := range ss.Infos[s.Name] {
		err := svc.Start()
		if err != nil {
			golog.Error(err)
		}
	}
}

func (s *Script) MakeServer() {
	// 通过script 生成 server
	s.MakeEnv()

	replica := s.Replicate
	if replica == 1 || replica == 0 {
		replica = 1
	}
	if ss.Infos[s.Name] == nil {
		ss.Infos[s.Name] = make(map[string]*Server)
	}
	portIndex := 0
	for i := 0; i < replica; i++ {
		// 根据副本数提取子名称
		env := make(map[string]string)
		for k, v := range s.TempEnv {
			env[k] = v
		}
		subname := fmt.Sprintf("%s_%d", s.Name, i)
		var svc *Server
		if s.Port > 0 {
			portIndex += probePort(s.Port)
			env["PORT"] = strconv.Itoa(s.Port + i + portIndex)
			svc = s.add(s.Port+i+portIndex, replica, subname)
		} else {
			env["PORT"] = "0"
			svc = s.add(0, replica, subname)
		}
		env["NAME"] = subname

		// 检测端口是否被占用， 如果占用了

		svc.Env = env
		ss.Infos[s.Name][subname] = svc
	}
}

func probePort(port int) int {
	// 检测端口
	index := 0
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf(":%d", port), time.Nanosecond*100)
		if err != nil {
			return index
		}
		if conn != nil {
			_ = conn.Close()
			return index
		}
		index++
	}

}
