package server

import (
	"context"
	"os/exec"
	"runtime"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/internal/config/alert"
	"github.com/hyahm/scs/internal/config/alert/to"
	"github.com/hyahm/scs/internal/config/liveness"
	"github.com/hyahm/scs/internal/config/scripts/cron"
	"github.com/hyahm/scs/internal/config/scripts/prestart"
	"github.com/hyahm/scs/internal/config/scripts/subname"
	"github.com/hyahm/scs/pkg/message"
	"github.com/hyahm/scs/status"
)

// 默认的间隔时间
const defaultContinuityInterval = time.Hour * 1

type Server struct {
	Index              int                            `json:"index"` // svc的索引
	Name               string                         `json:"name"`
	Dir                string                         `json:"dir,omitempty"`
	Command            string                         `json:"command"`
	Version            string                         `json:"version,omitempty"`
	Cron               *cron.Cron                     `json:"cron,omitempty"`    // 这个cron是新生成的
	IsLoop             bool                           `json:"is_loop,omitempty"` // 如果是定时任务
	Env                map[string]string              `json:"-"`
	Logger             *golog.Log                     `json:"-"`               // 日志
	Times              int                            `json:"times,omitempty"` // 记录循环的次数
	SubName            subname.Subname                `json:"subname,omitempty"`
	Cmd                *exec.Cmd                      `json:"-"`
	Replicate          int                            `json:"replicate,omitempty"`
	Status             *status.ServiceStatus          `json:"status,omitempty"`
	Alert              map[string]message.SendAlerter `json:"-"`
	AT                 *to.AlertTo                    `json:"at,omitempty"`
	Disable            bool                           `json:"disable,omitempty"`
	Port               int                            `json:"port,omitempty"`
	ContinuityInterval time.Duration                  `json:"continuity_interval,omitempty"`
	AI                 *alert.AlertInfo               `json:"-"` // 报警规则
	Exit               chan int                       `json:"-"` // 判断是否是主动退出的
	CancelProcess      chan bool                      `json:"-"` // 取消操作，
	// 停止后发出的信号, 9 主动退出， 10 重启， 11 主动退出并删除
	StopSigle    chan bool            `json:"-"`
	Ctx          context.Context      `json:"-"`
	Cancel       context.CancelFunc   `json:"-"`                // 结束定时器的上下文和日志的上下文
	Removed      bool                 `json:"-"`                // 标识是否已经被删除
	Update       string               `json:"update,omitempty"` // 更新的命令
	Liveness     *liveness.Liveness   `json:"-"`
	Ready        chan bool            `json:"-"`
	Always       bool                 `json:"always,omitempty"`
	DisableAlert bool                 `json:"disable_alert,omitempty"`
	PreStart     []*prestart.PreStart `json:"-"`
}

func newCommand(command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("powershell", "/C", command)
	} else {
		return exec.Command("/bin/bash", "-c", command)
	}
}

// update 的时候执行
func (svc *Server) shell(command string) error {
	var cmd *exec.Cmd

	command = internal.Format(command, svc.Env)
	cmd = newCommand(command)
	golog.Info(command)
	for k, v := range svc.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	cmd.Dir = svc.Dir

	read(cmd, svc)
	err := cmd.Start()
	if err != nil {
		golog.Error(err)
		return err
	}
	return cmd.Wait()
}

func (svc *Server) shellWithOutDir(command string) error {
	var cmd *exec.Cmd
	cmd = newCommand(command)
	for k, v := range svc.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	cmd.Dir = svc.Dir
	read(cmd, svc)
	err := cmd.Start()
	if err != nil {
		golog.Error(err)
		return err
	}
	return cmd.Wait()
}

// 这是未开发出的就绪状态
func (svc *Server) CheckReady(ctx context.Context) {
	if svc.Liveness == nil || svc.Liveness.Http != "" && svc.Liveness.Tcp != "" && svc.Liveness.Shell != "" {
		//
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Millisecond * 1):
			ok := svc.Liveness.Ready()
			if ok {
				svc.Ready <- true
				return
			}
		}
	}

}

// Restart  重动服务, 同步执行的
func (svc *Server) Restart() {
	if svc.IsLoop {
		svc.Cancel()
		// 如果是循环的就直接退出
		return
	}
	switch svc.Status.Status {
	case status.WAITSTOP:
		// 如果之前是等待停止的状态， 更改为重启状态
		<-svc.Exit
		svc.Exit <- 10
		svc.Status.Status = status.WAITRESTART
		return
	case status.RUNNING:
		svc.Exit <- 10
		svc.Status.Status = status.WAITRESTART
		svc.stop()
		return
	case status.STOP:
		svc.Start()
	}

}

// 同步删除
func (svc *Server) Remove() {
	defer svc.Cancel()
	if svc.IsLoop {
		svc.Cancel()
		return
	}
	switch svc.Status.Status {
	case status.WAITRESTART:
		// 结束发送的退出错误发出的信号
		<-svc.Exit
		// 结束停止的goroutine， 转为删除处理
		svc.CancelProcess <- true
		svc.Stop()
	case status.STOP:
		// svc.StopSigle <- true
	case status.INSTALL:
		// TODO 直接删除
		// svc.Stop()
		// go svc.remove()
		// DeleteServiceBySubName(svc.SubName)
	case status.RUNNING, status.WAITSTOP:
		svc.Stop()
	default:
		golog.Error("error status")
	}
}

// Stop  停止服务
func (svc *Server) Stop() {
	if svc.IsLoop {
		// 如果是定时任务， 直接停止
		golog.Infof("stop loop %s", svc.SubName)
		svc.Cancel()
	}
	switch svc.Status.Status {
	case status.RUNNING:
		svc.Exit <- 9
		svc.Status.Status = status.WAITSTOP
		svc.stop()
	case status.STOP:
		svc.Exit <- 9
	case status.WAITRESTART:
		<-svc.Exit
		svc.Exit <- 9
		svc.Status.Status = status.WAITSTOP
	}
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
		return
	}
	switch svc.Status.Status {
	case status.RUNNING:
		svc.Exit <- 9
		if err := svc.kill(); err != nil {
			golog.Error(err)
			// s.Cancel()
		}
	case status.WAITRESTART, status.WAITSTOP:
		<-svc.Exit
		svc.Exit <- 9
		svc.kill()
	}

}

func (svc *Server) stopStatus() {
	svc.Status.Status = status.STOP
	svc.Status.Pid = 0
	svc.Status.CanNotStop = false
	svc.Status.RestartCount = 0
	svc.Status.Start = 0
	// svc.Removed = false
}

func (s *Server) successAlert() {
	// 启动成功后恢复的通知
	if !s.AI.Broken {
		return
	}
	for {
		select {
		// 每3秒一次操作
		case <-time.After(time.Second * 3):
			am := &message.Message{
				Title:      "service recover",
				Pname:      s.Name,
				Name:       s.SubName.String(),
				BrokenTime: s.AI.Start.String(),
				FixTime:    time.Now().String(),
			}
			alert.AlertMessage(am, s.AT)
			s.AI.Broken = false
			return
		case <-s.Ctx.Done():
			return
		}
	}

}
