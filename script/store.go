package script

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/hyahm/golog"
)

// 保存的所有脚本相关的配置
var ss *Service

func init() {
	ss = &Service{
		// 由2层组成， 一级是name  二级是pname
		Infos:   make(map[string]map[string]*Server),
		Scripts: make(map[string]*Script), // 保存脚本
		Mu:      &sync.RWMutex{},
	}
}

type Service struct {
	Infos   map[string]map[string]*Server
	Scripts map[string]*Script
	Mu      *sync.RWMutex
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

func (s *Service) HassubName(pname, name string) bool {
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

func GetServerBySubname(name string) (*Server, error) {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for pname := range ss.Infos {
		if _, ok := ss.Infos[pname][name]; ok {
			return ss.Infos[pname][name], nil
		}
	}
	return nil, ErrFoundPnameOrName
}

func GetServersByName(name string) (map[string]*Server, error) {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Infos[name]; !ok {
		return nil, ErrFoundPnameOrName
	}

	return ss.Infos[name], ErrFoundPnameOrName
}

func GetServerByNameAndSubname(pname, name string) (*Server, error) {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Infos[pname]; !ok {
		return nil, ErrFoundPnameOrName
	}
	if _, ok := ss.Infos[pname][name]; ok {
		return ss.Infos[pname][name], nil
	}

	return nil, ErrFoundPnameOrName
}

func (s *Script) KillScript() error {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Infos[s.Name]; !ok {
		return ErrFoundPnameOrName
	}
	for _, svc := range ss.Infos[s.Name] {
		svc.kill()
	}
	return nil
}

func StopAllServer() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for _, s := range ss.Scripts {
		err := s.StopScript()
		if err != nil {
			golog.Error(err)
		}
	}
}
func WaitStopAllServer() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for _, s := range ss.Scripts {
		err := s.WaitStopScript()
		if err != nil {
			golog.Error(err)
		}
	}
}

func WaitKillAllServer() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for _, s := range ss.Scripts {
		err := s.WaitKillScript()
		if err != nil {
			golog.Error(err)
		}
	}
}

// 通过pname删除script,  如果要删除server， 只能通过server.Stop() 成功后删除
func DeleteServiceByName(pname string) {
	ss.Mu.Lock()
	defer ss.Mu.Unlock()
	delete(ss.Scripts, pname)
}

//
func DeleteServiceBySubName(subname string) error {
	// 删除server
	ss.Mu.Lock()
	defer ss.Mu.Unlock()
	index := strings.LastIndex(subname, "_")
	name := subname[:index]
	if _, ok := ss.Infos[name][subname]; ok {
		delete(ss.Infos[name], subname)
		if len(ss.Infos[name]) == 0 {
			delete(ss.Scripts, name)
			delete(ss.Infos, name)
		}
		return nil
	} else {
		return ErrFoundPnameOrName
	}

}

func RestartAllServer() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for pname := range ss.Infos {
		for _, svc := range ss.Infos[pname] {
			svc.Restart()
		}
	}
}

func (s *Script) WriteToFile() error {
	// 将script 写入文件
	return AddScriptToConfigFile(s)
}

// 通过script 生成和启动服务
func (s *Script) AddScript() error {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	// 添加一个script

	// 先判断script 是否存在
	if _, ok := ss.Scripts[s.Name]; ok {
		golog.Info("已经存在此脚本")
		return ErrFoundPnameOrName
	}
	// 值挂载给ss
	ss.Scripts[s.Name] = s
	// 通过script 生成server
	s.MakeServer()

	// s.StartServer()
	return nil

}

func (s *Service) MakeSubStruct(pname string) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if _, ok := s.Infos[pname]; !ok {
		s.Infos[pname] = make(map[string]*Server)
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

// 从 ss 中删除某一个subname
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

// 从 ss 中删除某一个pname
func (s *Service) DeletePname(pname string) {
	if s.Len() == 0 {
		return
	}
	s.Mu.Lock()
	defer s.Mu.Unlock()
	delete(s.Infos, pname)
}

// 从 ss 中删除某一个subname
func (s *Service) GetScriptFromPnameAndSubname(pname, subname string) *Server {
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
	defer s.Mu.RUnlock()
	for pname := range ss.Infos {
		for name := range ss.Infos[pname] {
			ss.Infos[pname][name].Start()
		}
	}
}

func Get(pname, name string) *ServiceStatus {
	if _, ok := ss.Infos[pname]; ok {
		if sv, ok := ss.Infos[pname][name]; ok {
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
	if ss.Mu == nil {
		ss.Mu = &sync.RWMutex{}
	}
	statuss := &StatusList{
		Data: make([]*ServiceStatus, 0),
	}
	// ss := make([]*ServiceStatus, 0)
	for pname := range ss.Infos {
		for _, s := range ss.Infos[pname] {
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
	if ss.Mu == nil {
		ss.Mu = &sync.RWMutex{}
	}
	statuss := &StatusList{
		Data: make([]*ServiceStatus, 0),
	}
	for _, s := range ss.Infos[pname] {
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
	if ss.Mu == nil {
		ss.Mu = &sync.RWMutex{}
	}

	statuss := &StatusList{
		Data: make([]*ServiceStatus, 0),
	}
	// statuss := make([]*script.ServiceStatus, 0)
	if _, pok := ss.Infos[pname]; pok {
		if s, ok := ss.Infos[pname][name]; ok {
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
