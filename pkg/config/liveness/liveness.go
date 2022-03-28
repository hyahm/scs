package liveness

import (
	"net"
	"net/http"
	"os/exec"
	"runtime"
)

type Liveness struct {
	Http  string `json:"http"`
	Tcp   string `json:"tcp"`
	Shell string `json:"shell"`
}

func newCommand(command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("powershell", "/C", command)
	} else {
		return exec.Command("/bin/bash", "-c", command)
	}
}

func (liveness *Liveness) Ready() bool {
	if liveness.Http != "" {
		resp, err := http.Get(liveness.Http)
		if err != nil {
			return false
		}
		return resp.StatusCode == 200
	}
	if liveness.Tcp != "" {
		conn, err := net.Dial("tcp", liveness.Tcp)
		if err != nil {
			return false
		}
		defer conn.Close()
		return true
	}

	if liveness.Shell != "" {
		cmd := newCommand(liveness.Shell)
		_, err := cmd.Output()
		return err == nil
	}
	return true
}
