package script

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/server/alert"
	"github.com/hyahm/scs/server/script/status"
	"gopkg.in/yaml.v2"
)

// 保存的所有脚本相关的配置
var ss *Service

type Service struct {
	Infos   map[string]map[string]*Server
	Scripts []*Script
	Mu      *sync.RWMutex
}

func RunServer(scripts []*Script) error {
	// 将 config 的数据加载到ss 中
	if ss == nil {
		ss = &Service{
			// 由2层组成， 一级是name  二级是pname
			Infos:   make(map[string]map[string]*Server),
			Scripts: scripts,
			Mu:      &sync.RWMutex{},
		}
	} else {
		// 对比不同值
	}
	for index := range ss.Scripts {
		// 如果名字为空， 或者 command 为空  那么就跳过
		if strings.Trim(ss.Scripts[index].Name, " ") == "" ||
			strings.Trim(ss.Scripts[index].Command, " ") == "" {
			continue
		}

		// 如果 ss 的key
		if _, ok := ss.Infos[ss.Scripts[index].Name]; !ok {
			ss.Infos[ss.Scripts[index].Name] = make(map[string]*Server)
		}
		// 第一次启动的时候
		fill(index)
	}
	return nil
}

func fill(index int) {
	// 加载环境变量
	baseEnv := make(map[string]string)

	// 填充系统环境变量到
	pathEnvName := "PATH"
	for _, v := range os.Environ() {
		kv := strings.Split(v, "=")
		if strings.ToUpper(kv[0]) == pathEnvName {
			pathEnvName = kv[0]
		}
		baseEnv[kv[0]] = kv[1]
	}
	golog.Info(pathEnvName)
	for k, v := range ss.Scripts[index].Env {
		// path 环境单独处理， 可以多个值， 其他环境变量多个值请以此写完
		if strings.ToLower(k) == strings.ToLower(pathEnvName) {
			if runtime.GOOS == "windows" {
				baseEnv[pathEnvName] = baseEnv[pathEnvName] + ";" + v
			} else {
				golog.Info(pathEnvName)
				baseEnv[pathEnvName] = baseEnv[pathEnvName] + ":" + v
			}
		} else {
			baseEnv[k] = v
		}
	}

	baseEnv["TOKEN"] = global.Token
	baseEnv["PNAME"] = ss.Scripts[index].Name
	// 结束当前索引值环境变量
	replica := ss.Scripts[index].Replicate
	if replica < 1 {
		replica = 1
	}
	// 通过config.sc 生成 server
	for i := 0; i < replica; i++ {

		// 根据server副本数提取子名称
		subname := fmt.Sprintf("%s_%d", ss.Scripts[index].Name, i)
		// if reload {
		// 	// 判断下值是否发生了变化， 发生了变化才重启， 否则不做任何操作

		// 	// 如果是加载配置文件， 那么删除已经有的
		// 	golog.Info("delete subname")
		// 	DelDelScript(subname)
		// }
		baseEnv["NAME"] = subname
		baseEnv["PORT"] = strconv.Itoa(ss.Scripts[index].Port + i)

		if ss.HassubName(ss.Scripts[index].Name, subname) {
			// 如果存在键值就修改，
			golog.Info("update")
			update(index, ss.Scripts[index].Port+i,
				subname, ss.Scripts[index].Command, baseEnv)
		} else {
			golog.Info("add")
			// 新增
			ss.MakeSubStruct(ss.Scripts[index].Name)
			add(index, ss.Scripts[index].Port+i, subname,
				ss.Scripts[index].Command, baseEnv)
		}

	}
	// 删除多余的副本
	go func() {
		pname := ss.Scripts[index].Name

		replicate := ss.Scripts[index].Replicate
		if replicate < 1 {
			replicate = 1
		}
		l := ss.Len()
		if l > 0 && l > replicate {
			for i := l - 1; i >= replicate; i-- {
				subname := fmt.Sprintf("%s_%d", pname, i)
				// if reload {
				// 	// 如果是加载配置文件， 那么删除已经有的
				// 	DelDelScript(subname)
				// }
				ss.GetScriptFromPnameAndSubname(pname, subname).Remove()
				// ss.Infos[pname][subname].Stop()
				// delete(ss.Infos[pname], subname)
			}
		}
	}()
}

