package client

import "testing"

func TestClient(t *testing.T) {
	cli := NewClient()
	cli.Domain = "http://127.0.0.1:11111"
	// b, err := cli.StatusAll()
	// if err != nil {
	// 	t.Log(err)
	// }
	// t.Log(string(b))
}
