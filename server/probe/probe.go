package probe

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/scs/global"
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

type Config struct {
	// 通过静态配置生成动态配置
	Monitor   []string
	Monitored []string
	Mem       float64 `yaml:"mem"`
	// cpu使用率, 默认90
	Cpu float64 `yaml:"cpu"`
	// 硬盘使用率， 默认90
	Disk float64              `yaml:"disk"`
	Dp   []disk.PartitionStat `yaml:"excludeDisk"`
	// 检测间隔， 默认10秒
	Interval time.Duration `yaml:"interval"`
	// 下次报警时间间隔， 如果恢复了就重置
	ContinuityInterval time.Duration `yaml:"continuityInterval"`
}

//  保存配置文件信息
type Probe struct {
	Monitor   []string `yaml:"monitor"`
	Monitored []string `yaml:"monitored"`
	// 内存使用率 默认90
	Mem float64 `yaml:"mem"`
	// cpu使用率, 默认90
	Cpu float64 `yaml:"cpu"`
	// 硬盘使用率， 默认90
	Disk        float64  `yaml:"disk"`
	ExcludeDisk []string `yaml:"excludeDisk"`
	// 检测间隔， 默认10秒
	Interval time.Duration `yaml:"interval"`
	// 下次报警时间间隔， 如果恢复了就重置
	ContinuityInterval time.Duration `yaml:"continuityInterval"`
}

func Run(p *Probe) {
	if healthDetector == nil {
		healthDetector = &Detector{
			Probe: p,
			Config: &Config{
				Monitored: make([]string, 0),
			},
			Cps: make([]CheckPointer, 4),
		}
	} else {
		//healthDetector.Cancel()
	}

	healthDetector.Ctx, healthDetector.Cancel = context.WithCancel(context.Background())

	initConfig()
}

func initConfig() {
	if healthDetector.Probe.Interval == 0 {
		healthDetector.Config.Interval = time.Second * 10
	}
	if healthDetector.Probe.ContinuityInterval == 0 {
		healthDetector.Config.ContinuityInterval = time.Hour * 1
	}
	healthDetector.Config.Dp = getDisk()
	healthDetector.Config.Monitored = healthDetector.Probe.Monitored

	if healthDetector.Probe.Cpu == 0 {
		healthDetector.Config.Cpu = 90
	}
	if healthDetector.Probe.Mem == 0 {
		healthDetector.Config.Mem = 90
	}
	if healthDetector.Probe.Disk == 0 {
		healthDetector.Config.Disk = 85
	}
	healthDetector.Config.Monitor = healthDetector.Probe.Monitor
	global.Monitored = healthDetector.Probe.Monitored

	if healthDetector.Config.Cpu > 0 ||
		healthDetector.Config.Mem > 0 ||
		healthDetector.Config.Disk > 0 ||
		len(healthDetector.Config.Monitor) > 0 {
		go CheckHardWare()
	}

}

func CheckHardWare() {
	if healthDetector.Config.Cpu > 0 {
		if IsNil(healthDetector.Cps[0]) {
			healthDetector.Cps[0] = NewCpu()
		} else {
			healthDetector.Cps[0].Update()
		}
	} else {
		healthDetector.Cps[0] = nil
	}
	if healthDetector.Config.Mem > 0 {
		if IsNil(healthDetector.Cps[1]) {
			healthDetector.Cps[1] = NewMem()
		} else {
			healthDetector.Cps[1].Update()
		}

	} else {
		healthDetector.Cps[1] = nil
	}
	if healthDetector.Config.Disk > 0 {
		if IsNil(healthDetector.Cps[2]) {
			healthDetector.Cps[2] = NewDisk()
		} else {
			healthDetector.Cps[2].Update()
		}
	} else {
		healthDetector.Cps[2] = nil
	}
	if len(healthDetector.Config.Monitor) > 0 {
		if IsNil(healthDetector.Cps[3]) {
			golog.Info("new")
			healthDetector.Cps[3] = NewMonitor()
		} else {
			healthDetector.Cps[3].Update()
		}
	} else {
		healthDetector.Cps[3] = nil
	}
	for {
		select {
		case <-healthDetector.Ctx.Done():
			golog.Info("exit check")
			return
		case <-time.After(healthDetector.Config.Interval):
			for _, check := range healthDetector.Cps {
				if IsNil(check) {
					continue
				}
				check.Check()
			}
		}
	}

}

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	return !vi.IsValid() || vi.IsNil()
}

func getDisk() []disk.PartitionStat {
	dp := make([]disk.PartitionStat, 0)
	parts, err := disk.Partitions(true)
	if err != nil {
		golog.Error(err)
		return dp
	}
	excludePath := make(map[string]int)
	for _, he := range healthDetector.Probe.ExcludeDisk {
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
	for _, part := range dp {
		golog.Infof("alert dist: --%s--, type: %s", part.Mountpoint, part.Fstype)
	}
	return dp
}
