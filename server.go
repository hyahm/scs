package scs

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/hyahm/golog"
)

type Server struct {
	Script             *Script                `json:"-"`
	Command            string                 `json:"command"`
	Version            string                 `json:"version"`
	Cron               *Cron                  `json:"-"` // 这个cron是新生成的
	IsLoop             bool                   `json:"-"` // 如果是定时任务
	Env                map[string]string      `json:"-"`
	SubName            Subname                `json:"subname"`
	Log                []string               `json:"-"`
	cmd                *exec.Cmd              `json:"-"`
	Status             *ServiceStatus         `json:"status"`
	Alert              map[string]SendAlerter `json:"-"`
	AT                 *AlertTo               `json:"at"`
	Port               int                    `json:"port"`
	ContinuityInterval time.Duration          `json:"continuityInterval"`
	AI                 *AlertInfo             `json:"ai"` // 报警规则
	Exit               chan int               `json:"-"`  // 判断是否是主动退出的
	CancelProcess      chan bool              `json:"-"`  // 取消操作，
	StopSigle          chan bool              `json:"-"`  // 停止后发出的信号
	Ctx                context.Context        `json:"-"`
	Cancel             context.CancelFunc     `json:"-"`
	Msg                chan string            `json:"-"`
	removed            bool                   `json"-"` // 标识是否已经被删除
	Update             string                 `json:"update"`
	LogLocker          *sync.RWMutex          `json:"-"`
}

func getVersion(command string) string {
	var cmd *exec.Cmd
	golog.Info(command)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-c", command)
	} else {
		cmd = exec.Command("/bin/bash", "-c", command)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	s := new(string)
	ch := make(chan struct{})
	defer cancel()
	go func(s *string) {
		out, err := cmd.Output()
		if err != nil {
			return
		}
		*s = string(out)
		ch <- struct{}{}
	}(s)

	select {
	case <-ctx.Done():
		return ""
	case <-ch:
		output := strings.ReplaceAll(*s, "\n", "")
		output = strings.ReplaceAll(output, "\r", "")
		return output
	}

}

func (svc *Server) shell(command string) error {
	var cmd *exec.Cmd
	golog.Info(command)
	// s.co = strings.ReplaceAll(s.comm, "$NAME", subname)
	// command = strings.ReplaceAll(command, "$PNAME", c.SC[index].Name)
	// command = strings.ReplaceAll(command, "$PORT", strconv.Itoa(c.SC[index].Port+i))
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-c", command)
	} else {
		cmd = exec.Command("/bin/bash", "-c", command)
	}
	if cmd.Env == nil {
		cmd.Env = make([]string, 0, len(svc.Env))
	}
	cmd.Dir = svc.Script.Dir
	for k, v := range svc.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
		command = strings.ReplaceAll(command, "$"+k, v)
		command = strings.ReplaceAll(command, "${"+k+"}", v)
	}
	read(cmd, svc)
	err := cmd.Start()
	if err != nil {
		golog.Error(err)
		return err
	}
	return cmd.Wait()
}

func (svc *Server) cron() {
	svc.Status.Status = RUNNING

	for {
		select {
		case <-svc.Ctx.Done():
			svc.StopSigle <- true
			golog.Info("name:" + svc.SubName + " end cron")
			return
		case <-time.After(-time.Since(svc.Cron.StartTime)):
			golog.Info("start time: ", svc.Cron.StartTime)
			if err := svc.start(); err != nil {
				golog.Error(err)
				// 设置下载启动的时间
				svc.Cron.ComputerStartTime()
				continue
			}
			if svc.cmd != nil && svc.cmd.Process != nil {
				svc.Status.Pid = svc.cmd.Process.Pid
			}
			err := svc.wait()
			if err != nil {
				golog.Error(err)
				svc.Cron.ComputerStartTime()
				continue
			}
			svc.Cron.ComputerStartTime()
			golog.Infof("cron task: %s have been completed", svc.SubName)
			continue
		case <-svc.Cron.First:
			golog.Info("start time: ", time.Now())
			if err := svc.start(); err != nil {
				golog.Error(err)
				// 设置下载启动的时间
				svc.Cron.ComputerStartTime()
				continue
			}
			if svc.cmd != nil && svc.cmd.Process != nil {
				svc.Status.Pid = svc.cmd.Process.Pid
			}
			err := svc.wait()
			if err != nil {
				golog.Error(err)
				svc.Cron.ComputerStartTime()
				continue
			}

			svc.Cron.ComputerStartTime()
			golog.Infof("cron task: %s have been completed", svc.SubName)
			continue
		}
	}
}

