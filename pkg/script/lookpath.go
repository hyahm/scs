package script

import (
	"os/exec"

	"github.com/hyahm/golog"
)

func (s *Script) LookCommandPath() error {
	for _, v := range s.LookPath {
		_, err := exec.LookPath(v.Path)
		if err != nil {
			if err := Shell(v.Install, s.Env); err != nil {
				golog.Error(v.Install)
				return err
			}
		}
	}
	return nil
}
