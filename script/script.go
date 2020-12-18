package script

import (
	"context"
	"os"
	"os/exec"
	"scs/alert"
	"scs/internal"
	"time"

	"github.com/hyahm/golog"
)

// 脚本的信息
type Script struct {
	GetIfNotExist      string
	Name               string
	Dir                string
	Command            string
	Replicate          int
	Always             bool
	DisableAlert       bool
	Env                map[string]string
	SubName            string
	Log                []string
	cmd                *exec.Cmd
	Status             *ServiceStatus
	Alert              map[string]alert.SendAlerter
	AT                 internal.AlertTo
	Port               int
	ContinuityInterval time.Duration
	AI                 *alert.AlertInfo // 报警规则
	exit               bool             // 判断是否是主动退出的
	ctx                context.Context
	cancel             context.CancelFunc
	Email              []string
	KillTime           time.Duration
	IsScript           bool
	// exitCode chan int // 如果推出信号是9
}

func (s *Script) RunGetResource() error {
	if s.GetIfNotExist != "" && s.Dir != "" {
		if _, err := os.Open(s.Dir); os.IsNotExist(err) {
			s.IsScript = true
			defer func() {
				s.IsScript = false
			}()
			return s.Start(s.GetIfNotExist)
		}
	}
	return nil
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
	s.Start(s.Command)
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
				alert.AlertMessage(am, &s.AT)
			} else {
				// 间隔时间内才发送报警
				if time.Since(s.AI.AlertTime) >= s.ContinuityInterval {
					s.AI.AlertTime = time.Now()
					alert.AlertMessage(am, &s.AT)
				}
			}
		}
		golog.Debugf("serviceName: %s, subScript: %s, error: %v \n", s.Name, s.SubName, err)
		s.stopStatus()

		if !s.exit && s.Always && !s.IsScript {
			// 失败了， 每秒启动一次
			golog.Info("restart")
			time.Sleep(1 * time.Second)
			s.Status.RestartCount++
			s.Start(s.Command)
		}
		s.exit = false
		// s.Status.Last = false
		return err
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
	// index := strings.Index(s.Command, " ")
	// s.cmd = exec.Command(s.Command[:index], s.Command[index:])
	s.start(command)

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
