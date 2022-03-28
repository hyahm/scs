package probe

import (
	"os"
	"testing"

	"github.com/shirou/gopsutil/disk"
)

func TestDisk(t *testing.T) {
	us, _ := disk.Usage("C:")
	t.Log(us.InodesTotal)
	t.Log(us.InodesUsed)
	ios, err := disk.IOCounters("d:\\")
	if err != nil {
		t.Log(err)
		os.Exit(1)
	}
	for name, io := range ios {
		t.Log(name)
		t.Log(io.IoTime)
	}
}
