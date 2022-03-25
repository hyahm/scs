package server

// "github.com/hyahm/scs/internal/config/scripts"

// 保存的所有脚本相关的配置
// var ss *Service

// func init() {
// 	ss = &Service{
// 		// 由2层组成， 一级是name  二级是pname
// 		Servers: make(map[subname.Subname]*Server),
// 		Scripts: make(map[string]*scripts.Script), // 保存脚本
// 		Mu:      &sync.RWMutex{},
// 	}
// }

// type Service struct {
// 	Servers map[subname.Subname]*Server // 根据subname存放server的
// 	Scripts map[string]*scripts.Script  // 根据name存放脚本的
// 	Mu      *sync.RWMutex
// }

// func GetServers() []byte {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	send, err := json.MarshalIndent(ss.Servers, "", "  ")
// 	// send, err := json.Marshal(ss.Infos)
// 	if err != nil {
// 		golog.Error(err)
// 	}
// 	return send
// }

// func TempScript(temp map[string]struct{}) {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	for name := range ss.Scripts {
// 		temp[name] = struct{}{}
// 	}
// }

// func NeedStart(subname string) bool {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	if _, ok := ss.Scripts[subname]; !ok {
// 		return false
// 	}
// 	if ss.Scripts[subname].Disable {
// 		return false
// 	}
// 	return true
// }

// func DelScript(pname string) {
// 	ss.Mu.Lock()
// 	defer ss.Mu.Unlock()
// 	delete(ss.Scripts, pname)
// }

// func RemoveScript(pname string) error {
// 	// ss.Mu.RLock()
// 	// defer ss.Mu.RUnlock()

// 	if _, ok := ss.Scripts[pname]; ok {
// 		replicate := ss.Scripts[pname].Replicate
// 		if replicate == 0 {
// 			replicate = 1
// 		}

// 		for i := 0; i < replicate; i++ {
// 			subname := subname.NewSubname(pname, i)
// 			ss.Servers[subname].Remove()
// 		}

// 	} else {
// 		return errors.New("not found this pname:" + pname)
// 	}
// 	return nil
// }

// func AddAndStartServer(script *scripts.Script) {
// 	if !config.CheckScriptName(script.Name) {
// 		golog.Error("script name must be a word, " + script.Name)
// 		return
// 	}
// 	ss.Mu.Lock()
// 	ss.Scripts[script.Name] = script
// 	ss.Mu.Unlock()
// 	ss.Scripts[script.Name].MakeServer()
// }

// func GetScripts() []byte {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	send, err := json.MarshalIndent(ss.Scripts, "", "  ")
// 	// send, err := json.Marshal(ss.Scripts)
// 	if err != nil {
// 		golog.Error(err)
// 	}
// 	return send
// }

// func HasName(name string) bool {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	if _, ok := ss.Scripts[name]; ok {
// 		return true
// 	}
// 	return false
// }

// func HassubName(name subname.Subname) bool {

// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	if _, ok := ss.Servers[name]; ok {
// 		return true
// 	}
// 	return false
// }

// func GetServerBySubname(subname subname.Subname) (*server.Server, bool) {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	v, ok := ss.Servers[subname]
// 	return v, ok
// }

// func GetServersByName(name string) (map[subname.Subname]*server.Server, bool) {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	if _, ok := ss.Scripts[name]; !ok {
// 		return nil, false
// 	}
// 	servers := make(map[subname.Subname]*server.Server)
// 	replicate := ss.Scripts[name].Replicate
// 	if replicate == 0 {
// 		replicate = 1
// 	}
// 	for i := 0; i < replicate; i++ {
// 		subname := subname.NewSubname(name, i)
// 		servers[subname] = ss.Servers[subname]
// 	}

// 	return servers, true
// }

// // 获取脚本结构体
// func GetServerByNameAndSubname(name string, subname subname.Subname) (*server.Server, bool) {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	if _, ok := ss.Scripts[name]; !ok {
// 		return nil, false
// 	}
// 	if _, ok := ss.Servers[subname]; ok {
// 		return ss.Servers[subname], true
// 	}

