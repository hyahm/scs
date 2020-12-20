package script

import (
	"os/exec"

	"github.com/hyahm/golog"
)

func (s *Script) LookCommandPath() error {
	for _, v := range s.LookPath {
		_, err := exec.LookPath(v.Command)
		if err != nil {
			golog.Info(v.Install)
			if err := Shell(v.Install, s.Env); err != nil {
				return err
			}
		}
	}

	return nil
}
