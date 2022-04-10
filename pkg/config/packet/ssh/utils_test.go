package ssh

import (
	"net"
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
