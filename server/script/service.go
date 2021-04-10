package script

import (
	"context"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/hyahm/scs/server/alert"
	"github.com/hyahm/scs/server/at"
	"github.com/hyahm/scs/server/lookpath"
	"github.com/hyahm/scs/server/script/status"

	"github.com/hyahm/golog"
)

type Cron struct {
	// 开始执行的时间戳
	Start time.Time
	// 间隔的时间， 如果IsMonth 为true， loop 单位为月， 否则为秒
	IsMonth bool
	Loop    int
}

type Server struct {
	LookPath           []*lookpath.LoopPath
	Name               string
	Dir                string
	Command            string
	Replicate          int
	Always             bool
	Cron               *Cron
	loopTime           time.Time
	IsLoop             bool // 如果是定时任务
	DisableAlert       bool
	DeleteWhenExit     bool
	Env                map[string]string
	SubName            string
	Disable            bool
	Log                map[string][]string
	cmd                *exec.Cmd
	Status             *status.ServiceStatus
	Alert              map[string]alert.SendAlerter
	AT                 *at.AlertTo
	Port               int
	ContinuityInterval time.Duration
	AI                 *alert.AlertInfo // 报警规则
	Exit               chan int         // 9 是退出信号， 10是重启信号
	CancelSigle        chan bool        // 取消停止，重启的信号
	Ctx                context.Context
	Cancel             context.CancelFunc
	Email              []string
	Msg                chan string
	Update             string
	LogLocker          *sync.RWMutex
	Version            string
	SC                 *Script
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

func (s *Server) shell(command string, typ string) error {
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
		cmd.Env = make([]string, 0, len(s.Env))
	}
	cmd.Dir = s.Dir
	for k, v := range s.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
		command = strings.ReplaceAll(command, "$"+k, v)
		command = strings.ReplaceAll(command, "${"+k+"}", v)
	}
	t := time.Now().Format("2006/1/2 15:04:05")
	command = t + " -- " + command
	s.LogLocker.Lock()
	s.Log[typ] = append(s.Log[typ], command)
	s.LogLocker.Unlock()
	read(cmd, s, typ)
	err := cmd.Start()
	if err != nil {
		golog.Error(err)
		return err
	}
	return cmd.Wait()
}

func (s *Server) cron() {
	for {
		select {
		case <-s.Ctx.Done():
			return
		case <-time.After(-time.Since(s.Cron.Start)):
			s.loopTime = time.Now()
			if err := s.start(); err != nil {
				golog.Error(err)
				return
			}

			go s.wait()
			if s.cmd.Process != nil {
				s.Status.Pid = s.cmd.Process.Pid
			}
			return
		}
	}
}

// Start  启动服务
func (s *Server) Start() error {
	golog.Info("start")
	switch s.Status.Status {
	case status.WAITSTOP:
		// 如果之前是等待停止的状态， 更改为重启状态
		// 消耗掉等待停止的队列
		<-s.Exit
		s.Exit <- 10
		s.Status.Status = status.WAITRESTART
	case status.STOP:
		go func() {
			s.Status.Status = status.INSTALL
			if err := s.LookCommandPath(); err != nil {
				golog.Error(err)
				s.Status.Status = status.STOP
				return
			}
			// 创建goroutine 上下文管理
			s.Ctx, s.Cancel = context.WithCancel(context.Background())
			if s.Cron != nil && s.Cron.Loop > 0 {
				// 如果时间没填， 或者已经过去的时间了， 那么就直接启动
				if (s.Cron.Start != time.Time{}) && time.Since(s.Cron.Start) < 0 {
					go s.cron()
					return
				}
				s.loopTime = time.Now()
			}

			if err := s.start(); err != nil {
				s.stopStatus()
				return
			}

			go s.wait()
			if s.cmd.Process != nil {
				s.Status.Pid = s.cmd.Process.Pid
			}
		}()

	}
	return nil
}

// Restart  重动服务
func (s *Server) Restart() {
	if s.IsLoop {
		// 定时器不支持重启操作
		return
	}
	switch s.Status.Status {
	case status.WAITSTOP:
		// 如果之前是等待停止的状态， 更改为重启状态
		<-s.Exit
		s.Exit <- 10
		s.Status.Status = status.WAITRESTART
		return
	case status.RUNNING:
		s.Exit <- 10
		s.Status.Status = status.WAITRESTART
		s.stop()
		return
	case status.STOP:
		s.Start()
	}

}

func (s *Server) Remove() {
	if s == nil {
		return
	}
	switch s.Status.Status {
	case status.WAITRESTART, status.WAITSTOP:
		// 结束发送的退出错误发出的信号
		<-s.Exit
		// 结束停止的goroutine， 转为删除处理
		s.CancelSigle <- true
		go s.remove()
	case status.STOP:
		// 直接删除
		delete(ss.Infos[s.Name], s.SubName)
		if len(ss.Infos[s.Name]) == 0 {
			delete(ss.Infos, s.Name)
		}
	case status.RUNNING:
		go ss.Infos[s.Name][s.SubName].remove()
	default:
		golog.Error("error status")
	}
}

