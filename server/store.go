package server

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/probe"
	"github.com/hyahm/scs/status"
	"github.com/hyahm/scs/subname"
	"github.com/hyahm/scs/to"
	"github.com/hyahm/scs/utils"
)

// 保存的所有脚本相关的配置
var ss *Service

func init() {
	ss = &Service{
		// 由2层组成， 一级是name  二级是pname
		Infos:   make(map[subname.Subname]*Server),
		Scripts: make(map[string]*Script), // 保存脚本
		Mu:      &sync.RWMutex{},
	}
}

type Service struct {
	Infos   map[subname.Subname]*Server // 根据subname存放server的
	Scripts map[string]*Script          // 根据name存放脚本的
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

func HaveScript(pname string) bool {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	_, ok := ss.Scripts[pname]
	return ok
}

func NeedStart(subname string) bool {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Scripts[subname]; !ok {
		return false
	}
	if ss.Scripts[subname].Disable {
		return false
	}
	return true
}

func DelScript(pname string) {
	ss.Mu.Lock()
	defer ss.Mu.Unlock()
	delete(ss.Scripts, pname)
}

func RemoveScript(pname string) error {
	// del := make(chan bool)
	ss.Mu.Lock()
	defer ss.Mu.Unlock()

	if _, ok := ss.Scripts[pname]; ok {
		// go func() {
		// wg := &sync.WaitGroup{}
		replicate := ss.Scripts[pname].Replicate
		if replicate == 0 {
			replicate = 1
		}

		for i := 0; i < replicate; i++ {
			subname := subname.NewSubname(pname, i)
			ss.Infos[subname].Remove()
		}

	} else {
		return errors.New("not found this pname:" + pname)
	}
	return nil
}

func AddAndStartServer(script *Script) {
	ss.Scripts[script.Name] = script
	ss.Scripts[script.Name].MakeServer()
	replicate := script.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(script.Name, i)
		ss.Infos[subname].Start()
	}
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

func HassubName(name subname.Subname) bool {

	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Infos[name]; ok {
		return true
	}
	return false
}

func GetServerBySubname(subname subname.Subname) (*Server, bool) {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	v, ok := ss.Infos[subname]
	return v, ok
}

func GetServersByName(name string) (map[subname.Subname]*Server, bool) {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Scripts[name]; !ok {
		return nil, false
	}
	servers := make(map[subname.Subname]*Server)
	replicate := ss.Scripts[name].Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(name, i)
		servers[subname] = ss.Infos[subname]
	}

	return servers, true
}

func GetServerByNameAndSubname(name string, subname subname.Subname) (*Server, bool) {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Scripts[name]; !ok {
		return nil, false
	}
	if _, ok := ss.Infos[subname]; ok {
		return ss.Infos[subname], true
	}

	return nil, false
}

func StopAllServer() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for _, s := range ss.Scripts {
		err := StopScript(s)
		if err != nil {
			golog.Error(err)
		}
	}
}
func WaitStopAllServer() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for _, s := range ss.Scripts {
		WaitStopScript(s)
	}
}

func WaitKillAllServer() {
	// ss.ScriptLocker.RLock()
	// defer ss.ScriptLocker.RUnlock()
	for _, s := range ss.Scripts {
		WaitKillScript(s)
	}
}

// 只删除ss.infos 里面的额
func DeleteServiceBySubName(subname subname.Subname) error {
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
	return errors.New("")

}

func GetAllServer() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for _, svc := range ss.Infos {
		svc.Start()
	}
}

func UpdateAndRestartAllServer() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for _, s := range ss.Scripts {
		go UpdateAndRestartScript(s)
	}
}

func RestartAllServer() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for _, svc := range ss.Infos {
		go svc.Restart()
	}
}

func Len() int {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	return len(ss.Scripts)
}

// 从 ss 中删除某一个subname
func DeleteSubname(subname subname.Subname) {

	ss.Mu.Lock()
	defer ss.Mu.Unlock()
	// 以最后一个下划线来分割出pname
	if _, ok := ss.Infos[subname]; ok {
		delete(ss.Infos, subname)
	}
}