func add(index, port int, subname, command string, baseEnv map[string]string) {
	svc := &Server{
		Name:      ss.Scripts[index].Name,
		LookPath:  ss.Scripts[index].LookPath,
		Command:   command,
		Env:       baseEnv,
		Dir:       ss.Scripts[index].Dir,
		Replicate: ss.Scripts[index].Replicate,
		Log:       make(map[string][]string),
		LogLocker: &sync.RWMutex{},
		SubName:   subname,
		Status: &status.ServiceStatus{
			Name:    subname,
			PName:   ss.Scripts[index].Name,
			Status:  status.STOP,
			Path:    ss.Scripts[index].Dir,
			Version: getVersion(ss.Scripts[index].Version),
			Disable: ss.Scripts[index].Disable,
		},
		DeleteWhenExit:     ss.Scripts[index].DeleteWhenExit,
		Update:             ss.Scripts[index].Update,
		DisableAlert:       ss.Scripts[index].DisableAlert,
		ContinuityInterval: ss.Scripts[index].ContinuityInterval,
		Always:             ss.Scripts[index].Always,
		Disable:            ss.Scripts[index].Disable,
		AI:                 &alert.AlertInfo{},
		Port:               port,
		Exit:               make(chan int, 2),
		CancelSigle:        make(chan bool, 2),
		AT:                 ss.Scripts[index].AT,
	}
	// 生成对应的文件类型
	svc.Log["log"] = make([]string, 0, global.LogCount)
	svc.Log["lookPath"] = make([]string, 0, global.LogCount)
	svc.Log["update"] = make([]string, 0, global.LogCount)
	if ss.Scripts[index].Cron != nil {
		start, err := time.ParseInLocation("2006-01-02 15:04:05", ss.Scripts[index].Cron.Start, time.Local)
		if err != nil {
			start = time.Time{}
		}
		svc.Cron = &Cron{
			Start:   start,
			IsMonth: ss.Scripts[index].Cron.IsMonth,
			Loop:    ss.Scripts[index].Cron.Loop,
		}
	}
	ss.AddScript(subname, svc)

	if !svc.Disable {
		svc.Start()
	}
}

func (svc *Service) Delete(pname, name string) {
	svc.Mu.Lock()
	defer svc.Mu.Unlock()
	if _, ok := ss.Infos[pname]; ok {
		delete(svc.Infos[pname], name)
		if len(svc.Infos[pname]) == 0 {
			tmp := make([]*Script, 0)
			for i := range ss.Scripts {
				if ss.Scripts[i].Name == pname {
					tmp = append(tmp, ss.Scripts[:i]...)
					tmp = append(tmp, ss.Scripts[i+1:]...)
					ss.Scripts = tmp
					break
				}
			}
			// 并且删除配置文件的
		}
	}

}

func update(index, port int, subname, command string, baseEnv map[string]string) {
	scriptInfo := ss.GetScriptFromPnameAndSubname(ss.Scripts[index].Name, subname)

	scriptInfo.Env = baseEnv
	scriptInfo.LookPath = ss.Scripts[index].LookPath
	if ss.Scripts[index].Cron != nil {
		start, err := time.ParseInLocation("2006-01-02 15:04:05", ss.Scripts[index].Cron.Start, time.Local)
		if err != nil {
			start = time.Time{}
		}
		scriptInfo.Cron = &Cron{
			Start:   start,
			IsMonth: ss.Scripts[index].Cron.IsMonth,
			Loop:    ss.Scripts[index].Cron.Loop,
		}
	}

	scriptInfo.Command = command
	scriptInfo.DeleteWhenExit = ss.Scripts[index].DeleteWhenExit
	scriptInfo.Update = ss.Scripts[index].Update
	scriptInfo.Dir = ss.Scripts[index].Dir
	scriptInfo.Replicate = ss.Scripts[index].Replicate
	scriptInfo.Log = make(map[string][]string)
	ss.Infos[ss.Scripts[index].Name][subname].LogLocker = &sync.RWMutex{}
	ss.Infos[ss.Scripts[index].Name][subname].Log["log"] = make([]string, 0, global.LogCount)
	ss.Infos[ss.Scripts[index].Name][subname].Log["lookPath"] = make([]string, 0, global.LogCount)
	ss.Infos[ss.Scripts[index].Name][subname].Log["update"] = make([]string, 0, global.LogCount)
	ss.Infos[ss.Scripts[index].Name][subname].DisableAlert = ss.Scripts[index].DisableAlert
	ss.Infos[ss.Scripts[index].Name][subname].Always = ss.Scripts[index].Always
	ss.Infos[ss.Scripts[index].Name][subname].ContinuityInterval = ss.Scripts[index].ContinuityInterval
	ss.Infos[ss.Scripts[index].Name][subname].Port = ss.Scripts[index].Port + index
	ss.Infos[ss.Scripts[index].Name][subname].AT = ss.Scripts[index].AT
	ss.Infos[ss.Scripts[index].Name][subname].Disable = ss.Scripts[index].Disable
	ss.Infos[ss.Scripts[index].Name][subname].Status.Version = getVersion(ss.Scripts[index].Version)
	ss.Infos[ss.Scripts[index].Name][subname].Status.Disable = ss.Scripts[index].Disable
	// 更新的时候

	if ss.Infos[ss.Scripts[index].Name][subname].Status.Status == status.STOP {
		// 如果是停止的name就启动
		if strings.Trim(ss.Scripts[index].Command, " ") != "" && strings.Trim(ss.Scripts[index].Name, " ") != "" && !ss.Scripts[index].Disable {
			ss.Infos[ss.Scripts[index].Name][subname].Start()
		}
	}
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

func (s *Service) AddScript(subname string, server *Server) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	pname := getpname(subname)
	if pname == "" {
		return
	}
	if _, ok := s.Infos[pname]; !ok || s.Infos[pname] == nil {
		s.Infos[pname] = make(map[string]*Server)
	}
	s.Infos[pname][subname] = server
}

