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
		Infos:   make(map[Subname]*Server),
		Scripts: make(map[string]*Script), // 保存脚本
		Mu:      &sync.RWMutex{},
	}
}

type Service struct {
	Infos   map[Subname]*Server // 根据subname存放server的
	Scripts map[string]*Script  // 根据name存放脚本的
	Mu      *sync.RWMutex
}

func GetServers() []byte {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	send, err := json.Marshal(ss.Infos)
	if err != nil {
		golog.Error(err)
	}
	return send
}

func GetScripts() []byte {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	send, err := json.Marshal(ss.Scripts)
	if err != nil {
		golog.Error(err)
	}
	return send
}

func HasName(name string) bool {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Scripts[name]; ok {
		return true
	}
	return false
}

func HassubName(name Subname) bool {

	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Infos[name]; ok {
		return true
	}
	return false
}

func GetServerBySubname(subname Subname) (*Server, error) {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Infos[subname]; ok {
		return ss.Infos[subname], nil
	}
	return nil, ErrFoundName
}

func GetServersByName(name string) (map[Subname]*Server, error) {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Scripts[name]; !ok {
		return nil, ErrFoundPname
	}
	servers := make(map[Subname]*Server)
	replicate := ss.Scripts[name].Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := NewSubname(name, i)
		servers[subname] = ss.Infos[subname]
	}

	return servers, nil
}

func GetServerByNameAndSubname(name string, subname Subname) (*Server, error) {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Scripts[name]; !ok {
		return nil, ErrFoundPname
	}
	if _, ok := ss.Infos[subname]; ok {
		return ss.Infos[subname], nil
	}

	return nil, ErrFoundName
}

func (s *Script) KillScript() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()

	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := NewSubname(s.Name, i)
		ss.Infos[subname].Kill()
	}
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
		s.WaitStopScript()
	}
}

func WaitKillAllServer() {
	// ss.ScriptLocker.RLock()
	// defer ss.ScriptLocker.RUnlock()
	for _, s := range ss.Scripts {
		s.WaitKillScript()
	}
}

// 通过pname删除script,  如果要删除server， 只能通过server.Stop() 成功后删除
// func DeleteServiceByName(pname string) {
// 	ss.ScriptLocker.Lock()
// 	defer ss.ScriptLocker.Unlock()
// 	delete(ss.Scripts, pname)
// }

// 只删除ss.infos 里面的额
func DeleteServiceBySubName(subname Subname) error {
	// 删除server
	ss.Mu.Lock()
	defer ss.Mu.Unlock()
	if _, ok := ss.Infos[subname]; ok {
		delete(ss.Infos, subname)
		// 同时script 里面也要删除
		name := subname.GetName()
		if _, ok := ss.Scripts[name]; ok {
			ss.Scripts[name].Replicate--
			// 开发中， replicate =0 或 1 其实都是1 的意思， 所以减一后 <= 0 的其实就是都删除干净了的意思
			if ss.Scripts[name].Replicate < 1 {
				delete(ss.Scripts, name)
				delete(ss.Infos, subname)
			}
		}
		return nil
	}
	return ErrFoundName

}

func RestartAllServer() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for _, svc := range ss.Infos {
		go svc.Restart()
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
	// 先判断script 是否存在, 不存在的话应该是修改的
	// if _, ok := ss.Scripts[s.Name]; ok {
	// 	golog.Info("已经存在此脚本")
	// 	return ErrFoundPnameOrName
	// }
	// 值挂载给ss
	ss.Scripts[s.Name] = s
	// 通过script 生成server
	s.MakeServer()
	s.StartServer()
	return nil

}

func Len() int {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	return len(ss.Scripts)
}

