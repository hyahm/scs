package ssh

import (
	"net"
	"os"
	"testing"
)

func TestToip(t *testing.T) {
	a := net.IPv4(12, 51, 84, 68).To4()
	t.Log(Ipv4ToString(a))
}

func TestSshd(t *testing.T) {
	sshd := &Sshd{
		Dev:  "enp3s0",
		Port: 22,
	}
	sshd.Dump()
}

func TestLog(t *testing.T) {
	b, err := os.ReadFile("/var/log/btmp")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(b))
	t.Log(b)
}
