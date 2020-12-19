package script

import (
	"os/exec"
	"path/filepath"
	"runtime"
)

func (s *Script) LookCommandPath() {
	for _, v := range s.LookPath {
		path, err := exec.LookPath(v.Command)
		if err != nil {
			s.Start(v.Install)
		}

		dir := filepath.Dir(path)
		if runtime.GOOS == "windows" {
			s.Env["PATH"] += ";" + dir
		} else {
			s.Env["PATH"] += ":" + dir
		}

	}
}
