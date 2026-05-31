package config

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/hyahm/golog"
	"github.com/shirou/gopsutil/disk"
)

var healthDetector *Detector

type Detector struct {
	Probe  *Probe
	Ctx    context.Context
	Cancel context.CancelFunc
	Config *Config
	Cps    []CheckPointer
}

// type Config struct {
// 	// 通过静态配置生成动态配置
// 	Monitor   []string
// 	Monitored []string
// 	Mem       float64 `yaml:"mem"`
// 	// cpu使用率, 默认90
// 	Cpu float64 `yaml:"cpu"`
// 	IO  float64 `yaml:"io"`
// 	// 硬盘使用率， 默认90
// 	Disk float64              `yaml:"disk"`
// 	Dp   []disk.PartitionStat `yaml:"excludeDisk"`
// 	// 检测间隔， 默认10秒
// 	Interval time.Duration `yaml:"interval"`
// 	// 下次报警时间间隔， 如果恢复了就重置
// 	ContinuityInterval time.Duration `yaml:"continuityInterval"`
// }

// 保存配置文件信息
type Probe struct {
	Monitor   []string `yaml:"monitor,omitempty"`
	Monitored []string `yaml:"monitored,omitempty"`
	// 内存使用率 默认90
	Mem float64 `yaml:"mem,omitempty"`
	IO  float64 `yaml:"io,omitempty"`
	// cpu使用率, 默认90
	Cpu float64 `yaml:"cpu,omitempty"`
	// 硬盘使用率， 默认90
	Disk          float64  `yaml:"disk,omitempty"`
	ExcludeDisk   []string `yaml:"excludeDisk,omitempty"`
	DiskPartition []string `yaml:"diskPartition,omitempty"`
	// 检测间隔， 默认10秒
	Interval time.Duration `yaml:"interval,omitempty"`
	// 下次报警时间间隔， 如果恢复了就重置
	ContinuityInterval time.Duration `yaml:"continuityInterval,omitempty"`
}

func (p *Probe) initProbe() {
	if p == nil {
		p = &Probe{}
	}
	if p.Interval == 0 {
		p.Interval = time.Second * 10
	}
	if len(p.DiskPartition) == 0 {
		p.DiskPartition = getDisk()
	}
	if p.Mem == 0 {
		p.Mem = 90
	}
	if p.Cpu == 0 {
		p.Cpu = 90
	}
	if p.Disk == 0 {
		p.Disk = 90
	}
}

// func initConfig() {
// 	healthDetector.Config.Probe.Interval = healthDetector.Probe.Interval
// 	healthDetector.Config.Probe.ContinuityInterval = healthDetector.Probe.ContinuityInterval
// 	if healthDetector.Config.Probe.Interval == 0 {
// 		healthDetector.Config.Probe.Interval = time.Second * 10
// 	}
// 	if healthDetector.Config.Probe.ContinuityInterval == 0 {
// 		healthDetector.Config.Probe.ContinuityInterval = time.Hour * 1
// 	}
// 	// todo: healthDetector.Config.Dp = getDisk()
// 	healthDetector.Config.Probe.Monitored = healthDetector.Probe.Monitored
// 	healthDetector.Config.Probe.Cpu = healthDetector.Probe.Cpu
// 	healthDetector.Config.Probe.Mem = healthDetector.Probe.Mem
// 	healthDetector.Config.Probe.Disk = healthDetector.Probe.Disk
// 	healthDetector.Config.Probe.Monitor = healthDetector.Probe.Monitor
// 	if healthDetector.Config.Probe.Cpu == 0 {
// 		healthDetector.Config.Probe.Cpu = 90
// 	}
// 	if healthDetector.Config.Probe.Mem == 0 {
// 		healthDetector.Config.Probe.Mem = 90
// 	}
// 	if healthDetector.Config.Probe.Disk == 0 {
// 		healthDetector.Config.Probe.Disk = 85
// 	}

// 	if healthDetector.Config.Probe.Cpu > 0 ||
// 		healthDetector.Config.Probe.Mem > 0 ||
// 		healthDetector.Config.Probe.Disk > 0 ||
// 		healthDetector.Config.Probe.IO > 0 ||
// 		len(healthDetector.Config.Probe.Monitor) > 0 {
// 		go CheckHardWare()
// 	}

// }

// func (p *Probe) CheckHardWare() {

// 	if p.Cpu > 0 {
// 		if IsNil(healthDetector.Cps[0]) {
// 			healthDetector.Cps[0] = NewCpu()
// 		} else {
// 			healthDetector.Cps[0].Update()
// 		}
// 	} else {
// 		healthDetector.Cps[0] = nil
// 	}
// 	if healthDetector.Config.Probe.Mem > 0 {
// 		if IsNil(healthDetector.Cps[1]) {
// 			healthDetector.Cps[1] = NewMem()
// 		} else {
// 			healthDetector.Cps[1].Update()
// 		}

// 	} else {
// 		healthDetector.Cps[1] = nil
// 	}
// 	if healthDetector.Config.Probe.Disk > 0 {
// 		if IsNil(healthDetector.Cps[2]) {
// 			healthDetector.Cps[2] = NewDisk()
// 		} else {
// 			healthDetector.Cps[2].Update()
// 		}
// 	} else {
// 		healthDetector.Cps[2] = nil
// 	}
// 	if len(healthDetector.Config.Probe.Monitor) > 0 {
// 		if IsNil(healthDetector.Cps[3]) {
// 			healthDetector.Cps[3] = NewMonitor()
// 		} else {
// 			healthDetector.Cps[3].Update()
// 		}
// 	} else {
// 		healthDetector.Cps[3] = nil
// 	}

// 	// if healthDetector.Config.IO > 0 {
// 	// 	if IsNil(healthDetector.Cps[4]) {
// 	// 		healthDetector.Cps[4] = NewCpu()
// 	// 	} else {
// 	// 		healthDetector.Cps[4].Update()
// 	// 	}
// 	// } else {
// 	// 	healthDetector.Cps[4] = nil
// 	// }
// 	for {
// 		select {
// 		case <-healthDetector.Ctx.Done():
// 			golog.Info("exit check")
// 			return
// 		case <-time.After(healthDetector.Config.Probe.Interval):
// 			for _, check := range healthDetector.Cps {

// 				if IsNil(check) {
// 					continue
// 				}
// 				check.Check()
// 			}
// 		}
// 	}

// }

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	return !vi.IsValid() || vi.IsNil()
}

func getDisk() []string {
	dp := make([]disk.PartitionStat, 0)
	parts, err := disk.Partitions(true)
	if err != nil {
		golog.Error(err)
		return []string{}
	}
	excludePath := make(map[string]int)
	for _, he := range Cfg.Probe.ExcludeDisk {
		excludePath[strings.ToUpper(he)] = 0
	}

	mountNames := make(map[string]string)
	for _, part := range parts {
		if _, ok := excludePath[strings.ToUpper(part.Mountpoint)]; ok {
			continue
		}

		if _, ok := cludeType[strings.ToUpper(part.Fstype)]; ok {
			mountNames[part.Mountpoint] = part.Fstype
			dp = append(dp, part)
			continue
		}

	}
	list := make([]string, 0, len(dp))
	for _, part := range dp {
		list = append(list, part.Mountpoint)
		golog.Infof("alert dist: --%s--, type: %s", part.Mountpoint, part.Fstype)
	}
	return list
}
