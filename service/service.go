package service

import (
	"errors"

	"github.com/hyahm/scs/internal/server"
	sstore "github.com/hyahm/scs/internal/store"
	"github.com/hyahm/scs/pkg/config/alert"
	"github.com/hyahm/scs/pkg/config/scripts"
	"github.com/hyahm/scs/pkg/message"
)

var (
	ErrServerNotFound = errors.New("server not found")
	ErrScriptNotFound = errors.New("script not found")
)

type Service interface {
	// Script 操作
	AddScript(script *scripts.Script) error
	GetScript(name string) (*scripts.Script, bool)
	DeleteScript(name string)
	UpdateScript(script *scripts.Script) error
	GetAllScripts() map[string]*scripts.Script

	// Server 操作
	StartServer(name string) error
	StopServer(name string) error
	RestartServer(name string) error
	RemoveServer(name string) error
	KillServer(name string) error
	GetServer(name string) (*server.Server, bool)
	GetAllServers() map[string]*server.Server

	// 批量操作
	StartAll()
	StopAll()
	RestartAll()
	StopByScriptNames(names map[string]struct{})
	RestartByScriptNames(names map[string]struct{})
	GetServersByScriptNames(names map[string]struct{}) map[string]*server.Server

	// 状态 - 使用 server.Server 中的 Status 字段
	GetStatus(name string) (*server.Server, error)
	GetAllStatus() map[string]*server.Server

	// 告警
	GetAlerts() map[string]message.SendAlerter

	// 认证
	GetAuthByToken(token string) []AuthInfo
}

type AuthInfo struct {
	ServerName string
	ScriptName string
	Role       string
}

type service struct {
	store *sstore.Store
}

func New() Service {
	return &service{store: sstore.GetStore()}
}

func (s *service) AddScript(script *scripts.Script) error {
	s.store.SetScript(script)
	return nil
}

func (s *service) GetScript(name string) (*scripts.Script, bool) {
	return s.store.GetScriptByName(name)
}

func (s *service) DeleteScript(name string) {
	s.store.DeleteScriptByName(name)
}

func (s *service) UpdateScript(script *scripts.Script) error {
	s.store.SetScript(script)
	return nil
}

func (s *service) GetAllScripts() map[string]*scripts.Script {
	return s.store.GetAllScriptMap()
}

func (s *service) StartServer(name string) error {
	svc, ok := s.store.GetServerByName(name)
	if !ok {
		return ErrServerNotFound
	}
	svc.Start()
	return nil
}

func (s *service) StopServer(name string) error {
	svc, ok := s.store.GetServerByName(name)
	if !ok {
		return ErrServerNotFound
	}
	svc.Stop()
	return nil
}

func (s *service) RestartServer(name string) error {
	svc, ok := s.store.GetServerByName(name)
	if !ok {
		return ErrServerNotFound
	}
	svc.Restart()
	<-svc.StopSignal
	script, ok := s.store.GetScriptByName(svc.Name)
	if !ok {
		return ErrScriptNotFound
	}
	svc.MakeServer(script)
	svc.Start()
	return nil
}

func (s *service) RemoveServer(name string) error {
	svc, ok := s.store.GetServerByName(name)
	if !ok {
		return ErrServerNotFound
	}
	svc.Remove()
	<-svc.StopSignal
	s.store.DeleteServerByName(name)
	s.store.DeleteScriptIndex(svc.Name, svc.Index)
	if s.store.GetScriptLength(svc.Name) == 0 {
		s.store.DeleteScriptByName(svc.Name)
	}
	return nil
}

func (s *service) KillServer(name string) error {
	svc, ok := s.store.GetServerByName(name)
	if !ok {
		return ErrServerNotFound
	}
	svc.Kill()
	return nil
}

func (s *service) GetServer(name string) (*server.Server, bool) {
	return s.store.GetServerByName(name)
}

func (s *service) GetAllServers() map[string]*server.Server {
	return s.store.GetAllServerMap()
}

func (s *service) StartAll() {
	for _, svc := range s.store.GetAllServer() {
		svc.Start()
	}
}

func (s *service) StopAll() {
	for _, script := range s.store.GetAllScriptMap() {
		s.stopScript(script)
	}
}

func (s *service) RestartAll() {
	for _, svc := range s.store.GetAllServer() {
		if svc.Disable {
			continue
		}
		s.RestartServer(svc.SubName)
	}
}

func (s *service) StopByScriptNames(names map[string]struct{}) {
	for _, script := range s.store.GetScriptMapFilterByName(names) {
		s.stopScript(script)
	}
}

func (s *service) stopScript(script *scripts.Script) {
	if script.Disable {
		return
	}
	replicate := script.Replicate
	if replicate == 0 {
		replicate = 1
	}
	for i := 0; i < replicate; i++ {
		subname := script.Name + "_" + string(rune(i+'0'))
		if svc, ok := s.store.GetServerByName(subname); ok {
			svc.Stop()
		}
	}
}

func (s *service) RestartByScriptNames(names map[string]struct{}) {
	for pname := range names {
		for _, index := range s.store.GetScriptIndex(pname) {
			subname := pname + "_" + string(rune(index+'0'))
			if svc, ok := s.store.GetServerByName(subname); ok {
				svc.Restart()
			}
		}
	}
}

func (s *service) GetServersByScriptNames(names map[string]struct{}) map[string]*server.Server {
	return s.store.GetServerMapFilterScripts(names)
}

func (s *service) GetStatus(name string) (*server.Server, error) {
	svc, ok := s.store.GetServerByName(name)
	if !ok {
		return nil, ErrServerNotFound
	}
	return svc, nil
}

func (s *service) GetAllStatus() map[string]*server.Server {
	result := make(map[string]*server.Server)
	for name, svc := range s.store.GetAllServerMap() {
		result[name] = svc
	}
	return result
}

func (s *service) GetAlerts() map[string]message.SendAlerter {
	return alert.GetAlerts()
}

func (s *service) GetAuthByToken(token string) []AuthInfo {
	auths := make([]AuthInfo, 0)
	for name, srv := range s.store.GetAllServerMap() {
		if srv.ScriptToken == token {
			auths = append(auths, AuthInfo{
				ServerName: name,
				ScriptName: srv.Name,
				Role:       string(scripts.ScriptRole),
			})
		}
		if srv.SimpleToken == token {
			auths = append(auths, AuthInfo{
				ServerName: name,
				ScriptName: srv.Name,
				Role:       string(scripts.SimpleRole),
			})
		}
	}
	return auths
}
