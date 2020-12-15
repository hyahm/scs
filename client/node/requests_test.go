package node

import "testing"

func TestReq(t *testing.T) {
	b, err := Requests("POST", "https://127.0.0.1:11111/test", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}
