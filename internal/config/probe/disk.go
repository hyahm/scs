package probe

import (
	"fmt"
	"sync"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/internal/config/alert"
	"github.com/hyahm/scs/pkg/message"
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
	Percent  float64              // 硬盘使用百分百
	AI       *alert.AlertInfo     // 报警器信息
	Interval time.Duration        // 间隔时间
	Dp       []disk.PartitionStat // 监控的磁盘分区
	AlertDp  map[string]struct{}  // 当前报警的磁盘分区, 单线程监控磁盘，所以不用加锁
	AdLocker *sync.RWMutex
}

func NewDisk() *Disk {
	return &Disk{
		Percent: healthDetector.Config.Disk,
		Dp:      healthDetector.Config.Dp,
		AI: &alert.AlertInfo{
			AM:                 &message.Message{},
			ContinuityInterval: healthDetector.Config.ContinuityInterval,
		},
		Interval: healthDetector.Config.Interval,
		AlertDp:  make(map[string]struct{}),
		AdLocker: &sync.RWMutex{},
	}
}
func (d *Disk) Update() {
	d.Percent = healthDetector.Config.Disk
	d.Interval = healthDetector.Config.Interval
	d.AI.ContinuityInterval = healthDetector.Config.ContinuityInterval
}

var brokenMountPoint string
var brokenStat *disk.UsageStat

func (d *Disk) Check() {
	// 检测硬盘问题
	for _, part := range d.Dp {
		di, err := disk.Usage(part.Mountpoint)
		if err != nil {
			golog.Error(err)
			continue
		}
		currPercent := float64(di.Used) / float64(di.Total) * 100
		if currPercent >= d.Percent {
			d.AI.AM.DiskPath = part.Mountpoint
			brokenStat = di
			d.AdLocker.Lock()
			d.AlertDp[part.Mountpoint] = struct{}{}
			d.AdLocker.Unlock()
			d.AI.AM.Use = di.Used / 1024 / 1024 / 1024
			d.AI.AM.Total = di.Total / 1024 / 1024 / 1024
			d.AI.AM.UsePercent = float64(di.Used / di.Total)
			brokenMountPoint = part.Mountpoint
			d.AI.BreakDown(fmt.Sprintf("硬盘使用率超过 %.2f%%, 当前使用率 %.2f%%", d.Percent, currPercent))
			return
		}
		d.AI.AM.Use = di.Used / 1024 / 1024 / 1024

	}
	d.AdLocker.Lock()
	defer d.AdLocker.Unlock()
	if _, ok := d.AlertDp[brokenMountPoint]; ok {
		d.AI.AM.DiskPath = brokenMountPoint
		d.AlertDp[brokenMountPoint] = struct{}{}
		d.AI.AM.Use = brokenStat.Used / 1024 / 1024 / 1024
		d.AI.AM.Total = brokenStat.Total / 1024 / 1024 / 1024
		d.AI.AM.UsePercent = float64(brokenStat.Used / brokenStat.Total)
		delete(d.AlertDp, brokenMountPoint)
		brokenMountPoint = ""
		brokenStat = nil
		d.AI.Recover("硬盘恢复")
	}
}