func (s *Service) MakeSubStruct(pname string) {
	if s.Len() == 0 {
		s.Infos[pname] = make(map[string]*Server)
		return
	}
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
	if _, ok := s.Infos[pname]; ok {
		delete(s.Infos, pname)
	}
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
	return ss.Infos[pname][subname]
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
func Start() {
	// 先无视有启动顺序的脚本
	// 先启动无需顺序的脚本
	ss.Mu.RLock()
	for pname := range ss.Infos {
		for name := range ss.Infos[pname] {
			ss.Infos[pname][name].Start()
		}
	}
	ss.Mu.RUnlock()

}

func StopPname(pname string) bool {
	if _, ok := ss.Infos[pname]; ok {
		for name := range ss.Infos[pname] {
			ss.Infos[pname][name].Disable = true
			ss.Infos[pname][name].Status.Disable = true
			go ss.Infos[pname][name].Stop()
		}
		return true
	}
	return false
}

func StopServer(pname, name string) bool {
	if _, ok := ss.Infos[pname]; ok {
		if _, nok := ss.Infos[pname][name]; nok {
			go ss.Infos[pname][name].Stop()
		}
		return true
	}
	return false
}

func StopAllServer() {
	for pname := range ss.Infos {
		for name := range ss.Infos[pname] {
			ss.Infos[pname][name].Disable = true
			ss.Infos[pname][name].Status.Disable = true
			go ss.Infos[pname][name].Stop()
		}
	}

}

func Get(pname, name string) *status.ServiceStatus {
	if _, ok := ss.Infos[pname]; ok {
		if sv, ok := ss.Infos[pname][name]; ok {
			return sv.Status
		}

	}
	return nil
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

func All() []byte {
	if ss.Mu == nil {
		ss.Mu = &sync.RWMutex{}
	}
	statuss := &StatusList{
		Data: make([]*status.ServiceStatus, 0),
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
		Data: make([]*status.ServiceStatus, 0),
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
		Data: make([]*status.ServiceStatus, 0),
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

func AddScript(s *Script) error {

	golog.Infof("%+v", s)
	// 添加到配置文件
	for i, v := range ss.Scripts {
		if v.Name == s.Name {
			// 修改
			// c.updateConfig(s, i)
			fill(i)
			b, err := yaml.Marshal(s)
			if err != nil {
				return err
			}
			// 跟新配置文件
			return ioutil.WriteFile(global.Cfgfile, b, 0644)
		}
	}
	// 添加
	// 默认配置

	ss.Scripts = append(ss.Scripts, s)
	index := len(ss.Scripts) - 1
	fill(index)

	b, err := yaml.Marshal(ss.Scripts)
	if err != nil {
		return err
	}
	// 跟新配置文件
	return ioutil.WriteFile(global.Cfgfile, b, 0644)
}

func DelScript(pname string) error {
	// del := make(chan bool)
	ss.Mu.Lock()
	defer ss.Mu.Unlock()
	if _, ok := ss.Infos[pname]; ok {
		// go func() {
		// wg := &sync.WaitGroup{}
		for name := range ss.Infos[pname] {
			ss.Infos[pname][name].Remove()
		}

	} else {
		return errors.New("not found this pname:" + pname)
	}
	for i, s := range ss.Scripts {
		if s.Name == pname {
			ss.Scripts = append(ss.Scripts[:i], ss.Scripts[i+1:]...)
			delete(ss.Infos, pname)
			break
		}
	}
	b, err := yaml.Marshal(ss.Scripts)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(global.Cfgfile, b, 0644)
}
