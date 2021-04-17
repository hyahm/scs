package scs

import (
	"context"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/hyahm/golog"
)

type Server struct {
	Script             *Script
	Command            string
	Cron               *Cron // 这个cron是新生成的
	IsLoop             bool  // 如果是定时任务
	Env                map[string]string
	SubName            string
	Log                []string
	cmd                *exec.Cmd
	Status             *ServiceStatus
	Alert              map[string]SendAlerter
	AT                 *AlertTo
	Port               int
	ContinuityInterval time.Duration
	AI                 *AlertInfo // 报警规则
	Exit               chan int   // 判断是否是主动退出的
	CancelProcess      chan bool  // 取消操作，
	StopSigle          chan bool  // 停止后发出的信号
	Ctx                context.Context
	Cancel             context.CancelFunc
	Email              []string
	Msg                chan string
	Update             string
	LogLocker          *sync.RWMutex
}

func getVersion(command string) string {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-c", command)
	} else {
		cmd = exec.Command("/bin/bash", "-c", command)
	}
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	output := strings.ReplaceAll(string(out), "\n", "")
	output = strings.ReplaceAll(output, "\r", "")
	return output
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

func (s *Server) cron() {
	s.Status.Status = RUNNING
	for {
		select {
		case <-s.Ctx.Done():
			golog.Info("end loop")
			return
		case <-time.After(-time.Since(s.Cron.StartTime)):
			if err := s.start(); err != nil {
				golog.Error(err)
				// 设置下载启动的时间
				s.Cron.ComputerStartTime()
				continue
			}

			err := s.wait()
			if err != nil {
				s.Cron.ComputerStartTime()
				continue
			}
			if s.cmd.Process != nil {
				s.Status.Pid = s.cmd.Process.Pid
			}
			s.Cron.ComputerStartTime()
			continue
		case <-s.Cron.First:
			if err := s.start(); err != nil {
				golog.Error(err)
				// 设置下载启动的时间
				s.Cron.ComputerStartTime()
				continue
			}

			err := s.wait()
			if err != nil {
				s.Cron.ComputerStartTime()
				continue
			}
			if s.cmd.Process != nil {
				s.Status.Pid = s.cmd.Process.Pid
			}
			s.Cron.ComputerStartTime()
			continue
		}
	}
}

// Start  启动服务
func (svc *Server) Start() error {
	if svc.Script.Disable {
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
			svc.Status.Status = STOP
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
			golog.Info("start time: ", svc.Cron.StartTime)
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

// Restart  重动服务
func (s *Server) Restart() {
	if s.IsLoop {
		s.Cancel()
		s.stopStatus()
	}
	switch s.Status.Status {
	case WAITSTOP:
		// 如果之前是等待停止的状态， 更改为重启状态
		<-s.Exit
		s.Exit <- 10
		s.Status.Status = WAITRESTART
		return
	case RUNNING:
		s.Exit <- 10
		s.Status.Status = WAITRESTART
		s.stop()
		return
	case STOP:
		s.Start()
	}

}

func (s *Server) Remove() {
	switch s.Status.Status {
	case WAITRESTART:
		// 结束发送的退出错误发出的信号
		<-s.Exit
		// 结束停止的goroutine， 转为删除处理
		s.CancelProcess <- true
		go s.remove()
	case STOP, INSTALL:
		// 直接删除
		DeleteServiceBySubName(s.SubName)
	case RUNNING, WAITSTOP:
		go s.remove()
	default:
		golog.Error("error status")
	}
}

func (s *Server) remove() {
	s.Stop()
	// 等待停止信号
	<-s.StopSigle
	DeleteServiceBySubName(s.SubName)
}

// Stop  停止服务
func (s *Server) Stop() {
	if s.IsLoop {
		// 如果是定时任务， 直接停止
		s.Cancel()
		s.stopStatus()
	}
	switch s.Status.Status {
	case RUNNING:
		s.Exit <- 9
		s.Status.Status = WAITSTOP
		s.stop()
	case WAITRESTART:
		<-s.Exit
		s.Exit <- 9
		s.Status.Status = WAITSTOP
	}
}

func GetScriptByPname(name string) (*Script, error) {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	if v, ok := ss.Scripts[name]; ok {
		return v, nil
	} else {
		return nil, ErrFoundPnameOrName
	}
}

func UpdateAndRestartAllServer() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for _, s := range ss.Scripts {
		s.UpdateAndRestartScript()
	}
}

func StartAllServer() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for pname := range ss.Scripts {
		for _, svc := range ss.Infos[pname] {
			svc.Start()
		}
		// s.StartServer()
	}
}

