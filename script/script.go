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

// 脚本的信息
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
	Exit               chan bool
	ctx                context.Context
	cancel             context.CancelFunc
	Email              []string
	KillTime           time.Duration
	// exitCode chan int // 如果推出信号是9
}

func (s *Script) Restart() {
	s.Status.Status = WAITRESTART
	s.Stop()
	// 先要停止， 然后再启动
	// 判断是否已经停止了
	for {
		if s.Status.Status == STOP {
			break
		}
		time.Sleep(s.KillTime)
	}
	s.exit = false
	s.Start()
}

func (s *Script) GetEnv() []string {
	return s.cmd.Env
}

func (s *Script) wait() error {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	go s.successAlert()

	// 这个goroutine 随 kill() 关闭
	if err := s.cmd.Wait(); err != nil {
		// 执行脚本后环境的错误
		s.cmd = nil
		golog.Info("time error")
		time.Sleep(1 * time.Second)
		s.cancel()
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
			// 失败了， 每秒启动一次
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
		select {
		case <-time.After(time.Duration(sleep) * time.Second):
			golog.Info("%s have been loop at %v", s.Name, time.Now())
			s.Start()
			s.exit = false
			return nil
		case <-s.Exit:
		}

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

func (s *Script) Install(command string) {
	s.exit = false
	golog.Info(s.Log)
	s.start()

	go s.waitinstall()
}

func (s *Script) waitinstall() {

	// 这个goroutine 随 kill() 关闭
	if err := s.cmd.Wait(); err != nil {
		// 执行脚本后环境的错误
		s.cmd = nil
		time.Sleep(1 * time.Second)

		return
	}
	// s.stopStatus()

}
