package scs

import (
	"os"
	"os/exec"
	"strings"

	"github.com/hyahm/golog"
)

func (svc *Server) LookCommandPath() error {
	for _, v := range svc.Script.LookPath {
		if strings.Trim(v.Path, " ") == "" && strings.Trim(v.Command, " ") == "" {
			continue
		}
		if strings.Trim(v.Path, " ") != "" {
			golog.Info("check path: ", v.Path)
			_, err := os.Stat(v.Path)
			if !os.IsNotExist(err) {
				continue
			}
		}
		if strings.Trim(v.Command, " ") != "" {
			golog.Info("check command: ", v.Command)
			_, err := exec.LookPath(v.Command)
			if err == nil {
				continue
			}
		}
		golog.Info("exec: ", v.Install)
		if err := svc.shell(v.Install); err != nil {
			golog.Error(v.Install)
			return err
		}
	}
	return nil
}