// 	return nil, false
// }

// func StopAllServer() {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	for _, s := range ss.Scripts {
// 		err := StopScript(s)
// 		if err != nil {
// 			golog.Error(err)
// 		}
// 	}
// }
// func WaitStopAllServer() {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	for _, s := range ss.Scripts {
// 		WaitStopScript(s)
// 	}
// }

// func WaitKillAllServer() {
// 	// ss.ScriptLocker.RLock()
// 	// defer ss.ScriptLocker.RUnlock()
// 	for _, s := range ss.Scripts {
// 		WaitKillScript(s)
// 	}
// }

// // 只删除ss.infos 里面的额
// func DeleteServiceBySubName(subname subname.Subname) error {
// 	// 删除server
// 	ss.Mu.Lock()
// 	defer ss.Mu.Unlock()
// 	if _, ok := ss.Servers[subname]; ok {
// 		delete(ss.Servers, subname)
// 		// 同时script 里面也要删除
// 		name := subname.GetName()
// 		if _, ok := ss.Scripts[name]; ok {
// 			// ss.Scripts[name].Replicate--

// 			if ss.Scripts[name].Replicate == 1 {
// 				ss.Scripts[name].Replicate = 0
// 			}

// 			golog.Info(ss.Scripts[name].Replicate)
// 			// 开发中， replicate =0 或 1 其实都是1 的意思， 所以减一后 <= 0 的其实就是都删除干净了的意思
// 			if ss.Scripts[name].Replicate <= 0 {
// 				delete(ss.Scripts, name)
// 				delete(ss.Servers, subname)
// 			}
// 		}
// 		return nil
// 	}
// 	return errors.New("")

// }

// func GetAllServer() {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	for _, svc := range ss.Servers {
// 		svc.Start()
// 	}
// }

// func UpdateAndRestartAllServer() {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	for _, s := range ss.Scripts {
// 		go UpdateAndRestartScript(s)
// 	}
// }

// func RestartAllServer() {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	for _, svc := range ss.Servers {
// 		go svc.Restart()
// 	}
// }

// func Len() int {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	return len(ss.Scripts)
// }

// // 从 ss 中删除某一个subname
// func DeleteSubname(subname subname.Subname) {

// 	ss.Mu.Lock()
// 	defer ss.Mu.Unlock()
// 	// 以最后一个下划线来分割出pname
// 	delete(ss.Servers, subname)
// }

// func (sl *StatusList) Filter(filter []string) {
// 	temp := make([]*status.ServiceStatus, 0, len(sl.Data))

// 	for _, s := range sl.Data {
// 		for _, f := range filter {
// 			if strings.Contains(s.Name, f) {
// 				temp = append(temp, s)
// 				break
// 			}
// 		}
// 	}
// 	sl.Data = temp
// }

// // 获取所有服务的状态
// func All() []byte {
// 	// ss.Mu.RLock()
// 	// defer ss.Mu.RUnlock()
// 	statuss := &StatusList{
// 		Data: make([]*status.ServiceStatus, 0),
// 	}
// 	// ss := make([]*ServiceStatus, 0)
// 	for subname := range ss.Servers {
// 		pname := subname.GetName()
// 		if _, ok := ss.Scripts[pname]; !ok {
// 			golog.Debug("not found name: " + pname)
// 			continue
// 		}

// 		status := &status.ServiceStatus{
// 			PName:        pname,
// 			Name:         subname.String(),
// 			Command:      ss.Servers[subname].Command,
// 			Always:       ss.Scripts[pname].Always,
// 			Version:      ss.Servers[subname].Version,
// 			Status:       ss.Servers[subname].Status.Status,
// 			CanNotStop:   ss.Servers[subname].Status.CanNotStop,
// 			Path:         ss.Servers[subname].Script.Dir,
// 			Start:        ss.Servers[subname].Status.Start,
// 			RestartCount: ss.Servers[subname].Status.RestartCount,
// 			Up:           ss.Servers[subname].Status.Up,
// 			Disable:      ss.Scripts[pname].Disable,
// 		}
// 		if ss.Servers[subname].Cmd != nil && ss.Servers[subname].Cmd.Process != nil {
// 			status.Pid = ss.Servers[subname].Status.Pid
// 			status.Cpu, status.Mem, _ = probe.GetProcessInfo(int32(ss.Servers[subname].Status.Pid))

