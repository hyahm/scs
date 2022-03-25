package pkg

import (
	"context"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func GetVersion(command string) string {
	var cmd *exec.Cmd

	if command == "" {
		return ""
	}
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-c", command)
	} else {
		cmd = exec.Command("/bin/bash", "-c", command)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	s := new(string)
	ch := make(chan struct{})
	defer cancel()
	go func(s *string) {
		out, err := cmd.Output()
		if err != nil {
			return
		}
		*s = string(out)
		ch <- struct{}{}
	}(s)

	select {
	case <-ctx.Done():
		return ""
	case <-ch:
		output := strings.ReplaceAll(*s, "\n", "")
		output = strings.ReplaceAll(output, "\r", "")
		return output
	}

}
