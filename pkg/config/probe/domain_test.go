package probe

import (
	"net/url"
	"testing"
)

func TestDomain(t *testing.T) {
	uri, _ := url.ParseRequestURI("https://127.0.0.1:11111/534")
	t.Log(uri.Host)
}