// Start  启动服务 异步的
func (svc *Server) Start() error {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Scripts[svc.SubName.GetName()]; ok && ss.Scripts[svc.SubName.GetName()].Disable {
		return nil
	}

	switch svc.Status.Status {
	case WAITSTOP:
		// 如果之前是等待停止的状态， 更改为重启状态
		<-svc.Exit
		svc.Exit <- 10
		svc.Status.Status = WAITRESTART
	case STOP:
		svc.Status.Status = INSTALL
		err := svc.LookCommandPath()
		if err != nil {
			golog.Error(err)
			svc.stopStatus()
			return err
		}
		svc.Exit = make(chan int, 2)
		svc.CancelProcess = make(chan bool, 2)
		svc.Ctx, svc.Cancel = context.WithCancel(context.Background())
		if svc.Cron != nil && svc.Cron.Loop > 0 {
			svc.IsLoop = true
			// 循环的起止时间可以只设置时分秒， 自动补齐今天的日期
			svc.Cron.Start = strings.Trim(svc.Cron.Start, " ")
			if svc.Cron.Start != "" {
				// 计算下次启动的时间
				index := strings.Index(svc.Cron.Start, " ")
				if index < 0 {
					// 如果只有时间， 自动获取今天的年月日
					svc.Cron.Start = strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")[0] + " " + svc.Cron.Start
				}
				svc.Cron.StartTime, err = time.ParseInLocation("2006-01-02 15:04:05", svc.Cron.Start, time.Local)
				if err != nil {
					golog.Error(err)
				}
				// 比较是否过了时间点， 如果过了就重新计算， 否则就是给定的时间
				if time.Since(svc.Cron.StartTime) > 0 {
					svc.Cron.loopTime = time.Duration(svc.Cron.Loop) * time.Second
					svc.Cron.StartTime = svc.Cron.StartTime.Add(svc.Cron.loopTime)
					if svc.Cron.IsMonth {
						svc.Cron.StartTime = svc.Cron.StartTime.AddDate(0, svc.Cron.Loop, 0)
					}
				}
			} else {
				// 如果没设置， 设置下此启动的时间为当前时间
				svc.Cron.First = make(chan bool, 1)
				svc.Cron.First <- true
			}
			// 如果有定时任务， 那么时间到了就执行
			// 保留时间

			go svc.cron()
			return nil
		}

		if err := svc.start(); err != nil {
			svc.stopStatus()
			return err
		}

		go svc.wait()
		if svc.cmd.Process != nil {
			svc.Status.Pid = svc.cmd.Process.Pid
		}

	}
	return nil
}

// Restart  重动服务, 同步执行的
func (svc *Server) Restart() {
	if svc.IsLoop {
		svc.Cancel()
		// 如果是循环的就直接退出
		return
	}
	switch svc.Status.Status {
	case WAITSTOP:
		// 如果之前是等待停止的状态， 更改为重启状态
		<-svc.Exit
		svc.Exit <- 10
		svc.Status.Status = WAITRESTART
		return
	case RUNNING:
		svc.Exit <- 10
		svc.Status.Status = WAITRESTART
		svc.stop()
		return
	case STOP:
		svc.Start()
	}

}

// 异步删除
func (svc *Server) Remove() {
	if svc.removed {
		return
	}
	svc.removed = true
	switch svc.Status.Status {
	case WAITRESTART:
		// 结束发送的退出错误发出的信号
		<-svc.Exit
		// 结束停止的goroutine， 转为删除处理
		svc.CancelProcess <- true
		go svc.remove()
	case STOP, INSTALL:
		// 直接删除
		DeleteServiceBySubName(svc.SubName)
	case RUNNING, WAITSTOP:
		go svc.remove()
	default:
		golog.Error("error status")
	}
}