// 		}
// 		statuss.Data = append(statuss.Data, status)
// 	}
// 	statuss.Code = 200
// 	send, err := json.MarshalIndent(statuss, "", "\t")
// 	if err != nil {
// 		golog.Error(err)
// 	}
// 	return send
// }

// func ScriptPname(pname string) []byte {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	statuss := &StatusList{
// 		Data: make([]*status.ServiceStatus, 0),
// 	}
// 	if _, ok := ss.Scripts[pname]; !ok {
// 		statuss.Msg = "not found " + pname
// 		send, err := json.MarshalIndent(statuss, "", "\n")

// 		if err != nil {
// 			golog.Error(err)
// 		}
// 		return send
// 	}
// 	replicate := ss.Scripts[pname].Replicate
// 	if replicate == 0 {
// 		replicate = 1
// 	}

// 	for i := 0; i < replicate; i++ {
// 		subname := subname.NewSubname(pname, i)
// 		status := &status.ServiceStatus{
// 			PName:        subname.GetName(),
// 			Name:         subname.String(),
// 			Command:      ss.Servers[subname].Status.Command,
// 			Always:       ss.Scripts[subname.GetName()].Always,
// 			Version:      ss.Servers[subname].Status.Version,
// 			CanNotStop:   ss.Servers[subname].Status.CanNotStop,
// 			Path:         ss.Servers[subname].Status.Path,
// 			Status:       ss.Servers[subname].Status.Status,
// 			RestartCount: ss.Servers[subname].Status.RestartCount,
// 			Up:           ss.Servers[subname].Status.Up,
// 			Disable:      ss.Scripts[subname.GetName()].Disable,
// 		}
// 		if ss.Servers[subname].Cmd != nil && ss.Servers[subname].Cmd.Process != nil {
// 			status.Pid = ss.Servers[subname].Cmd.Process.Pid
// 			status.Cpu, status.Mem, _ = probe.GetProcessInfo(int32(ss.Servers[subname].Cmd.Process.Pid))

// 		}
// 		statuss.Data = append(statuss.Data, status)
// 	}

// 	statuss.Code = 200
// 	send, err := json.MarshalIndent(statuss, "", "\n")

// 	if err != nil {
// 		golog.Error(err)
// 	}
// 	return send
// }

// func ScriptName(pname string, subname subname.Subname) []byte {
// 	// ss.Mu.RLock()
// 	// defer ss.Mu.RUnlock()
// 	statuss := &StatusList{
// 		Data: make([]*status.ServiceStatus, 0),
// 	}
// 	if _, ok := ss.Scripts[pname]; !ok {
// 		statuss.Msg = "not found " + pname
// 		send, err := json.MarshalIndent(statuss, "", "\n")

// 		if err != nil {
// 			golog.Error(err)
// 		}
// 		return send
// 	}

// 	status := &status.ServiceStatus{
// 		PName:        subname.GetName(),
// 		Name:         subname.String(),
// 		Command:      ss.Servers[subname].Status.Command,
// 		Always:       ss.Scripts[subname.GetName()].Always,
// 		Version:      ss.Servers[subname].Status.Version,
// 		CanNotStop:   ss.Servers[subname].Status.CanNotStop,
// 		Path:         ss.Servers[subname].Status.Path,
// 		Status:       ss.Servers[subname].Status.Status,
// 		RestartCount: ss.Servers[subname].Status.RestartCount,
// 		Up:           ss.Servers[subname].Status.Up,
// 		Disable:      ss.Scripts[subname.GetName()].Disable,
// 	}
// 	if ss.Servers[subname].Cmd != nil && ss.Servers[subname].Cmd.Process != nil {
// 		status.Pid = ss.Servers[subname].Cmd.Process.Pid
// 		status.Cpu, status.Mem, _ = probe.GetProcessInfo(int32(ss.Servers[subname].Cmd.Process.Pid))

