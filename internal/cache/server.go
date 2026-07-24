package cache

import (
	"context"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal"

	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
)

// 默认的间隔时间
const defaultContinuityInterval = time.Hour * 1

type Server struct {
	mu             sync.RWMutex
	canNotOpration bool
	Index          int `json:"index"` // svc的索引
	// ScriptToken string            `json:"scriptToken"` // svc的token
	// SimpleToken string            `json:"simpleToken"` // svc的token
	// User       string            `json:"user"`
	// Group      string            `json:"group"`
	Name       string            `json:"name"`
	Dir        string            `json:"dir,omitempty"`
	Command    string            `json:"command"`
	Version    string            `json:"version,omitempty"`
	Cron       config.Cron       `json:"cron,omitempty"`    // 这个cron是新生成的
	IsCron     bool              `json:"is_loop,omitempty"` // 如果是定时任务
	Env        map[string]string `json:"-"`
	Logger     *golog.Log        `json:"-"`               // 日志
	Times      int               `json:"times,omitempty"` // 记录循环的次数
	SubName    string            `json:"subname,omitempty"`
	Cmd        *exec.Cmd         `json:"-"`
	AlwaysSign bool              `json:"-"` // 在停止的时候， always会变为false
	// StartTime  string            `json:"-"`
	// StopTime   string            `json:"-"`
	// 总副本数
	Replicate int        `json:"replicate,omitempty"`
	Status    pkg.Status `json:"status,omitempty"`
	// Alert     map[string]message.SendAlerter `json:"-"`
	//  todo: 感觉不够完善
	// AT      *config.AlertTo   `json:"at,omitempty"`
	Disable bool `json:"disable,omitempty"`
	Port    int  `json:"port,omitempty"`
	// AI      *config.AlertInfo `json:"-"` // 报警规则
	// 进程退出后是否自动重启（Restart 设置，wait 消费）
	restartOnExit bool `json:"-"`
	// 取消操作， 可以取消等待重启， 等待停止， 等待remove(暂时没实现)
	// CancelProcess chan bool `json:"-"`
	// 服务停止后的信号， 比如  restart, remove 操作， 因为停止后还有下一步操作
	// StopSignal chan bool `json:"-"`
	// 这2个上上下文
	Ctx    context.Context    `json:"-"`
	Cancel context.CancelFunc `json:"-"` // 结束定时器的上下文和日志的上下文
	// 更新的命令
	Update string `json:"update,omitempty"`
	// 暂时无视
	Liveness *config.Liveness `json:"-"`
	Ready    chan bool        `json:"-"`
	// 是否一直重启， 应该还需要一个retry次数的字段才对
	Always bool `json:"always,omitempty"`
	// 取消报警的感觉没用， 谁没事了会取消报警
	DisableAlert bool `json:"disable_alert,omitempty"`
	// 启动前的准备工作
	PreStart []config.PreStart `json:"-"`
	// 执行完成就自动删除
	DeleteWhenExit bool `json:"deleteWhenExit,omitempty"`
	// 执行完成就remove的信号
	DeleteWhenExitSingle chan bool `json:"-"`
}

func (s *Server) SetCanNotOperation(v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.canNotOpration = v
}

// CanOperation 安全地读取 canOpration
func (s *Server) GetCanNotOperation() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.canNotOpration
}

func (s *Server) SetCanNotStop(v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Status.CanNotStop = v
}

// CanOperation 安全地读取 canOpration
func (s *Server) GetCanNotStop() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Status.CanNotStop
}

func (s *Server) setRestartOnExit(v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.restartOnExit = v
}

func (s *Server) getRestartOnExit() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.restartOnExit
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

func (svc *Server) GetStatus() pkg.ServiceStatus {
	svc.mu.RLock()
	defer svc.mu.RUnlock()
	status := pkg.ServiceStatus{
		PName:        svc.Name,
		Name:         svc.SubName,
		IsCron:       svc.IsCron,
		Command:      svc.Status.Command,
		Version:      svc.Status.Version,
		CanNotStop:   svc.Status.CanNotStop,
		Path:         svc.Dir,
		Status:       svc.Status.Status,
		RestartCount: svc.Status.RestartCount,
		Pid:          svc.Status.Pid,
		Disable:      svc.Disable,
		OS:           runtime.GOOS,
		Start:        svc.Status.Start,
	}

	status.Cpu, status.Mem, _ = config.GetProcessInfo(int32(status.Pid))
	return status
}

// Restart 重启服务：先等 cannotstop 解除，再按当前状态决定——
// STOP 直接 StartAsync 拉起；RUNNING 置 restartOnExit 并 Cancel，由 wait() 在进程退出后重启。
func (svc *Server) Restart() {
	if svc.IsCron {
		// 定时任务直接取消循环
		svc.Cancel()
		return
	}
	// 等 cannotstop 解除
	if svc.GetCanNotStop() {
		<-svc.Status.ChStop
	}
	if svc.Always {
		svc.Always = false
	}
	svc.mu.Lock()
	switch svc.Status.Status {
	case pkg.STOP:
		svc.mu.Unlock()
		// 已停止，直接拉起
		svc.StartAsync()
	case pkg.RUNNING:
		// 标记本轮退出后重启，Cancel 触发 wait() 执行重启
		svc.restartOnExit = true
		svc.mu.Unlock()
		svc.Cancel()
	default:
		// WAITSTOP / WAITRESTART 等中间态：忽略，避免状态错乱
		svc.mu.Unlock()
	}
}

