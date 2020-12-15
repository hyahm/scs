// +build !aix,!darwin,!linux,!freebsd,!openbsd,!windows

// package script

// import (
// 	"context"
// 	"log"
// 	"os/exec"
// 	"runtime"
// 	"scs/alert"
// 	"scs/hardware"
// 	"scs/internal"
// 	"syscall"
// 	"time"

// 	"github.com/hyahm/golog"
// )

// // 脚本的信息
// type Script struct {
// 	Name               string
// 	Command            string
// 	Dir                string
// 	Replicate          int
// 	Always             bool
// 	DisableAlert       bool
// 	Env                []string
// 	SubName            string
// 	Log                []string
// 	cmd                *exec.Cmd
// 	Status             *ServiceStatus
// 	Alert              map[string]alert.SendAlerter
// 	AT                 internal.AlertTo
// 	Port               int
// 	ContinuityInterval time.Duration
// 	AI                 *hardware.AlertInfo // 报警规则
// 	exit               bool                // 判断是否是主动退出的
// 	ctx                context.Context
// 	cancel             context.CancelFunc
// 	Email              []string
// 	KillTime           time.Duration

// 	// exitCode chan int // 如果推出信号是9
// }

// func (s *Script) Stop() {
// 	if s.Status.Status == RUNNING {
// 		s.Status.Status = WAITSTOP
// 	}
// 	defer func() {
// 		if err := recover(); err != nil {
// 			golog.Info("脚本已经停止了")
// 		}
// 	}()
// 	for {
// 		time.Sleep(time.Millisecond * 10)
// 		if !s.Status.CanNotStop {
// 			s.exit = true
// 			s.cancel()
// 			err := syscall.Kill(-s.cmd.Process.Pid, syscall.SIGKILL)
// 			if err != nil {
// 				// 正常来说，不会进来的，特殊问题以后再说
// 				golog.Error(err)
// 			}

// 			golog.Infof("stop %s\n", s.SubName)

// 			// 预留3秒代码退出的时间
// 			time.Sleep(s.KillTime)
// 			return
// 		}
// 	}

// }

// func (s *Script) Kill() {
// 	// 数组存日志
// 	// s.Log = make([]string, Config.LogCount)
// 	// s.cancel()
// 	s.exit = true
// 	var err error

// 	err = syscall.Kill(-s.cmd.Process.Pid, syscall.SIGKILL)
// 	// err = s.cmd.Process.Kill()
// 	// err = exec.Command("kill", "-9", fmt.Sprint(s.cmd.Process.Pid)).Run()
// 	// err := s.cmd.Process.Kill()

// 	if err != nil {
// 		// 正常来说，不会进来的，特殊问题以后再说
// 		golog.Error(err)
// 		// return
// 	}
// 	s.stopStatus()

// 	return

// }

// func (s *Script) Restart() {
// 	s.Status.Status = WAITRESTART
// 	s.Stop()
// 	// 先要停止， 然后再启动
// 	// 判断是否已经停止了
// 	for {
// 		if s.Status.Status == STOP {
// 			break
// 		}
// 		time.Sleep(s.KillTime)
// 	}
// 	s.exit = false
// 	s.Start()
// }

// func (s *Script) stopStatus() {
// 	s.Status.Status = STOP
// 	s.Status.Ppid = 0
// 	s.Status.Up = 0
// }

// func (s *Script) Start() {
// 	s.exit = false
// 	golog.Info(s.Status.Version)
// 	s.Status.Status = RUNNING
// 	// index := strings.Index(s.Command, " ")
// 	// s.cmd = exec.Command(s.Command[:index], s.Command[index:])
// 	s.Status.Command = s.Command
// 	if runtime.GOOS == "windows" {
// 		log.Fatal("not support windows")
// 	} else {
// 		s.cmd = exec.Command("/bin/bash", "-c", s.Command)
// 	}
// 	s.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
// 	s.cmd.Dir = s.Dir
// 	s.cmd.Env = s.Env
// 	// 等待初始化完成完成后向后执行
// 	s.read()
// 	s.Status.Up = time.Now().Unix() // 设置启动状态是成功的
// 	if err := s.cmd.Start(); err != nil {
// 		// 执行脚本前的错误, 改变状态
// 		golog.Error(err)
// 		s.stopStatus()
// 		return
// 	}

// 	if s.cmd.Process == nil {
// 		s.stopStatus()
// 		return
// 	}

// 	s.Status.Ppid = s.cmd.Process.Pid
// 	go s.wait()
// }

// func (s *Script) wait() {
// 	s.ctx, s.cancel = context.WithCancel(context.Background())
// 	go s.successAlert()

// 	// 这个goroutine 随 kill() 关闭
// 	if err := s.cmd.Wait(); err != nil {
// 		// 执行脚本后环境的错误
// 		s.cmd = nil
// 		time.Sleep(1 * time.Second)
// 		s.cancel()
// 		if !s.exit && !s.DisableAlert {
// 			am := &alert.Message{
// 				Title:      "service error stop",
// 				Pname:      s.Name,
// 				Name:       s.SubName,
// 				Reason:     err.Error(),
// 				BrokenTime: s.AI.Start.String(),
// 			}
// 			if !s.AI.Broken {
// 				// 第一次
// 				s.AI.Start = time.Now()
// 				s.AI.Interval = time.Now()
// 				s.AI.Broken = true
// 				alert.AlertMessage(am, &s.AT)
// 			} else {
// 				// 间隔时间内才发送报警
// 				if time.Since(s.AI.Interval) >= s.ContinuityInterval {
// 					s.AI.Interval = time.Now()
// 					alert.AlertMessage(am, &s.AT)
// 				}
// 			}
// 		}
// 		golog.Debugf("serviceName: %s, subScript: %s, error: %v \n", s.Name, s.SubName, err)
// 		s.stopStatus()

// 		if !s.exit && s.Always {
// 			// 失败了， 每秒启动一次
// 			golog.Info("restart")
// 			time.Sleep(1 * time.Second)
// 			s.Status.RestartCount++
// 			s.Start()
// 		}
// 		s.exit = false
// 		// s.Status.Last = false
// 		return
// 	}
// 	s.stopStatus()

// }