// 	}
// 	statuss.Data = append(statuss.Data, status)
// 	send, err := json.MarshalIndent(statuss, "", "\t")
// 	if err != nil {
// 		golog.Error(err)
// 	}
// 	statuss.Code = 200
// 	return send
// }

// func NeedStop(s *scripts.Script) bool {
// 	// 更新server
// 	// 判断值是否相等
// 	if s.Dir != ss.Scripts[s.Name].Dir ||
// 		s.Command != ss.Scripts[s.Name].Command ||
// 		s.Replicate != ss.Scripts[s.Name].Replicate ||
// 		s.Always != ss.Scripts[s.Name].Always ||
// 		s.DisableAlert != ss.Scripts[s.Name].DisableAlert ||
// 		!pkg.CompareMap(s.Env, ss.Scripts[s.Name].Env) ||
// 		s.Port != ss.Scripts[s.Name].Port ||
// 		s.Version != ss.Scripts[s.Name].Version ||
// 		s.Disable != ss.Scripts[s.Name].Disable ||
// 		s.Update != ss.Scripts[s.Name].Update ||
// 		s.DeleteWhenExit != ss.Scripts[s.Name].DeleteWhenExit ||
// 		!s.Cron.IsEqual(ss.Scripts[s.Name].Cron) ||
// 		!IsEqual(s.Name, s.AT) {
// 		// 如果值有变动， 那么需要重新server
// 		// 先同步停止之前的server， 然后启动新的server
// 		// server 是单独的， 在通知后需要同步更新server
// 		return true
// 	}
// 	return false
// }

// // 启动方法， 异步执行
// func StartServer(s *scripts.Script) {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	replicate := s.Replicate
// 	if replicate == 0 {
// 		replicate = 1
// 	}
// 	for i := 0; i < replicate; i++ {
// 		subname := subname.NewSubname(s.Name, i)
// 		ss.Servers[subname].Start()
// 	}
// }

// func KillScript(s *scripts.Script) {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	replicate := s.Replicate
// 	if replicate == 0 {
// 		replicate = 1
// 	}
// 	for i := 0; i < replicate; i++ {
// 		subname := subname.NewSubname(s.Name, i)
// 		ss.Servers[subname].Kill()
// 	}
// }

// // 通过script 生成和启动服务
// func AddScript(s *scripts.Script) {
// 	// 通过script 生成server
// 	s.MakeServer()
// 	StartServer(s)
// }

// // 异步重启
// func RestartScript(s *scripts.Script) error {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	// 禁用 script 所在的所有server
// 	if _, ok := ss.Scripts[s.Name]; !ok {
// 		return errors.New("")
// 	}
// 	replicate := s.Replicate
// 	if replicate == 0 {
// 		replicate = 1
// 	}
// 	for i := 0; i < replicate; i++ {
// 		subname := subname.NewSubname(s.Name, i)
// 		golog.Info("restart: ", subname)
// 		go ss.Servers[subname].Restart()
// 	}
// 	return nil
// }

// // 同步杀掉
// func WaitKillScript(s *scripts.Script) {
// 	// ss.ServerLocker.RLock()
// 	// defer ss.ServerLocker.RUnlock()
// 	// 禁用 script 所在的所有server
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	// 禁用 script 所在的所有server
// 	replicate := s.Replicate
// 	if replicate == 0 {
// 		replicate = 1
// 	}
// 	for i := 0; i < replicate; i++ {
// 		subname := subname.NewSubname(s.Name, i)
// 		ss.Servers[subname].Kill()
// 	}
// }

// // 同步停止
// func WaitStopScript(s *scripts.Script) {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	// 禁用 script 所在的所有server
// 	replicate := s.Replicate
// 	if replicate == 0 {
// 		replicate = 1
// 	}
// 	for i := 0; i < replicate; i++ {
// 		subname := subname.NewSubname(s.Name, i)
// 		ss.Servers[subname].Stop()
// 	}
// }

