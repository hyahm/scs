package script

import (
	"encoding/json"
	"sync"

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

type StatusList struct {
	Data []ServiceStatus `json:"data"`
	Code int             `json:"code"`
}

func All() []byte {
	if SS.mu == nil {
		SS.mu = &sync.RWMutex{}
	}
	statuss := &StatusList{
		Data: make([]ServiceStatus, 0),
	}
	// ss := make([]*ServiceStatus, 0)
	for pname := range SS.Infos {
		for _, s := range SS.Infos[pname] {
			s.Status.Command = s.Command
			statuss.Data = append(statuss.Data, *s.Status)
		}

	}
	statuss.Code = 200
	send, err := json.MarshalIndent(statuss, "", "\t")
	if err != nil {
		golog.Error(err)
	}
	return send
}

func ScriptPname(pname string) []byte {
	if SS.mu == nil {
		SS.mu = &sync.RWMutex{}
	}
	statuss := &StatusList{
		Data: make([]ServiceStatus, 0),
	}
	for _, s := range SS.Infos[pname] {
		statuss.Data = append(statuss.Data, *s.Status)
	}
	statuss.Code = 200
	send, err := json.MarshalIndent(statuss, "", "\n")

	if err != nil {
		golog.Error(err)
	}
	return send
}

func ScriptName(pname, name string) []byte {
	if SS.mu == nil {
		SS.mu = &sync.RWMutex{}
	}

	statuss := &StatusList{
		Data: make([]ServiceStatus, 0),
	}
	// statuss := make([]*script.ServiceStatus, 0)
	if _, pok := SS.Infos[pname]; pok {
		if s, ok := SS.Infos[pname][name]; ok {
			statuss.Data = append(statuss.Data, *s.Status)
		}
	}
	statuss.Code = 200

	send, err := json.MarshalIndent(statuss, "", "\t")
	if err != nil {
		golog.Error(err)
	}
	return send
}
