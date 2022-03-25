package pkg

import "testing"

func TestCompareSlice(t *testing.T) {
	a := []string{"1111"}
	b := []string{"2222"}
	t.Log(CompareSlice(a, b))
}