// // 异步执行停止脚本
// func StopScript(s *scripts.Script) error {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	if _, ok := ss.Scripts[s.Name]; !ok {
// 		return errors.New("")
// 	}
// 	// 禁用 script 所在的所有server
// 	replicate := s.Replicate
// 	if replicate == 0 {
// 		replicate = 1
// 	}
// 	for i := 0; i < replicate; i++ {
// 		subname := subname.NewSubname(s.Name, i)
// 		go ss.Servers[subname].Stop()
// 	}
// 	return nil
// }

// func GetScriptByPname(name string) (*scripts.Script, bool) {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	v, ok := ss.Scripts[name]
// 	return v, ok

// }

// // 返回成功还是失败
// func UpdateAndRestartScript(s *scripts.Script) bool {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	if _, ok := ss.Scripts[s.Name]; !ok {
// 		return false
// 	}
// 	replicate := s.Replicate
// 	if replicate == 0 {
// 		replicate = 1
// 	}
// 	for i := 0; i < replicate; i++ {
// 		subname := subname.NewSubname(s.Name, i)
// 		go ss.Servers[subname].UpdateAndRestart()
// 	}
// 	return true
// }

// func EnableScript(s *scripts.Script) bool {
// 	ss.Mu.Lock()
// 	defer ss.Mu.Unlock()
// 	// 禁用 script 所在的所有server
// 	if _, ok := ss.Scripts[s.Name]; !ok {
// 		return false
// 	}
// 	ss.Scripts[s.Name].Disable = false
// 	replicate := s.Replicate
// 	if replicate == 0 {
// 		replicate = 1
// 	}
// 	for i := 0; i < replicate; i++ {
// 		subname := subname.NewSubname(s.Name, i)
// 		go ss.Servers[subname].Start()
// 	}
// 	return true
// }

// func DisableScript(s *scripts.Script) bool {
// 	ss.Mu.Lock()
// 	defer ss.Mu.Unlock()
// 	// 禁用 script 所在的所有server
// 	if _, ok := ss.Scripts[s.Name]; !ok {
// 		return false
// 	}
// 	ss.Scripts[s.Name].Disable = true
// 	replicate := s.Replicate
// 	if replicate == 0 {
// 		replicate = 1
// 	}
// 	for i := 0; i < replicate; i++ {
// 		subname := subname.NewSubname(s.Name, i)
// 		go ss.Servers[subname].Stop()
// 	}
// 	return true
// }

// func AddInfo(name subname.Subname, svc *server.Server) {

// }

// // 比较新的与之前的是否相等， 调用者必须是新的
// func IsEqual(pname string, at *to.AlertTo) bool {
// 	if at == nil && ss.Scripts[pname].AT == nil {
// 		return true
// 	}
// 	if (at == nil && ss.Scripts[pname].AT != nil) || (at != nil && ss.Scripts[pname].AT == nil) {
// 		return false
// 	}
// 	if !pkg.CompareSlice(at.Email, ss.Scripts[pname].AT.Email) ||
// 		!pkg.CompareSlice(at.Rocket, ss.Scripts[pname].AT.Rocket) ||
// 		!pkg.CompareSlice(at.Telegram, ss.Scripts[pname].AT.Telegram) ||
// 		!pkg.CompareSlice(at.WeiXin, ss.Scripts[pname].AT.WeiXin) {
// 		return false
// 	}
// 	return true
// }

// func StartAllServer() {
// 	ss.Mu.RLock()
// 	defer ss.Mu.RUnlock()
// 	for _, v := range ss.Servers {
// 		v.Start()
// 	}
// }

// // 扩缩容script, 接受2个参数， name string  replicate int
// func MakeReplicateServer(name string, replicate int) {
// 	ss.Mu.Lock()
// 	defer ss.Mu.Unlock()
// 	ss.Scripts[s.Name].Replicate = s.Replicate
// }
