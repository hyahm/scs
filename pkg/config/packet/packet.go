package packet

import "github.com/hyahm/scs/pkg/config/packet/ssh"

type Packet struct {
	Ssh *ssh.Sshd `yaml:"ssh"`
}

type PacketInterface interface {
	Dump()
}
