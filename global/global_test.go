package global

import "testing"

type m struct {
	A []string
	B []string
}

func TestGlobal(t *testing.T) {
	var a []string
	t.Log(len(a))
	x := &m{
		B: nil,
	}
	t.Log(len(x.A))
	t.Log(len(x.B))
}
