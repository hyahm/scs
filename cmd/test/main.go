package main

import "github.com/hyahm/scs/pkg/config/packet/ssh"

func main() {
	sshd := &ssh.Sshd{
		Dev:  "enp3s0",
		Port: 22,
	}
	sshd.Dump()
}
