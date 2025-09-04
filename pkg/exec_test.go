package pkg

import (
	"os/exec"
	"testing"
)

func TestExec(t *testing.T) {
	cmd := exec.Command("powershell", "/C", "")
	cmd.Dir = "E:\\test"
}
