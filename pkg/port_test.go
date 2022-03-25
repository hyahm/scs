package pkg

import (
	"testing"

	"github.com/hyahm/golog"
)

func TestPort(t *testing.T) {
	defer golog.Sync()
	t.Log(GetAvailablePort(11111))
}
