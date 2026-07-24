package cache

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/hyahm/scs/pkg"
	"github.com/hyahm/scs/pkg/config"
)

func newCommand(command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("powershell.exe", "-Command", command)
	} else {
		return exec.Command("/bin/sh", "-c", command)
	}
}

func AddScript(s config.Script) {
	// 直接将现在的配置文件保存下来
	SetScript(s.Name, s)
	// if s.ScriptToken == "" {
	// 	s.ScriptToken = pkg.RandomToken()
	// }

	// 将scripts填充到store中
	// so.SetScript(s)
	// 初始化脚本的副本数
	replicate := s.Replicate
	if replicate == 0 {
		replicate = 1
	}
	// 生成这个脚本的的environment
	env := s.MakeTempEnv()
	// 假设设置的端口是可用的
	// 对于每个script 都生成对应的
	// availablePort := s.Port

	for i := 0; i < replicate; i++ {
		subname := fmt.Sprintf("%s_%d", s.Name, i)
		env["SCS_INDEX"] = fmt.Sprintf("%d", i)
		serverStore := &Server{
			Index:   i,
			Name:    s.Name,
			Dir:     s.Dir,
			Command: s.Command,
			Version: s.Version,
			Disable: s.Disable,
			Status: pkg.Status{
				Status: pkg.STOP,
				ChStop: make(chan struct{}),
			},
			Cron:    s.Cron,
			IsCron:  s.Cron.LoopTime > 0,
			Env:     env,
			SubName: subname,
		}
		// env["SCS_INDEX"] = fmt.Sprintf("%d", i)
		// serverStore.Env["SCS_INDEX"] = fmt.Sprintf("%d", i)

		// svc := so.InitServer(i, s.Name, subname)
		// so.SetScriptIndex(s.Name, i)
		// 保存基础数据
		SetServer(subname, serverStore)
		// s.Port = availablePort
		// s.MakeServer()

		// availablePort = s.Port + 1
		if s.Disable {
			// 如果是禁用的 ，那么不用生成多个副本
			// 不在上面就是因为， 是为了看到状态
			return
		}
		// so.SetServer(subname, svc)
		// store.GetStore()
		serverStore.StartAsync()

	}

}
