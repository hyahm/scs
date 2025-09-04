package ssh

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// ssh请求

type Sshd struct {
	Port int    `yaml:"port"`
	Dev  string `yaml:"dev"`
}

func (ssh *Sshd) Dump() error {
	handle, err := pcap.OpenLive(ssh.Dev, 65535, false, pcap.BlockForever)
	if err != nil {
		return err
	}

	if err := handle.SetBPFFilter(fmt.Sprintf("tcp and src port %d", ssh.Port)); err != nil {
		return err
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	seqnum := 0
	for packet := range packetSource.Packets() {
		ssh.decodePacket(packet, seqnum)
		seqnum++
	}
	return nil
}

func (ssh *Sshd) decodePacket(packet gopacket.Packet, seqnum int) {
	// iplayer := packet.Layer(layers.LayerTypeIPv4)
	// [69 0 0 124 56 73 64 0 128 6 219 5 192 168 50 226 192 168 50 250]
	// if iplayer == nil {
	// 	return
	// }
	// if uint8(iplayer.LayerContents()[3]) > 40 {
	// 	return
	// }
	// ip, _ := iplayer.(*layers.IPv4)
	// fmt.Println(Ipv4ToString(ip.SrcIP.To4()))

	tcpLayer := packet.Layer(layers.LayerTypeTCP)

	// src port(16)  dst port(16)  seq(32) ack(32) xxx(64)
	// [239 120 0 22 161 202 90 24 106 151 50 196 80 16 32 19 31 199 0 0]
	// fmt.Println(tcpLayer.LayerContents())
	fmt.Println(seqnum)
	fmt.Println(tcpLayer.LayerPayload())
	// authlayer := packet.Layer(layers.LayerTypeDot11MgmtAuthentication)
	// if authlayer == nil {
	// 	return
	// }
	// auth, _ := authlayer.(*layers.Dot11MgmtAuthentication)
	// fmt.Println(auth.Status.String())
}
