package script

import (
	"context"
	"math"
	"os/exec"
	"scs/alert"
	"scs/internal"
	"time"

	"github.com/hyahm/golog"
)

type Script struct {
	LookPath           []*internal.LoopPath
	Name               string
	Dir                string
	Command            string
	Replicate          int
	Always             bool
	Loop               int
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
	KillTime           time.Duration
	Msg                chan string
}

// Start  启动服务
func (s *Script) Start() error {
	s.Exit = make(chan int, 2)
	s.EndStop = make(chan bool, 2)
	s.Ctx, s.Cancel = context.WithCancel(context.Background())
	golog.Info("start")
	if s.Loop > 0 {
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
	return nil
}

// Restart  重动服务
func (s *Script) Restart() {
	s.Exit <- 10
	s.stop()
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
	s.Exit <- 9
	s.stop()
}

// Stop  杀掉服务
func (s *Script) Kill() {
	s.Exit <- 9
	s.kill()
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
				time.Sleep(s.KillTime)
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
				time.Sleep(s.KillTime)
				s.Status.RestartCount++
				return s.Start()

			}
		}
		golog.Debugf("serviceName: %s, subScript: %s, error: %v \n", s.Name, s.SubName, err)
		s.stopStatus()
		// s.Status.Last = false
		return err
	}
	if s.Loop > 0 {
		sleep := math.Ceil(float64(s.Loop) - time.Now().Sub(s.loopTime).Seconds())
		if sleep > 0 {
			// 允许循环， 每s.Loop秒启动一次
			time.Sleep(time.Duration(sleep) * time.Second)
		}
		s.stopStatus()
		s.Start()
		return nil
	}
	s.stopStatus()
	return nil

}

func (s *Script) stopStatus() {
	s.Status.Status = STOP
	s.Status.RestartCount = 0
	s.Status.Pid = 0
	s.cmd = nil
	s.Status.Up = 0
}
