package script

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/hyahm/golog"
)

// 保存的所有脚本相关的配置
var SS *Service

const (
	STOP        string = "Stop"
	RUNNING     string = "Running"
	WAITSTOP    string = "Waiting Stop"
	WAITRESTART string = "Waiting Restart"
)

func init() {
	SS = &Service{
		Infos: make(map[string]map[string]*Script),
		mu:    &sync.RWMutex{},
	}
}

type Service struct {
	Infos map[string]map[string]*Script
	mu    *sync.RWMutex
}

type ServiceStatus struct {
	Name         string `json:"name"`
	Ppid         int    `json:"ppid"`
	Status       string `json:"status"`
	Command      string `json:"command"`
	PName        string `json:"pname"`
	Path         string `json:"path"`
	CanNotStop   bool   `json:"cannotStop"` //
	Stoping      bool   `json:"stoping,omitempty"`
	Up           int64
	Start        string `json:"start"` // 启动的时间
	Version      string `json:"version"`
	Always       bool
	RestartCount int `json:"restartCount"` // 记录失败重启的次数
}

type status struct {
	s     bool
	pname string
	name  string
}

// 第一次启动
func (s *Service) Start() {
	// 先无视有启动顺序的脚本
	// 先启动无需顺序的脚本
	s.mu.RLock()
	for pname := range SS.Infos {
		for name := range SS.Infos[pname] {
			SS.Infos[pname][name].Start(SS.Infos[pname][name].Command)
		}
	}
	s.mu.RUnlock()

}

func Get(pname, name string) *ServiceStatus {
	if _, ok := SS.Infos[pname]; ok {
		if sv, ok := SS.Infos[pname][name]; ok {
			return sv.Status
		}

	}
	return nil
}

func All() []byte {
	if SS.mu == nil {
		SS.mu = &sync.RWMutex{}
	}

	ss := make([]*ServiceStatus, 0)
	for pname := range SS.Infos {
		for _, s := range SS.Infos[pname] {
			if s.Status.Up > 0 {
				s.Status.Start = (time.Duration(time.Now().Unix()-s.Status.Up) * time.Second).String()
			} else {
				s.Status.Start = "0s"
			}
			ss = append(ss, s.Status)
		}

	}
	send, err := json.MarshalIndent(ss, "", "\t")
	if err != nil {
		golog.Error(err)
	}
	return send
}
