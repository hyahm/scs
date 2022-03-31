package controller

import "testing"

func TestStroe(t *testing.T) {
	if _, ok := store.serverIndex["aaa"]; !ok {
		store.serverIndex["aaa"] = make(map[int]struct{})
	}
	for i := 0; i < 1000000; i++ {

	}
}
