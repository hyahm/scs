package subname

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
	end := strings.LastIndex(s.String(), "_")
	if end < 0 {
		return ""
	}
	return s.String()[:end]
}
