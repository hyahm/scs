package server

import (
	"context"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal"
	"github.com/hyahm/scs/internal/server/status"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/alert"
	"github.com/hyahm/scs/pkg/config/alert/to"
	"github.com/hyahm/scs/pkg/config/liveness"
	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/scs/pkg/config/scripts/cron"
	"github.com/hyahm/scs/pkg/config/scripts/prestart"
	"github.com/hyahm/scs/pkg/message"
)

// 默认的间隔时间
const defaultContinuityInterval = time.Hour * 1

type Server struct {
	Index      int               `json:"index"` // svc的索引
	Token      string            `json:"token"` // svc的token
	Name       string            `json:"name"`
	Dir        string            `json:"dir,omitempty"`
	Command    string            `json:"command"`
	Version    string            `json:"version,omitempty"`
	Cron       *cron.Cron        `json:"cron,omitempty"`    // 这个cron是新生成的
	IsCron     bool              `json:"is_loop,omitempty"` // 如果是定时任务
	Env        map[string]string `json:"-"`
	Logger     *golog.Log        `json:"-"`               // 日志
	Times      int               `json:"times,omitempty"` // 记录循环的次数
	SubName    string            `json:"subname,omitempty"`
	Cmd        *exec.Cmd         `json:"-"`
	AlwaysSign bool              `json:"always"` // 在停止的时候， always会变为false
	// 总副本数
	Replicate int            `json:"replicate,omitempty"`
	Status    *status.Status `json:"status,omitempty"`
	// Alert     map[string]message.SendAlerter `json:"-"`
	//  todo: 感觉不够完善
	AT      *to.AlertTo      `json:"at,omitempty"`
	Disable bool             `json:"disable,omitempty"`
	Port    int              `json:"port,omitempty"`
	AI      *alert.AlertInfo `json:"-"` // 报警规则
	// 主动退出的信号， kill: 9, restart: 10, stop: 11, remove: 12
	Exit chan int `json:"-"`
	// 取消操作， 可以取消等待重启， 等待停止， 等待remove(暂时没实现)
	CancelProcess chan bool `json:"-"`
	// 服务停止后的信号， 比如  restart, remove 操作， 因为停止后还有下一步操作
	StopSignal chan bool `json:"-"`
	// 这2个上上下文
	Ctx    context.Context    `json:"-"`
	Cancel context.CancelFunc `json:"-"` // 结束定时器的上下文和日志的上下文
	// 更新的命令
	Update string `json:"update,omitempty"`
	// 暂时无视
	Liveness *liveness.Liveness `json:"-"`
	Ready    chan bool          `json:"-"`
	// 是否一直重启， 应该还需要一个retry次数的字段才对
	Always bool `json:"always,omitempty"`
	// 取消报警的感觉没用， 谁没事了会取消报警
	DisableAlert bool `json:"disable_alert,omitempty"`
	// 启动前的准备工作
	PreStart []*prestart.PreStart `json:"-"`
	// 执行完成就自动删除
	DeleteWhenExit bool `json:"deleteWhenExit,omitempty"`
	// 执行完成就remove的信号
	DeleteWhenExitSingle chan bool `json:"-"`
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
	cmd := newCommand(command)
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

// Restart send stop single
func (svc *Server) Restart() {

	if svc.IsCron {
		svc.Cancel()
		// 如果是循环的就直接退出
		return
	}
	if svc.Always {
		svc.Always = false
	}
	golog.Info(svc.Status.Status)
	switch svc.Status.Status {
	case status.WAITSTOP:
		// 如果之前是等待停止的状态， 更改为重启状态
		golog.Info("waiting stop")
		<-svc.Exit
		svc.Exit <- 10
		svc.Status.Status = status.WAITRESTART
		return
	case status.RUNNING:
		svc.Exit <- 10
		svc.Status.Status = status.WAITRESTART
		go svc.stop()
		return
	case status.STOP:
		golog.Debug("ready send stop single")
		svc.StopSignal <- true
		golog.Debug("send stop single")
	}

}

func (svc *Server) MakeServer(script *scripts.Script, availablePort int) int {
	// 将环境变量填充到server中
	script.MakeEnv()
	env := make(map[string]string)
	for k, v := range script.TempEnv {
		env[k] = v
	}
	if script.Port > 0 {
		// 顺序拿到可用端口
		availablePort = pkg.GetAvailablePort(availablePort)
		env["PORT"] = strconv.Itoa(availablePort)
	} else {
		env["PORT"] = "0"
	}
	svc.fillServer(script)
	env["OS"] = runtime.GOOS
	env["NAME"] = svc.SubName
	env["PROJECT_HOME"] = svc.Dir
	// 格式化 SCS_TPL 开头的环境变量
	for k := range env {
		if len(k) > 8 && k[:7] == "SCS_TPL" {
			env[k] = internal.Format(env[k], env)
		}
	}
	svc.Env = env
	svc.Port = availablePort
	return availablePort
}

func (svc *Server) fillServer(script *scripts.Script) {
	// 填充server
	svc.Token = script.Token
	svc.Command = script.Command
	svc.Disable = script.Disable
	// Log:       make([]string, 0, global.GetLogCount()),
	svc.Dir = script.Dir
	if svc.Status == nil {
		svc.Status = &status.Status{
			Status: status.STOP,
		}
	}

	svc.Logger = golog.NewLog(
		filepath.Join(global.LogDir, svc.SubName+".log"), 10<<10, false, global.CleanLog)
	svc.Update = script.Update
	svc.AI = &alert.AlertInfo{}
	svc.AT = script.AT
	svc.StopSignal = make(chan bool, 1)

	svc.Liveness = script.Liveness
	svc.Ready = make(chan bool, 1)
	svc.Always = script.Always
	svc.AlwaysSign = script.Always
	svc.DeleteWhenExit = script.DeleteWhenExit

	// svc.DisableAlert = script.DisableAlert
	svc.PreStart = script.PreStart

	svc.Logger.Format = global.FORMAT
	if script.Cron != nil {
		svc.Cron = &cron.Cron{
			Start:   script.Cron.Start,
			Loop:    script.Cron.Loop,
			IsMonth: script.Cron.IsMonth,
			Times:   script.Cron.Times,
		}
	}
}

// 同步删除
func (svc *Server) Remove() {
	defer svc.Cancel()
	if svc.IsCron {
		svc.Cancel()
		return
	}
	if svc.Always {
		svc.Always = false
	}
	switch svc.Status.Status {
	case status.WAITRESTART:
		// 结束发送的退出错误发出的信号
		<-svc.Exit
		// 结束停止的goroutine， 转为删除处理
		svc.Exit <- 12
		svc.Stop()
	case status.STOP:
		golog.Debug("ready send stop single")
		svc.StopSignal <- true
		// DeleteServiceBySubName(svc.SubName)
	case status.RUNNING:
		svc.Exit <- 12
		svc.Stop()
	case status.WAITSTOP:
		<-svc.Exit
		// 结束停止的goroutine， 转为删除处理
		svc.Exit <- 12
	default:
		golog.Error("error status")
	}
}

// Stop  停止服务
func (svc *Server) Stop() {
	if svc.Disable {
		return
	}
	if svc.IsCron {
		// 如果是定时任务， 直接停止
		golog.Infof("stop loop %s", svc.SubName)
		svc.Cancel()
	}
	if svc.Always {
		svc.Always = false
	}
	switch svc.Status.Status {
	case status.RUNNING:
		svc.Exit <- 9
		svc.Status.Status = status.WAITSTOP
		svc.stop()
	// case status.STOP:
	// svc.Exit <- 9
	case status.WAITRESTART:
		// 将退出信号设置为waiting stop
		<-svc.Exit
		svc.Exit <- 9
		svc.Status.Status = status.WAITSTOP
	}
}

// 同步更新并重启
func (svc *Server) UpdateServer() {
	updateCommand := "git pull"
	if svc.Update != "" {
		updateCommand = svc.Update
	}
	if err := svc.shell(updateCommand); err != nil {
		golog.Error(err)
		return
	}

}

// Stop  杀掉服务, 没有产生goroutine， 直接杀死
func (svc *Server) Kill() {
	if svc.IsCron {
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
	golog.UpFunc(1, "stop")
	svc.Status.Status = status.STOP
	svc.Status.Pid = 0
	svc.Status.CanNotStop = false
	svc.Status.RestartCount = 0
	svc.Status.Start = 0
	svc.Logger.Close()
	svc.Cmd = nil
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
				Name:       s.SubName,
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
