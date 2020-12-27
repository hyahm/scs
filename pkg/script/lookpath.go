package script

import (
	"os"
	"os/exec"
	"strings"

	"github.com/hyahm/golog"
)

func (s *Script) LookCommandPath() error {
	for _, v := range s.LookPath {
		if strings.Trim(v.Path, " ") != "" {
			golog.Info("check path: ", v.Path)
			_, err := os.Stat(v.Path)
			if !os.IsNotExist(err) {
				continue
			}
		}
		if strings.Trim(v.Command, " ") != "" {
			golog.Info("check command: ", v.Path)
			_, err := exec.LookPath(v.Command)
			if err == nil {
				continue
			}
		}
		golog.Info("exec: ", v.Install)
		if err := Shell(v.Install, s.Env); err != nil {
			golog.Error(v.Install)
			return err
		}
		// check command
		command, err := exec.LookPath(v.Path)
		if err != nil {
			golog.Info(command)

			if err == os.ErrPermission && command != v.Path {
				golog.Error(err)
				golog.Info("exec: ", v.Install)
				if err := Shell(v.Install, s.Env); err != nil {
					golog.Error(v.Install)
					return err
				}
			}
		}
	}
	return nil
}