type StatusList struct {
	Data []*status.ServiceStatus `json:"data"`
	Code int                     `json:"code"`
	Msg  string                  `json:"msg"`
}

func (sl *StatusList) Filter(filter []string) {
	temp := make([]*status.ServiceStatus, 0, len(sl.Data))

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
		Data: make([]*status.ServiceStatus, 0),
	}
	// ss := make([]*ServiceStatus, 0)
	for subname := range ss.Infos {
		pname := subname.GetName()
		if _, ok := ss.Scripts[pname]; !ok {
			golog.Debug("not found name: " + pname)
			continue
		}

		status := &status.ServiceStatus{
			PName:        pname,
			Name:         subname.String(),
			Command:      ss.Infos[subname].Command,
			Always:       ss.Scripts[pname].Always,
			Version:      ss.Infos[subname].Version,
			Status:       ss.Infos[subname].Status.Status,
			CanNotStop:   ss.Infos[subname].Status.CanNotStop,
			Path:         ss.Infos[subname].Script.Dir,
			Start:        ss.Infos[subname].Status.Start,
			RestartCount: ss.Infos[subname].Status.RestartCount,
			Up:           ss.Infos[subname].Status.Up,
			Disable:      ss.Scripts[pname].Disable,
		}
		if ss.Infos[subname].Cmd != nil && ss.Infos[subname].Cmd.Process != nil {
			status.Pid = ss.Infos[subname].Status.Pid
			status.Cpu, status.Mem, _ = probe.GetProcessInfo(int32(ss.Infos[subname].Status.Pid))

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
		Data: make([]*status.ServiceStatus, 0),
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
		subname := subname.NewSubname(pname, i)
		status := &status.ServiceStatus{
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
		if ss.Infos[subname].Cmd != nil && ss.Infos[subname].Cmd.Process != nil {
			status.Pid = ss.Infos[subname].Cmd.Process.Pid
			status.Cpu, status.Mem, _ = probe.GetProcessInfo(int32(ss.Infos[subname].Cmd.Process.Pid))

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

func ScriptName(pname string, subname subname.Subname) []byte {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	statuss := &StatusList{
		Data: make([]*status.ServiceStatus, 0),
	}
	if _, ok := ss.Scripts[pname]; !ok {
		statuss.Msg = "not found " + pname
		send, err := json.MarshalIndent(statuss, "", "\n")

		if err != nil {
			golog.Error(err)
		}
		return send
	}

	status := &status.ServiceStatus{
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
	if ss.Infos[subname].Cmd != nil && ss.Infos[subname].Cmd.Process != nil {
		status.Pid = ss.Infos[subname].Cmd.Process.Pid
		status.Cpu, status.Mem, _ = probe.GetProcessInfo(int32(ss.Infos[subname].Cmd.Process.Pid))

	}
	statuss.Data = append(statuss.Data, status)
	send, err := json.MarshalIndent(statuss, "", "\t")
	if err != nil {
		golog.Error(err)
	}
	statuss.Code = 200
	return send
}

func NeedStop(s *Script) bool {
	// 更新server
	// 判断值是否相等
	if s.Dir != ss.Scripts[s.Name].Dir ||
		s.Command != ss.Scripts[s.Name].Command ||
		s.Replicate != ss.Scripts[s.Name].Replicate ||
		s.Always != ss.Scripts[s.Name].Always ||
		s.DisableAlert != ss.Scripts[s.Name].DisableAlert ||
		!utils.CompareMap(s.Env, ss.Scripts[s.Name].Env) ||
		s.Port != ss.Scripts[s.Name].Port ||
		s.Version != ss.Scripts[s.Name].Version ||
		s.Disable != ss.Scripts[s.Name].Disable ||
		s.Update != ss.Scripts[s.Name].Update ||
		s.DeleteWhenExit != ss.Scripts[s.Name].DeleteWhenExit ||
		!s.Cron.IsEqual(ss.Scripts[s.Name].Cron) ||
		!IsEqual(s.Name, s.AT) {
		// 如果值有变动， 那么需要重新server
		// 先同步停止之前的server， 然后启动新的server
		// server 是单独的， 在通知后需要同步更新server
		return true
	}
	return false
}

// 启动方法， 异步执行
func StartServer(s *Script) {
	ss.Mu.Lock()
	ss.Mu.Unlock()
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(s.Name, i)
		ss.Infos[subname].Start()
	}
}

func KillScript(s *Script) {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(s.Name, i)
		ss.Infos[subname].Kill()
	}
}

// 通过script 生成和启动服务
func AddScript(s *Script) {
	// 通过script 生成server
	s.MakeServer()
	StartServer(s)
}

// 异步重启
func RestartScript(s *Script) error {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	// 禁用 script 所在的所有server
	if _, ok := ss.Scripts[s.Name]; !ok {
		return errors.New("")
	}
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(s.Name, i)
		go ss.Infos[subname].Restart()
	}
	return nil
}

// 同步杀掉
func WaitKillScript(s *Script) {
	// ss.ServerLocker.RLock()
	// defer ss.ServerLocker.RUnlock()
	// 禁用 script 所在的所有server
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	// 禁用 script 所在的所有server
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(s.Name, i)
		ss.Infos[subname].Kill()
	}
}

// 同步停止
func WaitStopScript(s *Script) {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	// 禁用 script 所在的所有server
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(s.Name, i)
		ss.Infos[subname].Stop()
	}
}

