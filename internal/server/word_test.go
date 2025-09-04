package server

import (
	"regexp"
	"testing"
)

func TestWord(t *testing.T) {
	reg, err := regexp.Compile(`^\w+$`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reg.MatchString("spidf_asdf"))
}
