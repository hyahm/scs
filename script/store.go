package script

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/hyahm/golog"
)

// 保存的所有脚本相关的配置
var SS *Service

func init() {
	SS = &Service{
		// 由2层组成， 一级是name  二级是pname
		Infos: make(map[string]map[string]*Script),
		Mu:    &sync.RWMutex{},
	}
}

type Service struct {
	Infos map[string]map[string]*Script
	Mu    *sync.RWMutex
}

func getpname(subname string) string {
	i := strings.LastIndex(subname, "_")
	if i < 0 {
		return ""
	}
	return subname[:i]
}

func (s *Service) HasName(name string) bool {
	if s.Len() == 0 {
		return false
	}
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if _, ok := s.Infos[name]; !ok {
		return false
	}
	return false
}

func (s *Service) HasSubName(pname, name string) bool {
	if s.Len() == 0 {
		return false
	}
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if _, ok := s.Infos[pname]; !ok {
		return false
	}
	if _, ok := s.Infos[pname][name]; ok {
		return true
	}
	return false
}

func (s *Service) AddScript(subname string, script *Script) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	pname := getpname(subname)
	if pname == "" {
		return
	}
	if _, ok := s.Infos[pname]; !ok || s.Infos[pname] == nil {
		s.Infos[pname] = make(map[string]*Script)
	}
	s.Infos[pname][subname] = script
}

func (s *Service) MakeSubStruct(pname string) {
	if s.Len() == 0 {
		s.Infos[pname] = make(map[string]*Script)
		return
	}
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if _, ok := s.Infos[pname]; !ok {
		s.Infos[pname] = make(map[string]*Script)
	}
}

func (s *Service) Len() int {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return len(s.Infos)
}

func (s *Service) PnameLen(name string) int {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if v, ok := s.Infos[name]; ok {
		return len(v)
	}
	return 0
}

// 从 SS 中删除某一个subname
func (s *Service) DeleteSubname(subname string) {
	if s.Len() == 0 {
		return
	}
	s.Mu.Lock()
	defer s.Mu.Unlock()
	// 以最后一个下划线来分割出pname
	i := strings.LastIndex(subname, "_")
	if i < 0 {
		golog.Error("not found this subname :" + subname)
		return
	}
	pname := subname[:i]
	if _, ok := s.Infos[pname]; ok {
		delete(s.Infos[pname], subname)
	}
}

// 从 SS 中删除某一个pname
func (s *Service) DeletePname(pname string) {
	if s.Len() == 0 {
		return
	}
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if _, ok := s.Infos[pname]; ok {
		delete(s.Infos, pname)
	}
}

// 从 SS 中删除某一个subname
func (s *Service) GetScriptFromPnameAndSubname(pname, subname string) *Script {
	if s.Len() == 0 {
		return nil
	}
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	// 以最后一个下划线来分割出pname
	if _, ok := s.Infos[pname]; !ok {
		return nil
	}
	if _, ok := s.Infos[pname][subname]; !ok {
		return nil
	}
	return s.Infos[pname][subname]
}

type status struct {
	s     bool
	pname string
	name  string
}

func (s *Service) Copy() map[string]string {
	keys := make(map[string]string)
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	for pname := range s.Infos {
		for name := range s.Infos[pname] {
			keys[name] = pname
		}
	}
	return keys
}

// 第一次启动
func (s *Service) Start() {
	// 先无视有启动顺序的脚本
	// 先启动无需顺序的脚本
	s.Mu.RLock()
	for pname := range SS.Infos {
		for name := range SS.Infos[pname] {
			SS.Infos[pname][name].Start()
		}
	}
	s.Mu.RUnlock()

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
	Data []*ServiceStatus `json:"data"`
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
}

func (sl *StatusList) Filter(filter []string) {
	temp := make([]*ServiceStatus, 0, len(sl.Data))

	for _, s := range sl.Data {
		for _, f := range filter {
			if strings.Contains(s.Name, f) {
				temp = append(temp, s)
				break
			}
		}
	}
	sl.Data = temp
}

func All() []byte {
	if SS.Mu == nil {
		SS.Mu = &sync.RWMutex{}
	}
	statuss := &StatusList{
		Data: make([]*ServiceStatus, 0),
	}
	// ss := make([]*ServiceStatus, 0)
	for pname := range SS.Infos {
		for _, s := range SS.Infos[pname] {
			s.Status.Command = s.Command
			statuss.Data = append(statuss.Data, s.Status)
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
	if SS.Mu == nil {
		SS.Mu = &sync.RWMutex{}
	}
	statuss := &StatusList{
		Data: make([]*ServiceStatus, 0),
	}
	for _, s := range SS.Infos[pname] {
		statuss.Data = append(statuss.Data, s.Status)
	}
	statuss.Code = 200
	send, err := json.MarshalIndent(statuss, "", "\n")

	if err != nil {
		golog.Error(err)
	}
	return send
}

func ScriptName(pname, name string) []byte {
	if SS.Mu == nil {
		SS.Mu = &sync.RWMutex{}
	}

	statuss := &StatusList{
		Data: make([]*ServiceStatus, 0),
	}
	// statuss := make([]*script.ServiceStatus, 0)
	if _, pok := SS.Infos[pname]; pok {
		if s, ok := SS.Infos[pname][name]; ok {
			statuss.Data = append(statuss.Data, s.Status)
		}
	}

	send, err := json.MarshalIndent(statuss, "", "\t")
	if err != nil {
		golog.Error(err)
	}
	statuss.Code = 200
	return send
}
