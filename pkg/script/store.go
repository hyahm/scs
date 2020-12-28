package script

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/hyahm/golog"
)

// 保存的所有脚本相关的配置
var SS *Service

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
			SS.Infos[pname][name].Start()
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
			if s.Status.Status != STOP && s.Status.Status != INSTALL {
				s.Status.Start = int64(time.Since(s.Status.Up).Seconds())
			} else {
				s.Status.Start = 0
			}
			s.Status.Command = s.Command
			ss = append(ss, s.Status)
		}

	}

	send, err := json.MarshalIndent(ss, "", "\t")
	if err != nil {
		golog.Error(err)
	}
	return send
}