// 填充到server中
func (svc *Server) MakeServer(script config.Script) {

	env := make(map[string]string)
	for k, v := range script.TempEnv {
		env[k] = v
	}
	if svc.Port > 0 {
		// 顺序拿到可用端口
		svc.Port = pkg.GetAvailablePort(svc.Port)
		env["PORT"] = strconv.Itoa(svc.Port)
	} else {
		env["PORT"] = "0"
	}
	// 填充server
	svc.FillServer(script)
	env["NAME"] = svc.SubName

	svc.Env = env
}

// 填充server
func (svc *Server) FillServer(script config.Script) {

	// svc.ScriptToken = script.ScriptToken
	// svc.SimpleToken = script.SimpleToken
	// if svc.SimpleToken == "" {
	// 	svc.SimpleToken = pkg.RandomToken()
	// }
	svc.Command = script.Command
	// svc.User = script.User
	// svc.Group = script.Group
	svc.Disable = script.Disable
	// Log:       make([]string, 0, global.GetLogCount()),
	svc.Dir = script.Dir
	svc.Status = pkg.Status{
		Status: pkg.STOP,
	}
	// svc.StartTime = script.StartTime
	// svc.StopTime = script.StopTime

	svc.Update = script.Update
	// svc.AI = &config.AlertInfo{}
	// svc.AT = script.AT
	// svc.StopSignal = make(chan bool, 1)

	// svc.Liveness = script.Liveness
	svc.Ready = make(chan bool, 1)
	svc.Always = script.Always
	svc.AlwaysSign = script.Always
	svc.DeleteWhenExit = script.DeleteWhenExit

	// svc.DisableAlert = script.DisableAlert
	svc.PreStart = script.PreStart
	svc.Cron = script.Cron
	svc.Cron = config.Cron{
		Start:   script.Cron.Start,
		Loop:    script.Cron.Loop,
		IsMonth: script.Cron.IsMonth,
		Times:   script.Cron.Times,
	}
}

// Remove 删除服务：撤销待重启意图，等 cannotstop 解除后停止进程，再从全局 store 移除。
func (svc *Server) Remove() {
	// 撤销任何待重启意图
	svc.setRestartOnExit(false)
	if svc.IsCron {
		// 定时任务直接取消循环
		golog.Infof("stop loop %s", svc.SubName)
		svc.Cancel()
	} else {
		if svc.Always {
			svc.Always = false
		}
		if svc.GetCanNotStop() {
			<-svc.Status.ChStop
		}
		switch svc.Status.Status {
		case pkg.RUNNING, pkg.WAITRESTART, pkg.WAITSTOP:
			svc.Cancel()
			if err := svc.kill(); err != nil {
				golog.Error(err)
			}
		}
	}
	// 从全局 store 删除
	RemoveServer(svc.SubName)
}

// Stop  停止服务：撤销待重启意图后，等 cannotstop 解除再 Cancel。
func (svc *Server) Stop() {
	// 撤销任何待重启意图，避免进程退出后被 wait() 自动拉起
	svc.setRestartOnExit(false)
	if svc.Disable {
		return
	}

	if svc.GetCanNotStop() {
		<-svc.Status.ChStop
	}
	if svc.IsCron {
		// 定时任务直接取消循环
		svc.Cancel()
		return
	}
	if svc.Always {
		svc.Always = false
	}

	switch svc.Status.Status {
	case pkg.RUNNING:
		svc.Cancel()
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

// Kill 杀掉服务：撤销待重启意图，等 cannotstop 解除后 Cancel + kill 强杀整个进程组。
func (svc *Server) Kill() {
	if svc.IsCron {
		svc.Cancel()
		return
	}
	// 撤销任何待重启意图，避免退出后被 wait() 自动拉起
	svc.setRestartOnExit(false)
	if svc.GetCanNotStop() {
		<-svc.Status.ChStop
	}
	switch svc.Status.Status {
	case pkg.RUNNING, pkg.WAITRESTART, pkg.WAITSTOP:
		svc.Cancel()
		if err := svc.kill(); err != nil {
			golog.Error(err)
		}
	}
}

func (svc *Server) stopStatus() {
	svc.Status.Status = pkg.STOP
	svc.Status.Pid = 0
	svc.Status.CanNotStop = false
	svc.Status.RestartCount = 0
	svc.SetCanNotOperation(false)
	svc.Status.Start = 0
	svc.Cmd = nil
	// if svc.Logger != nil {
	// 	svc.Logger.Sync()
	// 	svc.Logger = nil
	// }
}

func (s *Server) successAlert() {
	// 启动成功后恢复的通知
	// if !s.AI.Broken {
	// 	return
	// }
	for {
		select {
		// 每3秒一次操作
		case <-time.After(time.Second * 3):
			// am := &message.Message{
			// 	Title: "service recover",
			// 	Pname: s.Name,
			// 	Name:  s.SubName,
			// 	// BrokenTime: s.AI.Start.String(),
			// 	FixTime: time.Now().String(),
			// }
			// config.AlertMessage(am, s.AT)
			// s.AI.Broken = false
			return
		case <-s.Ctx.Done():
			return
		}
	}

}
