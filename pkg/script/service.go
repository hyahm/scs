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

var SS map[string]map[string]*Script

func init() {
	SS = make(map[string]map[string]*Script)
}

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
	exit               bool             // 判断是否是主动退出的
	Ctx                context.Context
	Cancel             context.CancelFunc
	Email              []string
	KillTime           time.Duration
	Msg                chan string
	// exitCode chan int // 如果推出信号是9
}

// Start  启动服务
func (s *Script) Start() {
	s.Ctx, s.Cancel = context.WithCancel(context.Background())
	for {
		select {
		case <-s.Ctx.Done():
			golog.Infof("stop service : %s", s.SubName)
			return
		}
	}
}

// Restart  重动服务
func (s *Script) Restart() {
	for {
		select {
		case <-s.Ctx.Done():
			golog.Infof("stop service : %s", s.SubName)
			return
		}
	}
}

// Stop  停止服务
func (s *Script) Stop() {
	for {
		select {
		case <-s.Ctx.Done():
			golog.Infof("stop service : %s", s.SubName)
			return
		}
	}
}

func (s *Script) waitCmd() chan error {
	err := make(chan error)
	err <- s.cmd.Wait()
	return err
}

func (s *Script) wait() error {
	go s.successAlert()
	if err := s.cmd.Wait(); err != nil {
		// 执行脚本后环境的错误
		s.cmd = nil
		golog.Info("time error")
		time.Sleep(1 * time.Second)
		if !s.exit && !s.DisableAlert {
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
		golog.Debugf("serviceName: %s, subScript: %s, error: %v \n", s.Name, s.SubName, err)
		s.stopStatus()
		golog.Info(s.exit)
		if !s.exit && s.Loop > 0 {
			goto loop
		}
		if !s.exit && s.Always {
			golog.Info(time.Now())
			// 失败了， 每秒启动一次
			golog.Info("restart")
			time.Sleep(s.KillTime)
			s.Status.RestartCount++
			s.Start()
		}
		s.exit = false
		// s.Status.Last = false
		return err
	}
loop:
	if s.Loop > 0 {
		sleep := math.Ceil(float64(s.Loop) - time.Now().Sub(s.loopTime).Seconds())
		if sleep > 0 {
			// 允许循环， 每s.Loop秒启动一次
			time.Sleep(time.Duration(sleep) * time.Second)
		}

		s.Start()
		s.exit = false
		return nil
	}
	s.stopStatus()
	return nil

}

func (s *Script) stopStatus() {
	s.Status.Status = STOP
	s.Status.RestartCount = 0
	s.Status.Ppid = 0
	s.Status.Up = 0
}