func (svc *Server) remove() {
	svc.Stop()

	// 等待停止信号
	<-svc.StopSigle
	golog.Infof("%s stoped", svc.SubName)
	DeleteServiceBySubName(svc.SubName)
}

// Stop  停止服务
func (svc *Server) Stop() {
	if svc.IsLoop {
		// 如果是定时任务， 直接停止
		golog.Infof("stop loop %s", svc.SubName)
		svc.Cancel()
		svc.stopStatus()
	}
	switch svc.Status.Status {
	case RUNNING:
		svc.Exit <- 9
		svc.Status.Status = WAITSTOP
		svc.stop()
	case WAITRESTART:
		<-svc.Exit
		svc.Exit <- 9
		svc.Status.Status = WAITSTOP
	}
}

func GetScriptByPname(name string) (*Script, error) {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if v, ok := ss.Scripts[name]; ok {
		return v, nil
	}
	return nil, ErrFoundPname

}

func UpdateAndRestartAllServer() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for _, s := range ss.Scripts {
		go s.UpdateAndRestartScript()
	}
}

func StartAllServer() {
	for _, svc := range ss.Infos {
		svc.Start()
	}
}

func (s *Script) RemoveScript() error {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Scripts[s.Name]; !ok {
		return ErrFoundPname
	}
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := NewSubname(s.Name, i)
		golog.Warn(subname)
		ss.Infos[subname].Remove()
	}
	return nil
}

func (s *Script) UpdateAndRestartScript() error {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Scripts[s.Name]; !ok {
		return ErrFoundPname
	}
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := NewSubname(s.Name, i)
		go ss.Infos[subname].UpdateAndRestart()
	}
	return nil
}

func (s *Script) EnableScript() error {
	ss.Mu.Lock()
	defer ss.Mu.Unlock()
	// 禁用 script 所在的所有server
	if _, ok := ss.Scripts[s.Name]; !ok {
		return ErrFoundPname
	}
	ss.Scripts[s.Name].Disable = true
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := NewSubname(s.Name, i)
		go ss.Infos[subname].Start()
	}
	return nil
}

func (s *Script) DisableScript() error {
	ss.Mu.Lock()
	defer ss.Mu.Unlock()
	// 禁用 script 所在的所有server
	if _, ok := ss.Scripts[s.Name]; !ok {
		return ErrFoundPname
	}
	ss.Scripts[s.Name].Disable = true
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := NewSubname(s.Name, i)
		go ss.Infos[subname].Stop()
	}
	return nil
}

// 异步执行停止脚本
func (s *Script) StopScript() error {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if _, ok := ss.Scripts[s.Name]; !ok {
		return ErrFoundPname
	}
	// 禁用 script 所在的所有server
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := NewSubname(s.Name, i)
		go ss.Infos[subname].Stop()
	}
	return nil
}

// 同步停止
func (svc *Script) WaitStopScript() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	// 禁用 script 所在的所有server
	replicate := svc.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := NewSubname(svc.Name, i)
		ss.Infos[subname].Stop()
	}
}

// 同步杀掉
func (svc *Script) WaitKillScript() {
	// ss.ServerLocker.RLock()
	// defer ss.ServerLocker.RUnlock()
	// 禁用 script 所在的所有server
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	// 禁用 script 所在的所有server
	replicate := svc.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := NewSubname(svc.Name, i)
		ss.Infos[subname].Kill()
	}
}

// 异步重启
func (svc *Script) RestartScript() error {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	// 禁用 script 所在的所有server
	if _, ok := ss.Scripts[svc.Name]; !ok {
		return ErrFoundPname
	}
	replicate := svc.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := NewSubname(svc.Name, i)
		go ss.Infos[subname].Restart()
	}
	return nil
}

// 同步更新并重启
func (svc *Server) UpdateAndRestart() {
	updateCommand := "git pull"
	if svc.Update != "" {
		updateCommand = svc.Update
	}
	if err := svc.shell(updateCommand); err != nil {
		golog.Error(err)
		return
	}
	svc.Restart()
}

