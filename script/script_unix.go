// +build !windows

package script

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/hyahm/golog"
)

func Shell(command string, env map[string]string) error {
	cmd := exec.Command("/bin/bash", "-c", command)
	baseEnv := make(map[string]string)
	for _, v := range os.Environ() {
		kv := strings.Split(v, "=")
		baseEnv[kv[0]] = kv[1]
	}
	for k, v := range env {
		if k == "PATH" {
			baseEnv[k] = baseEnv[k] + ";" + v
		} else {
			baseEnv[k] = v
		}
	}
	for k, v := range baseEnv {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	read(cmd)
	err := cmd.Start()
	if err != nil {
		golog.Error(err)
		return err
	}
	defer func() {
		golog.Error(cmd.ProcessState.ExitCode())
	}()
	return cmd.Wait()
}

func (s *Script) Stop() {
	if s.Status.Status == RUNNING {
		s.Status.Status = WAITSTOP
	}
	defer func() {
		if err := recover(); err != nil {
			golog.Info("脚本已经停止了")
		}
	}()
	for {
		time.Sleep(time.Millisecond * 10)
		if !s.Status.CanNotStop {
			s.exit = true
			s.cancel()
			err := syscall.Kill(-s.cmd.Process.Pid, syscall.SIGKILL)
			if err != nil {
				// 正常来说，不会进来的，特殊问题以后再说
				golog.Error(err)
			}

			golog.Infof("stop %s\n", s.SubName)
			s.Status.RestartCount = 0
			// 默认预留1秒代码退出的时间
			time.Sleep(s.KillTime)
			return
		}
	}

}

func (s *Script) Kill() {
	// 数组存日志
	// s.Log = make([]string, Config.LogCount)
	// s.cancel()
	s.exit = true
	var err error

	err = syscall.Kill(-s.cmd.Process.Pid, syscall.SIGKILL)
	// err = s.cmd.Process.Kill()
	// err = exec.Command("kill", "-9", fmt.Sprint(s.cmd.Process.Pid)).Run()
	// err := s.cmd.Process.Kill()

	if err != nil {
		// 正常来说，不会进来的，特殊问题以后再说
		golog.Error(err)
		// return
	}
	s.stopStatus()

	return

}

func (s *Script) start() error {
	s.cmd = exec.Command("/bin/bash", "-c", s.Command)
	s.cmd.Dir = s.Dir
	s.Command = "cd " + s.cmd.Dir + " && " + s.Command
	baseEnv := make(map[string]string)
	for _, v := range os.Environ() {
		kv := strings.Split(v, "=")
		baseEnv[kv[0]] = kv[1]
	}
	for k, v := range s.Env {
		if k == "PATH" {
			baseEnv[k] = baseEnv[k] + ":" + v
		} else {
			baseEnv[k] = v
		}
	}

	for k, v := range baseEnv {
		s.cmd.Env = append(s.cmd.Env, k+"="+v)
	}
	s.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	s.read()
	if err := s.cmd.Start(); err != nil {
		// 执行脚本前的错误, 改变状态
		golog.Error(err)
		return err
	}

	if s.cmd.Process == nil {
		return errors.New("not running")
	}
	return nil
}

func (s *Script) Start() error {
	s.exit = false
	s.Status.Status = RUNNING
	// index := strings.Index(s.Command, " ")
	// s.cmd = exec.Command(s.Command[:index], s.Command[index:])
	if err := s.start(); err != nil {
		return err
	}
	// 等待初始化完成完成后向后执行

	s.Status.Up = time.Now().Unix() // 设置启动状态是成功的
	s.Status.Ppid = s.cmd.Process.Pid
	go s.wait()
	return nil

}

// func (s *Script) Install(command string) {
// 	s.exit = false
// 	golog.Info(s.Log)
// 	// index := strings.Index(s.Command, " ")
// 	// s.cmd = exec.Command(s.Command[:index], s.Command[index:])
// 	s.start(command)
// 	// 等待初始化完成完成后向后执行

// 	go s.waitinstall()
// }
