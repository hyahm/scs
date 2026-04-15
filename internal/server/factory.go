package server

import (
	"path/filepath"
	"strconv"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
	"github.com/hyahm/scs/internal/server/status"
	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config/alert"
	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/scs/pkg/config/scripts/cron"
)

type ServerOption func(s *Server)

func WithIndex(index int) ServerOption {
	return func(s *Server) { s.Index = index }
}

func WithName(name string) ServerOption {
	return func(s *Server) { s.Name = name }
}

func WithSubName(subname string) ServerOption {
	return func(s *Server) { s.SubName = subname }
}

func WithPort(port int) ServerOption {
	return func(s *Server) { s.Port = port }
}

func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		Status: &status.Status{
			Status: status.STOP,
		},
		StopSignal: make(chan bool, 1),
		Ready:      make(chan bool, 1),
		Env:        make(map[string]string),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) ApplyConfig(script *scripts.Script) {
	s.ScriptToken = script.ScriptToken
	if s.SimpleToken == "" {
		s.SimpleToken = pkg.RandomToken()
	}
	s.Command = script.Command
	s.User = script.User
	s.Group = script.Group
	s.Disable = script.Disable
	s.Dir = script.Dir
	s.StartTime = script.StartTime
	s.StopTime = script.StopTime
	s.Update = script.Update
	s.AI = &alert.AlertInfo{}
	s.AT = script.AT
	s.Liveness = script.Liveness
	s.Always = script.Always
	s.AlwaysSign = script.Always
	s.DeleteWhenExit = script.DeleteWhenExit
	s.PreStart = script.PreStart
	s.DisableAlert = script.DisableAlert

	if script.Cron != nil {
		s.Cron = &cron.Cron{
			Start:   script.Cron.Start,
			Loop:    script.Cron.Loop,
			IsMonth: script.Cron.IsMonth,
			Times:   script.Cron.Times,
		}
	}

	for k, v := range script.TempEnv {
		s.Env[k] = v
	}

	if s.Port > 0 {
		s.Port = pkg.GetAvailablePort(s.Port)
		s.Env["PORT"] = strconv.Itoa(s.Port)
	} else {
		s.Env["PORT"] = "0"
	}

	if s.SubName != "" {
		s.Env["NAME"] = s.SubName
	}
}

func (s *Server) InitLogger() {
	if global.CS.LogDir == "" {
		global.CS.LogDir = "log"
	}
	s.Logger = golog.NewLog(
		filepath.Join(global.CS.LogDir, s.SubName+".log"), 0, true)
}

func (s *Server) PrepareStart() {
	s.InitLogger()

	s.Exit = make(chan int, 2)
	s.CancelProcess = make(chan bool, 2)
	s.Status.Command = s.Command
}

func (s *Server) Close() {
	if s.Logger != nil {
		s.Logger.Close()
	}
	s.Cmd = nil
}

func (s *Server) ResetStatus() {
	s.Status.Status = status.STOP
	s.Status.Pid = 0
	s.Status.CanNotStop = false
	s.Status.RestartCount = 0
	s.Status.Start = 0
	s.Status.Command = ""
}