// Stop  杀掉服务, 没有产生goroutine， 直接杀死
func (svc *Server) Kill() {
	if svc.IsLoop {
		svc.Cancel()
		svc.stopStatus()
	}
	switch svc.Status.Status {
	case RUNNING:
		svc.Exit <- 9
		if err := svc.kill(); err != nil {
			golog.Error(err)
			// s.Cancel()
		}
	case WAITRESTART, WAITSTOP:
		<-svc.Exit
		svc.Exit <- 9
		svc.kill()
	}

}

func (svc *Server) wait() error {
	go svc.successAlert()
	if err := svc.cmd.Wait(); err != nil {

		// 脚本退出后才会执行这里的代码
		select {
		case ec := <-svc.Exit:
			svc.Cancel()
			switch ec {
			case 9:
				// 主动退出, kill， stop
				svc.Status.RestartCount = 0
				svc.stopStatus()
				return nil
			case 10:
				// 重启 restart
				svc.Status.RestartCount = 0
				svc.stopStatus()
				return svc.Start()
			}
		default:
			// 意外退出
			golog.Error("error stop")
			if !svc.Script.DisableAlert && HaveAlert() {
				am := &Message{
					Title:  "service error stop",
					Pname:  svc.Script.Name,
					Name:   svc.SubName.String(),
					Reason: err.Error(),
				}
				if !svc.AI.Broken {
					// 第一次
					svc.AI.Start = time.Now()
					am.BrokenTime = svc.AI.Start.String()
					svc.AI.AlertTime = time.Now()
					svc.AI.Broken = true
					AlertMessage(am, svc.AT)
				} else {
					// 间隔时间内才发送报警
					ci := svc.ContinuityInterval
					if ci == 0 {
						ci = time.Hour * 1
					}
					if time.Since(svc.AI.AlertTime) >= ci {
						svc.AI.AlertTime = time.Now()
						AlertMessage(am, svc.AT)
					}
				}
			}
			// 如果是定时器的话， 直接结束
			if svc.Cron != nil && svc.Cron.Loop > 0 {
				svc.cmd = nil
				return nil
			}
			svc.Cancel()
			if svc.Script.Always {
				golog.Info("restart +1")
				svc.Status.Status = STOP
				svc.Status.Pid = 0
				svc.Status.Start = 0
				// 失败了， 每秒启动一次
				svc.Status.RestartCount++
				return svc.Start()

			}
		}
		if svc.Script.DeleteWhenExit {
			return Cfg.DelScript(svc.Script.Name)
		}
		golog.Debugf("serviceName: %s, subScript: %s, error: %v \n", svc.Script.Name, svc.SubName, err)
		svc.stopStatus()
		return err
	}
	if svc.IsLoop {
		// 如果是个定时器， 那么不修改为停止
		svc.cmd = nil
		return nil
	}
	svc.Cancel()
	if svc.Script.DeleteWhenExit {
		return Cfg.DelScript(svc.Script.Name)
	}
	svc.stopStatus()
	return nil

}

func (svc *Server) stopStatus() {
	svc.Status.Status = STOP
	svc.Status.Pid = 0
	svc.Status.CanNotStop = false
	svc.Status.RestartCount = 0
	svc.Status.Start = 0
	svc.IsLoop = false
	svc.removed = false
}

func (svc *Server) LookCommandPath() error {
	for _, v := range svc.Script.LookPath {
		if strings.Trim(v.Path, " ") == "" && strings.Trim(v.Command, " ") == "" {
			continue
		}
		if strings.Trim(v.Path, " ") != "" {
			golog.Info("check path: ", v.Path)
			_, err := os.Stat(v.Path)
			if !os.IsNotExist(err) {
				continue
			}
		}
		if strings.Trim(v.Command, " ") != "" {
			golog.Info("check command: ", v.Command)
			_, err := exec.LookPath(v.Command)
			if err == nil {
				continue
			}
		}
		golog.Info("exec: ", v.Install)
		if err := svc.shell(v.Install); err != nil {
			golog.Error(v.Install)
			return err
		}
	}
	return nil
}
