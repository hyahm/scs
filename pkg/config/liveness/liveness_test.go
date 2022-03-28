package liveness

import "testing"

func TestLiveness(t *testing.T) {
	ln := Liveness{
		Shell: "ss",
	}
	t.Log(ln.Ready())
}
