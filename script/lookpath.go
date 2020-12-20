package script

import (
	"os/exec"
)

func (s *Script) LookCommandPath() {
	s.IsScript = true
	for _, v := range s.LookPath {
		_, err := exec.LookPath(v.Command)
		if err != nil {
			s.Start(v.Install)
		}
	}
	s.IsScript = false
}
