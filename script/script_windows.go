// +build windows
package script

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hyahm/golog"
)

func Shell(command string) ([]byte, error) {
	cmd := exec.Command("cmd", "/C", command)
	return cmd.Output()
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
			err := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(s.cmd.Process.Pid)).Run()
			if err != nil {
				// 正常来说，不会进来的，特殊问题以后再说
				golog.Error(err)
			}

			golog.Infof("stop %s\n", s.SubName)
			s.Status.RestartCount = 0
			// 预留3秒代码退出的时间
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

	err := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(s.cmd.Process.Pid)).Run()
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

func (s *Script) start(command string) error {
	s.cmd = exec.Command("cmd", "/C", command)
	baseEnv := make(map[string]string)
	for _, v := range os.Environ() {
		kv := strings.Split(v, "=")
		baseEnv[kv[0]] = kv[1]
	}
	for k, v := range s.Env {
		if k == "PATH" {
			baseEnv[k] = baseEnv[k] + ";" + v
		} else {
			baseEnv[k] = v
		}
	}
	s.cmd.Dir = s.Dir
	for k, v := range baseEnv {
		s.cmd.Env = append(s.cmd.Env, k+"="+v)
	}

	// 等待初始化完成完成后向后执行
	s.read()
	s.Status.Up = time.Now().Unix() // 设置启动状态是成功的
	if err := s.cmd.Start(); err != nil {
		// 执行脚本前的错误, 改变状态
		golog.Error(err)
		s.stopStatus()
		return err
	}

	if s.cmd.Process == nil {
		s.stopStatus()
		return errors.New("not start")
	}
	return nil
}

func (s *Script) Start(command string) error {
	s.exit = false
	s.Status.Status = RUNNING
	// index := strings.Index(s.Command, " ")
	// s.cmd = exec.Command(s.Command[:index], s.Command[index:])
	if err := s.start(command); err != nil {
		return err
	}
	s.Status.Ppid = s.cmd.Process.Pid
	go s.wait()
	return nil
}
