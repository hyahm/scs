package scs

import (
	"fmt"
	"strings"
)

type Subname string

func NewSubname(s string, i int) Subname {
	return Subname(fmt.Sprintf("%s_%d", s, i))
}

func (s Subname) String() string {
	return string(s)
}

func (s Subname) GetName() string {
	sn := strings.Split(s.String(), "_")
	if len(sn) < 2 {
		return ""
	}
	return sn[0]
}
