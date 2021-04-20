package scs

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
		Infos:        make(map[string]map[string]*Server),
		Scripts:      make(map[string]*Script), // 保存脚本
		ServerLocker: &sync.RWMutex{},
		ScriptLocker: &sync.RWMutex{},
	}
}

type Service struct {
	Infos        map[string]map[string]*Server
	Scripts      map[string]*Script
	ScriptLocker *sync.RWMutex
	ServerLocker *sync.RWMutex
}

func (s *Service) HasName(name string) bool {
	if s.Len() == 0 {
		return false
	}
	s.ScriptLocker.RLock()
	defer s.ScriptLocker.RUnlock()
	if _, ok := s.Scripts[name]; !ok {
		return false
	}
	return false
}

func (s *Service) HassubName(pname, name string) bool {
	if s.Len() == 0 {
		return false
	}
	s.ServerLocker.RLock()
	defer s.ServerLocker.RUnlock()
	if _, ok := s.Infos[pname]; !ok {
		return false
	}
	if _, ok := s.Infos[pname][name]; ok {
		return true
	}
	return false
}

func GetServerBySubname(name string) (*Server, error) {
	ss.ServerLocker.RLock()
	defer ss.ServerLocker.RUnlock()
	for pname := range ss.Infos {
		if _, ok := ss.Infos[pname][name]; ok {
			return ss.Infos[pname][name], nil
		}
	}
	return nil, ErrFoundPnameOrName
}

func GetServersByName(name string) (map[string]*Server, error) {
	ss.ServerLocker.RLock()
	defer ss.ServerLocker.RUnlock()
	if _, ok := ss.Infos[name]; !ok {
		return nil, ErrFoundPnameOrName
	}

	return ss.Infos[name], ErrFoundPnameOrName
}

func GetServerByNameAndSubname(pname, name string) (*Server, error) {
	ss.ServerLocker.RLock()
	defer ss.ServerLocker.RUnlock()
	if _, ok := ss.Infos[pname]; !ok {
		return nil, ErrFoundPnameOrName
	}
	if _, ok := ss.Infos[pname][name]; ok {
		return ss.Infos[pname][name], nil
	}

	return nil, ErrFoundPnameOrName
}

func (s *Script) KillScript() error {
	ss.ServerLocker.RLock()
	defer ss.ServerLocker.RUnlock()
	if _, ok := ss.Infos[s.Name]; !ok {
		return ErrFoundPnameOrName
	}
	for _, svc := range ss.Infos[s.Name] {
		svc.kill()
	}
	return nil
}

func StopAllServer() {
	ss.ScriptLocker.RLock()
	defer ss.ScriptLocker.RUnlock()
	for _, s := range ss.Scripts {
		err := s.StopScript()
		if err != nil {
			golog.Error(err)
		}
	}
}
func WaitStopAllServer() {
	ss.ScriptLocker.RLock()
	defer ss.ScriptLocker.RUnlock()
	for _, s := range ss.Scripts {
		err := s.WaitStopScript()
		if err != nil {
			golog.Error(err)
		}
	}
}

func WaitKillAllServer() {
	ss.ScriptLocker.RLock()
	defer ss.ScriptLocker.RUnlock()
	for _, s := range ss.Scripts {
		err := s.WaitKillScript()
		if err != nil {
			golog.Error(err)
		}
	}
}

// 通过pname删除script,  如果要删除server， 只能通过server.Stop() 成功后删除
// func DeleteServiceByName(pname string) {
// 	ss.ScriptLocker.Lock()
// 	defer ss.ScriptLocker.Unlock()
// 	delete(ss.Scripts, pname)
// }

// 只删除ss.infos 里面的额
func DeleteServiceBySubName(subname string) error {
	// 删除server
	ss.ServerLocker.Lock()
	defer ss.ServerLocker.Unlock()
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
	ss.ServerLocker.RLock()
	defer ss.ServerLocker.RUnlock()
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
	ss.ScriptLocker.RLock()
	defer ss.ScriptLocker.RUnlock()
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
	s.StartServer()
	return nil

}

func (s *Service) MakeSubStruct(pname string) {
	s.ServerLocker.RLock()
	defer s.ServerLocker.RUnlock()
	if _, ok := s.Infos[pname]; !ok {
		s.Infos[pname] = make(map[string]*Server)
	}
}

func (s *Service) Len() int {
	s.ServerLocker.RLock()
	defer s.ServerLocker.RUnlock()
	return len(s.Scripts)
}

func (s *Service) PnameLen(name string) int {
	s.ServerLocker.RLock()
	defer s.ServerLocker.RUnlock()
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
	s.ServerLocker.Lock()
	defer s.ServerLocker.Unlock()
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
// func (s *Service) DeletePname(pname string) {
// 	if s.Len() == 0 {
// 		return
// 	}
// 	s.ServerLocker.Lock()
// 	defer s.ServerLocker.Unlock()
// 	delete(s.Infos, pname)
// }

// 从 ss 中删除某一个subname
func GetScriptFromPnameAndSubname(pname, subname string) *Server {
	if ss.Len() == 0 {
		return nil
	}
	ss.ServerLocker.RLock()
	defer ss.ServerLocker.RUnlock()
	// 以最后一个下划线来分割出pname
	if _, ok := ss.Infos[pname]; !ok {
		return nil
	}
	if _, ok := ss.Infos[pname][subname]; !ok {
		return nil
	}
	return ss.Infos[pname][subname]
}

// func Copy() map[string]string {
// 	keys := make(map[string]string)
// 	ss.ServerLocker.RLock()
// 	defer ss.ServerLocker.RUnlock()
// 	for pname := range ss.Infos {
// 		for name := range ss.Infos[pname] {
// 			keys[name] = pname
// 		}
// 	}
// 	return keys
// }

// 第一次启动
// func (s *Service) Start() {
// 	// 先无视有启动顺序的脚本
// 	// 先启动无需顺序的脚本
// 	ss.ServerLocker.RLock()
// 	defer ss.ServerLocker.RUnlock()
// 	for pname := range ss.Infos {
// 		for name := range ss.Infos[pname] {
// 			ss.Infos[pname][name].Start()
// 		}
// 	}
// }

func Get(pname, name string) *ServiceStatus {
	ss.ServerLocker.RLock()
	defer ss.ServerLocker.RUnlock()
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
	ss.ServerLocker.Lock()
	ss.ScriptLocker.RLock()
	defer ss.ScriptLocker.RUnlock()
	defer ss.ServerLocker.Unlock()
	statuss := &StatusList{
		Data: make([]*ServiceStatus, 0),
	}
	// ss := make([]*ServiceStatus, 0)
	for pname := range ss.Infos {
		for name := range ss.Infos[pname] {
			ss.Infos[pname][name].Status.Cpu, ss.Infos[pname][name].Status.Mem, _ = GetProcessInfo(int32(ss.Infos[pname][name].cmd.Process.Pid))
			ss.Infos[pname][name].Status.Command = ss.Scripts[pname].Command
			ss.Infos[pname][name].Status.PName = ss.Scripts[pname].Name
			ss.Infos[pname][name].Status.Always = ss.Scripts[pname].Always
			ss.Infos[pname][name].Status.Disable = ss.Scripts[pname].Disable
			statuss.Data = append(statuss.Data, ss.Infos[pname][name].Status)
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
	ss.ServerLocker.RLock()
	defer ss.ServerLocker.RUnlock()
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
	ss.ServerLocker.RLock()
	defer ss.ServerLocker.RUnlock()
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
