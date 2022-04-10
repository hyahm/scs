package ssh

import (
	"fmt"
	"net"
	"strings"
)

func Ipv4ToString(ip net.IP) string {
	if len(ip) > 4 {
		return ""
	}
	b := make([]string, 4)
	for i, v := range ip {
		b[i] = fmt.Sprintf("%v", v)

	}
	return strings.Join(b, ".")
}