// 异步执行停止脚本
func StopScript(s *Script) error {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Scripts[s.Name]; !ok {
		return errors.New("")
	}
	// 禁用 script 所在的所有server
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(s.Name, i)
		go ss.Infos[subname].Stop()
	}
	return nil
}

func GetScriptByPname(name string) (*Script, bool) {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	v, ok := ss.Scripts[name]
	return v, ok

}

// 返回成功还是失败
func UpdateAndRestartScript(s *Script) bool {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Scripts[s.Name]; !ok {
		return false
	}
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(s.Name, i)
		go ss.Infos[subname].UpdateAndRestart()
	}
	return true
}

func EnableScript(s *Script) bool {
	ss.Mu.Lock()
	defer ss.Mu.Unlock()
	// 禁用 script 所在的所有server
	if _, ok := ss.Scripts[s.Name]; !ok {
		return false
	}
	ss.Scripts[s.Name].Disable = true
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(s.Name, i)
		go ss.Infos[subname].Start()
	}
	return true
}

func DisableScript(s *Script) bool {
	ss.Mu.Lock()
	defer ss.Mu.Unlock()
	// 禁用 script 所在的所有server
	if _, ok := ss.Scripts[s.Name]; !ok {
		return false
	}
	ss.Scripts[s.Name].Disable = true
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := subname.NewSubname(s.Name, i)
		go ss.Infos[subname].Stop()
	}
	return true
}

func AddInfo(name subname.Subname, svc *Server) {

}

// 比较新的与之前的是否相等， 调用者必须是新的
func IsEqual(pname string, at *to.AlertTo) bool {
	if at == nil && ss.Scripts[pname].AT == nil {
		return true
	}
	if (at == nil && ss.Scripts[pname].AT != nil) || (at != nil && ss.Scripts[pname].AT == nil) {
		return false
	}
	if !utils.CompareSlice(at.Email, ss.Scripts[pname].AT.Email) ||
		!utils.CompareSlice(at.Rocket, ss.Scripts[pname].AT.Rocket) ||
		!utils.CompareSlice(at.Telegram, ss.Scripts[pname].AT.Telegram) ||
		!utils.CompareSlice(at.WeiXin, ss.Scripts[pname].AT.WeiXin) {
		return false
	}
	return true
}

func StartAllServer() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for _, v := range ss.Infos {
		v.Start()
	}
}
