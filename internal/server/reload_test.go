package server

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	for {
		select {
		case <-time.After(10 * time.Second):
			t.Log(111)
			break
		}
	}
}