func (s *Script) RemoveScript() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for _, server := range ss.Infos[s.Name] {
		server.Remove()
	}
}

func (s *Script) UpdateAndRestartScript() {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	for _, server := range ss.Infos[s.Name] {
		server.UpdateAndRestart()
	}
}

func (s *Script) EnableScript() error {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	// 禁用 script 所在的所有server
	if _, ok := ss.Scripts[s.Name]; !ok {
		return ErrFoundPnameOrName
	}
	s.Disable = false
	for name := range ss.Infos[s.Name] {
		ss.Infos[s.Name][name].Script.Disable = false
		go ss.Infos[s.Name][name].Stop()
	}
	return nil
}

func (s *Script) DisableScript() error {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	// 禁用 script 所在的所有server
	if _, ok := ss.Scripts[s.Name]; !ok {
		return ErrFoundPnameOrName
	}
	s.Disable = true
	for name := range ss.Infos[s.Name] {
		ss.Infos[s.Name][name].Script.Disable = true
		go ss.Infos[s.Name][name].Stop()
	}
	return nil
}

func (s *Script) StopScript() error {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	// 禁用 script 所在的所有server
	if _, ok := ss.Scripts[s.Name]; !ok {
		return ErrFoundPnameOrName
	}
	for name := range ss.Infos[s.Name] {
		go ss.Infos[s.Name][name].Stop()
	}
	return nil
}

func (s *Script) WaitStopScript() error {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	// 禁用 script 所在的所有server
	if _, ok := ss.Scripts[s.Name]; !ok {
		return ErrFoundPnameOrName
	}
	for subname := range ss.Infos[s.Name] {
		golog.Info(subname)
		ss.Infos[s.Name][subname].Stop()
	}
	return nil
}

func (s *Script) WaitKillScript() error {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	// 禁用 script 所在的所有server
	if _, ok := ss.Scripts[s.Name]; !ok {
		return ErrFoundPnameOrName
	}
	for subname := range ss.Infos[s.Name] {
		golog.Info(subname)
		ss.Infos[s.Name][subname].Kill()
	}
	return nil
}

func (s *Script) RestartScript() error {
	ss.Mu.RLock()
	defer ss.Mu.RUnlock()
	// 禁用 script 所在的所有server
	if _, ok := ss.Scripts[s.Name]; !ok {
		return ErrFoundPnameOrName
	}
	for name := range ss.Infos[s.Name] {
		go ss.Infos[s.Name][name].Restart()
	}
	return nil
}

func (s *Server) UpdateAndRestart() {
	golog.Info(s.Update)
	updateCommand := "git pull"
	if s.Update != "" {
		updateCommand = s.Update
	}
	if err := s.shell(updateCommand); err != nil {
		golog.Error(err)
		return
	}
	s.Restart()
}

// Stop  杀掉服务
func (s *Server) Kill() {
	if s.IsLoop {
		s.Cancel()
		s.stopStatus()
	}
	switch s.Status.Status {
	case RUNNING:
		s.Exit <- 9
		if err := s.kill(); err != nil {
			golog.Error(err)
			// s.Cancel()
		}
	case WAITRESTART, WAITSTOP:
		<-s.Exit
		s.Exit <- 9
		s.kill()
	}

}

func (svc *Server) wait() error {
	go svc.successAlert()
	if err := svc.cmd.Wait(); err != nil {
		svc.Cancel()
		// 执行脚本后环境的错误
		select {
		case ec := <-svc.Exit:
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
			golog.Info("error stop")
			if !svc.Script.DisableAlert {
				am := &Message{
					Title:  "service error stop",
					Pname:  svc.Script.Name,
					Name:   svc.SubName,
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
				return errors.New("stoped")
			}
			if svc.Script.Always {
				svc.stopStatus()
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
	if svc.Script.DeleteWhenExit {
		return Cfg.DelScript(svc.Script.Name)
	}
	svc.stopStatus()
	return nil

}

func (s *Server) stopStatus() {
	s.Status.Status = STOP
	s.Status.Pid = 0
	s.Status.RestartCount = 0
	s.Status.Start = 0
	s.IsLoop = false
}
