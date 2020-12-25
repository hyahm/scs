package script

import (
	"context"
	"os/exec"
	"scs/alert"
	"scs/internal"
	"time"

	"github.com/hyahm/golog"
)

type Cron struct {
	// 开始执行的时间戳
	Start time.Time
	// 间隔的时间， 如果IsMonth 为true， loop 单位为月， 否则为秒
	IsMonth bool
	Loop    int
}

type Script struct {
	LookPath           []*internal.LoopPath
	Name               string
	Dir                string
	Command            string
	Replicate          int
	Always             bool
	Cron               *Cron
	loopTime           time.Time
	DisableAlert       bool
	Env                map[string]string
	SubName            string
	Disable            bool
	Log                []string
	cmd                *exec.Cmd
	Status             *ServiceStatus
	Alert              map[string]alert.SendAlerter
	AT                 *internal.AlertTo
	Port               int
	ContinuityInterval time.Duration
	AI                 *alert.AlertInfo // 报警规则
	Exit               chan int         // 判断是否是主动退出的
	EndStop            chan bool
	Ctx                context.Context
	Cancel             context.CancelFunc
	Email              []string
	Msg                chan string
}

func (s *Script) cron() {
	s.Status.Status = RUNNING
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
func (s *Script) Start() error {
	switch s.Status.Status {
	case WAITSTOP:
		// 如果之前是等待停止的状态， 更改为重启状态
		<-s.Exit
		s.Exit <- 10
		s.Status.Status = WAITRESTART
	case STOP:
		s.Exit = make(chan int, 2)
		s.EndStop = make(chan bool, 2)
		s.Ctx, s.Cancel = context.WithCancel(context.Background())
		if s.Cron != nil && s.Cron.Loop > 0 {
			// 如果时间没填， 或者已经过去的时间了， 那么就直接启动
			if (s.Cron.Start != time.Time{}) && time.Since(s.Cron.Start) < 0 {
				go s.cron()
				return nil
			}
			s.loopTime = time.Now()
		}

		s.Status.Status = RUNNING
		if err := s.start(); err != nil {
			return err
		}

		go s.wait()
		if s.cmd.Process != nil {
			s.Status.Pid = s.cmd.Process.Pid
		}

	}
	return nil
}

// Restart  重动服务
func (s *Script) Restart() {
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

func (s Script) Remove() {
	switch SS.Infos[s.Name][s.SubName].Status.Status {
	case WAITRESTART, WAITSTOP:
		// 结束发送的退出错误发出的信号
		<-SS.Infos[s.Name][s.SubName].Exit
		// 结束停止的goroutine， 转为删除处理
		SS.Infos[s.Name][s.SubName].EndStop <- true
		go SS.Infos[s.Name][s.SubName].remove()
	case STOP:
		// 直接删除
		delete(SS.Infos[s.Name], s.SubName)
		if len(SS.Infos[s.Name]) == 0 {
			delete(SS.Infos, s.Name)
		}
	case RUNNING:
		go SS.Infos[s.Name][s.SubName].remove()
	default:
		golog.Error("error status")
	}
}

func (s *Script) remove() {
	s.Stop()
	delete(SS.Infos[s.Name], s.SubName)
	if len(SS.Infos[s.Name]) == 0 {
		delete(SS.Infos, s.Name)
	}
}

// Stop  停止服务
func (s *Script) Stop() {
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

// Stop  杀掉服务
func (s *Script) Kill() {
	switch s.Status.Status {
	case RUNNING:
		s.Exit <- 9
		s.kill()
	case WAITRESTART, WAITSTOP:
		<-s.Exit
		s.Exit <- 9
		s.kill()
	}

}

func (s *Script) wait() error {
	go s.successAlert()
	if err := s.cmd.Wait(); err != nil {
		golog.Info("error stop")
		s.Cancel()
		// 执行脚本后环境的错误
		select {
		case ec := <-s.Exit:
			switch ec {
			case 9:
				// 主动退出, kill， stop
				s.stopStatus()
				return nil
			case 10:
				// 重启 restart
				s.stopStatus()
				return s.Start()
			}
		default:
			// 意外退出
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
					if time.Since(s.AI.AlertTime) >= s.ContinuityInterval {
						s.AI.AlertTime = time.Now()
						alert.AlertMessage(am, s.AT)
					}
				}
			}
			if s.Always {
				golog.Info(time.Now())
				// 失败了， 每秒启动一次
				s.Status.RestartCount++
				return s.Start()

			}
		}
		golog.Debugf("serviceName: %s, subScript: %s, error: %v \n", s.Name, s.SubName, err)
		s.stopStatus()
		if s.Cron != nil && s.Cron.Loop > 0 {
			goto loop
		}
		// s.Status.Last = false
		return err
	}
loop:
	if s.Cron != nil && s.Cron.Loop > 0 {
		if s.Cron.IsMonth {
			start := time.Since(s.loopTime.AddDate(0, s.Cron.Loop, 0))
			if start < 0 {
				for {
					select {
					case <-s.Ctx.Done():
						golog.Info("service stop and loop have been cancel")
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
						golog.Info("service stop and loop have been cancel")
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
	s.stopStatus()
	return nil

}

func (s *Script) stopStatus() {
	s.Status.Status = STOP
	s.Status.RestartCount = 0
	s.Status.Pid = 0
	s.cmd = nil
}