// 从 ss 中删除某一个subname
func DeleteSubname(subname Subname) {

	ss.Mu.Lock()
	defer ss.Mu.Unlock()
	// 以最后一个下划线来分割出pname
	if _, ok := ss.Infos[subname]; ok {
		delete(ss.Infos, subname)
	}
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

// 获取所有服务的状态
func All() []byte {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	statuss := &StatusList{
		Data: make([]*ServiceStatus, 0),
	}
	// ss := make([]*ServiceStatus, 0)
	for subname := range ss.Infos {

		if _, ok := ss.Scripts[subname.GetName()]; !ok {
			continue
		}

		status := &ServiceStatus{
			PName:        subname.GetName(),
			Name:         subname.String(),
			Command:      ss.Infos[subname].Command,
			Always:       ss.Scripts[subname.GetName()].Always,
			Version:      ss.Infos[subname].Version,
			Status:       ss.Infos[subname].Status.Status,
			CanNotStop:   ss.Infos[subname].Status.CanNotStop,
			Path:         ss.Infos[subname].Script.Dir,
			Start:        ss.Infos[subname].Status.Start,
			RestartCount: ss.Infos[subname].Status.RestartCount,
			Up:           ss.Infos[subname].Status.Up,
			Disable:      ss.Scripts[subname.GetName()].Disable,
		}
		if ss.Infos[subname].cmd != nil && ss.Infos[subname].cmd.Process != nil {
			status.Pid = ss.Infos[subname].Status.Pid
			status.Cpu, status.Mem, _ = GetProcessInfo(int32(ss.Infos[subname].Status.Pid))

		}
		statuss.Data = append(statuss.Data, status)
	}
	statuss.Code = 200
	send, err := json.MarshalIndent(statuss, "", "\t")
	if err != nil {
		golog.Error(err)
	}
	return send
}

func ScriptPname(pname string) []byte {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	statuss := &StatusList{
		Data: make([]*ServiceStatus, 0),
	}
	if _, ok := ss.Scripts[pname]; !ok {
		statuss.Msg = "not found " + pname
		send, err := json.MarshalIndent(statuss, "", "\n")

		if err != nil {
			golog.Error(err)
		}
		return send
	}
	replicate := ss.Scripts[pname].Replicate
	if replicate == 0 {
		replicate = 1
	}

	for i := 0; i < replicate; i++ {
		subname := NewSubname(pname, i)
		status := &ServiceStatus{
			PName:        subname.GetName(),
			Name:         subname.String(),
			Command:      ss.Infos[subname].Status.Command,
			Always:       ss.Scripts[subname.GetName()].Always,
			Version:      ss.Infos[subname].Status.Version,
			CanNotStop:   ss.Infos[subname].Status.CanNotStop,
			Path:         ss.Infos[subname].Status.Path,
			Status:       ss.Infos[subname].Status.Status,
			RestartCount: ss.Infos[subname].Status.RestartCount,
			Up:           ss.Infos[subname].Status.Up,
			Disable:      ss.Scripts[subname.GetName()].Disable,
		}
		if ss.Infos[subname].cmd != nil && ss.Infos[subname].cmd.Process != nil {
			status.Pid = ss.Infos[subname].cmd.Process.Pid
			status.Cpu, status.Mem, _ = GetProcessInfo(int32(ss.Infos[subname].cmd.Process.Pid))

		}
		statuss.Data = append(statuss.Data, status)
	}

	statuss.Code = 200
	send, err := json.MarshalIndent(statuss, "", "\n")

	if err != nil {
		golog.Error(err)
	}
	return send
}

func ScriptName(pname string, subname Subname) []byte {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	statuss := &StatusList{
		Data: make([]*ServiceStatus, 0),
	}
	if _, ok := ss.Scripts[pname]; !ok {
		statuss.Msg = "not found " + pname
		send, err := json.MarshalIndent(statuss, "", "\n")

		if err != nil {
			golog.Error(err)
		}
		return send
	}

	status := &ServiceStatus{
		PName:        subname.GetName(),
		Name:         subname.String(),
		Command:      ss.Infos[subname].Status.Command,
		Always:       ss.Scripts[subname.GetName()].Always,
		Version:      ss.Infos[subname].Status.Version,
		CanNotStop:   ss.Infos[subname].Status.CanNotStop,
		Path:         ss.Infos[subname].Status.Path,
		Status:       ss.Infos[subname].Status.Status,
		RestartCount: ss.Infos[subname].Status.RestartCount,
		Up:           ss.Infos[subname].Status.Up,
		Disable:      ss.Scripts[subname.GetName()].Disable,
	}
	if ss.Infos[subname].cmd != nil && ss.Infos[subname].cmd.Process != nil {
		status.Pid = ss.Infos[subname].cmd.Process.Pid
		status.Cpu, status.Mem, _ = GetProcessInfo(int32(ss.Infos[subname].cmd.Process.Pid))

	}
	statuss.Data = append(statuss.Data, status)
	send, err := json.MarshalIndent(statuss, "", "\t")
	if err != nil {
		golog.Error(err)
	}
	statuss.Code = 200
	return send
}
