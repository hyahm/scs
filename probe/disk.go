package probe

import (
	"fmt"
	"time"

	"github.com/hyahm/scs/alert"

	"github.com/hyahm/golog"
	"github.com/shirou/gopsutil/disk"
)

// 默认监控的磁盘
var cludeType = map[string]int{
	"EXT4": 0,
	"NTFS": 0,
	"NFS4": 0,
	"XFS":  0,
	"APFS": 0,
}

type Disk struct {
	Percent  float64
	AI       *alert.AlertInfo
	Interval time.Duration
	Dp       []disk.PartitionStat
}

func NewDisk(percent float64, dp []disk.PartitionStat, interval, continuityInterval time.Duration) *Disk {
	return &Disk{
		Percent: percent,
		Dp:      dp,
		AI: &alert.AlertInfo{
			AM:                 &alert.Message{},
			ContinuityInterval: continuityInterval,
		},
		Interval: interval,
	}
}
func (d *Disk) Update(probe *Probe) {
	d.Percent = probe.Disk
	d.Interval = probe.Interval
	d.AI.ContinuityInterval = probe.ContinuityInterval
}
func (d *Disk) Check() {
	for _, part := range d.Dp {
		di, err := disk.Usage(part.Mountpoint)
		if err != nil {
			golog.Error(err)
			continue
		}
		if float64(di.Used)/float64(di.Total)*100 >= d.Percent {
			d.AI.AM.DiskPath = part.Mountpoint
			d.AI.AM.Use = di.Used / 1024 / 1024 / 1024
			d.AI.AM.Total = di.Total / 1024 / 1024 / 1024
			d.AI.AM.UsePercent = float64(di.Used / di.Total)
			d.AI.BreakDown(fmt.Sprintf("硬盘使用率超过 %.2f%%", d.Percent))
			return
		}
	}
	d.AI.Recover("硬盘恢复")
}