func (s *Server) remove() {
	s.Stop()
	delete(ss.Infos[s.Name], s.SubName)
	if len(ss.Infos[s.Name]) == 0 {
		delete(ss.Infos, s.Name)
	}
}

// Stop  停止服务
func (s *Server) Stop() {
	if s.IsLoop {
		// 如果是定时器， 那么取消定时， 然后执行停止操作
		s.IsLoop = false
		// s.Cancel()
		// s.stopStatus()
	}
	// 只有运行和等待重启才会等待重启
	switch s.Status.Status {
	case status.RUNNING:
		s.Exit <- 9                       // 启动的话， 发送停止信号
		s.Status.Status = status.WAITSTOP // 更改状态为等待停止
		s.stop()
	case status.WAITRESTART:
		<-s.Exit // 取消restart 信号， 改成 等待停止
		s.Exit <- 9
		s.Status.Status = status.WAITSTOP
	}
}

func (s *Server) UpdateAndRestart() {
	golog.Info(s.Update)
	updateCommand := "git pull"
	if s.Update != "" {
		updateCommand = s.Update
	}
	if err := s.shell(updateCommand, "update"); err != nil {
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
	case status.RUNNING:
		s.Exit <- 9
		if err := s.kill(); err != nil {
			s.Cancel()
		}
	case status.WAITRESTART, status.WAITSTOP:
		<-s.Exit
		s.Exit <- 9
		s.kill()
	}

}

func (s *Server) wait() error {
	go s.successAlert()
	if err := s.cmd.Wait(); err != nil {
		s.Cancel()
		// 执行脚本后环境的错误
		select {
		case ec := <-s.Exit:
			switch ec {
			case 9:
				// 主动退出, kill， stop
				s.Status.RestartCount = 0
				s.stopStatus()
				return nil
			case 10:
				// 重启 restart
				s.Status.RestartCount = 0
				s.stopStatus()
				return s.Start()
			}
		default:
			// 意外退出
			golog.Info("error stop")
			if !s.DisableAlert {
				am := &alert.Message{
					Title:  "service error stop",
					Pname:  s.Name,
					Name:   s.SubName,
					Reason: err.Error(),
				}
				if !s.AI.Broken {
					// 第一次
					s.AI.Start = time.Now()
					am.BrokenTime = s.AI.Start.String()
					s.AI.AlertTime = time.Now()
					s.AI.Broken = true
					alert.AlertMessage(am, s.AT)
				} else {
					// 间隔时间内才发送报警
					ci := s.ContinuityInterval
					if ci == 0 {
						ci = time.Hour * 1
					}
					if time.Since(s.AI.AlertTime) >= ci {
						s.AI.AlertTime = time.Now()
						alert.AlertMessage(am, s.AT)
					}
				}
			}
			if s.Always {
				s.stopStatus()
				golog.Info(time.Now())
				// 失败了， 每秒启动一次
				s.Status.RestartCount++
				return s.Start()

			}
		}
		// if s.DeleteWhenExit {
		// 	return config.Cfg.DelScript(s.Name)
		// }
		golog.Debugf("serviceName: %s, subScript: %s, error: %v \n", s.Name, s.SubName, err)
		s.stopStatus()
		s.Status.RestartCount = 0
		// s.Status.Last = false
		return err
	}

	if s.Cron != nil && s.Cron.Loop > 0 {
		// s.stopStatus()
		s.IsLoop = true
		if s.Cron.IsMonth {
			start := time.Since(s.loopTime.AddDate(0, s.Cron.Loop, 0))
			if start < 0 {
				for {
					select {
					case <-s.Ctx.Done():
						golog.Info("loop service have been cancel")
						return nil
					case <-time.After(-start):
						s.stopStatus()
						s.Start()
						return nil
					}
				}
			}
			s.stopStatus()
			s.Start()
			return nil
		} else {
			// start := math.Ceil(float64(s.Cron.Loop) - time.Now().Sub(s.loopTime).Seconds())
			start := time.Since(s.loopTime.Add(time.Duration(s.Cron.Loop) * time.Second))
			if start < 0 {
				// 允许循环， 每s.Loop秒启动一次
				for {
					select {
					case <-s.Ctx.Done():
						golog.Info("loop service have been cancel")
						return nil
					case <-time.After(-start):
						s.stopStatus()
						s.Start()
						return nil
					}
				}
			}
			s.stopStatus()
			s.Start()
			return nil
		}

	}
	// if s.DeleteWhenExit {
	// 	return config.Cfg.DelScript(s.Name)
	// }
	s.stopStatus()
	return nil

}

func (s *Server) stopStatus() {
	s.Status.Status = status.STOP
	s.Status.Pid = 0
	s.Status.Start = 0
	s.cmd = nil
	s.IsLoop = false
}
