package probe

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"scs/alert"
	"strconv"
	"strings"
	"time"

	"github.com/hyahm/golog"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// 保存 probe 报警信息
var GlobalProbe *AlterTimer

//  设置报警时间间隔, 因为是单线程的检测， 所以不用加锁
type AlterTimer struct {
	AT    map[string]*alert.AlertInfo // 保存硬件监控信息
	Probe *Probe

	Exit chan bool
}

// 默认监控的磁盘
var cludeType = map[string]int{
	"EXT4": 0,
	"NTFS": 0,
	"NFS4": 0,
	"XFS":  0,
	"APFS": 0,
}

func init() {
	GlobalProbe = &AlterTimer{
		AT:   make(map[string]*alert.AlertInfo),
		Exit: make(chan bool),
	}
	GlobalProbe.AT["cpu"] = &alert.AlertInfo{}
	GlobalProbe.AT["mem"] = &alert.AlertInfo{}
	GlobalProbe.AT["disk"] = &alert.AlertInfo{}
	GlobalProbe.AT["server"] = &alert.AlertInfo{}
}

func CheckHardWare() {
	GlobalProbe.getDisk()
	for {
		select {
		case <-GlobalProbe.Exit:
			golog.Info("exit check")
			return
		case <-time.After(GlobalProbe.Probe.Interval):
			if GlobalProbe.Probe.Cpu > 0 {
				GlobalProbe.CheckCpu()
			}
			if GlobalProbe.Probe.Mem > 0 {
				GlobalProbe.CheckMem()
			}
			if GlobalProbe.Probe.Disk > 0 {
				GlobalProbe.CheckDisk()
			}
			if len(GlobalProbe.Probe.Monitor) > 0 {
				GlobalProbe.CheckServer()
			}
		}
	}

}

func (at *AlterTimer) CheckServer() {
	for _, server := range at.Probe.Monitor {
		var failed bool
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		//http cookie接口
		cookieJar, _ := cookiejar.New(nil)
		c := &http.Client{
			Jar:       cookieJar,
			Transport: tr,
			Timeout:   time.Second * 5,
		}

		resp, err := c.Get(server + "/probe")
		if err != nil {
			golog.Error(err)
			failed = true
		} else {
			if resp.StatusCode != 200 {
				golog.Error(resp.StatusCode)
				failed = true
			}
		}

		if failed {
			am := &alert.Message{
				Title:    fmt.Sprintf("服务器或scs服务出现问题: %s", server),
				HostName: server,
			}
			if !at.AT["server"].Broken {
				alert.AlertMessage(am, nil)
				at.AT["server"].Broken = true
				at.AT["server"].Start = time.Now()
				at.AT["server"].AlertTime = time.Now()
			} else {
				if time.Since(at.AT["server"].AlertTime) >= at.Probe.ContinuityInterval {
					at.AT["server"].AlertTime = time.Now()
					alert.AlertMessage(am, nil)
				}

			}
			continue
		}
		if at.AT["server"].Broken {
			am := &alert.Message{
				Title:      fmt.Sprintf("服务器或scs服务恢复: %s", server),
				BrokenTime: at.AT["server"].Start.String(),
				FixTime:    time.Now().Local().String(),
			}
			alert.AlertMessage(am, nil)
			at.AT["server"].Broken = false
		}
	}

}

func (at *AlterTimer) getDisk() {
	if at.Probe.dp == nil {
		at.Probe.dp = make([]disk.PartitionStat, 0)
	}
	parts, err := disk.Partitions(true)
	if err != nil {
		golog.Error(err)
		return
	}
	excludePath := make(map[string]int)
	for _, he := range at.Probe.ExcludeDisk {
		excludePath[strings.ToUpper(he)] = 0
	}

	mountNames := make(map[string]string)
	for _, part := range parts {
		if _, ok := excludePath[strings.ToUpper(part.Mountpoint)]; ok {
			continue
		}

		if _, ok := cludeType[strings.ToUpper(part.Fstype)]; ok {
			mountNames[part.Mountpoint] = part.Fstype
			at.Probe.dp = append(at.Probe.dp, part)
			continue
		}

	}
	for _, part := range at.Probe.dp {
		golog.Infof("alert dist: --%s--, type: %s", part.Mountpoint, part.Fstype)
	}

}

func (at *AlterTimer) CheckDisk() {
	for _, part := range at.Probe.dp {
		di, err := disk.Usage(part.Mountpoint)
		if err != nil {
			golog.Error(err)
			continue
		}
		if float64(di.Used)/float64(di.Total)*100 >= at.Probe.Disk {
			am := &alert.Message{
				Title: fmt.Sprintf("硬盘使用率超过%.2f%%", at.Probe.Disk),
			}
			am.DiskPath = part.Mountpoint
			am.Use = di.Used / 1024 / 1024 / 1024
			am.Total = di.Total / 1024 / 1024 / 1024
			am.UsePercent = float64(di.Used / di.Total)
			if !at.AT["disk"].Broken {
				alert.AlertMessage(am, nil)
				at.AT["disk"].Broken = true
				at.AT["disk"].Start = time.Now()
				at.AT["disk"].AlertTime = time.Now()

			} else {
				if time.Since(at.AT["disk"].AlertTime) >= at.Probe.ContinuityInterval {
					at.AT["disk"].AlertTime = time.Now()
					alert.AlertMessage(am, nil)
				}

			}
			continue
		}
		if at.AT["disk"].Broken {
			am := &alert.Message{
				Title:      "硬盘空间已恢复正常",
				BrokenTime: at.AT["disk"].Start.String(),
				FixTime:    time.Now().Local().String(),
			}
			alert.AlertMessage(am, nil)
			at.AT["disk"].Broken = false
		}
	}

}

func (at *AlterTimer) CheckCpu() {
	percents, err := cpu.Percent(time.Second*1, true)
	if err != nil {
		golog.Error(err)
		return
	}
	var totalPercents float64
	for _, percent := range percents {
		totalPercents += percent
	}

	if totalPercents >= at.Probe.Cpu*(float64)(len(percents)) {
		am := &alert.Message{
			Title: fmt.Sprintf("cpu 繁忙超过%.2f%%", at.Probe.Cpu*(float64)(len(percents))),
		}
		am.UsePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", totalPercents), 64)
		am.Top = TopCpu(1)[0].ToString()
		if !at.AT["cpu"].Broken {
			alert.AlertMessage(am, nil)
			at.AT["cpu"].Broken = true
			at.AT["cpu"].Start = time.Now()
			at.AT["cpu"].AlertTime = time.Now()

		} else {
			if time.Since(at.AT["cpu"].AlertTime) >= at.Probe.ContinuityInterval {
				am.BrokenTime = at.AT["cpu"].Start.String()
				at.AT["cpu"].AlertTime = time.Now()
				alert.AlertMessage(am, nil)

			}
		}
		return
	}
	if at.AT["cpu"].Broken {
		am := &alert.Message{
			Title: "cpu 恢复",
		}
		am.BrokenTime = at.AT["cpu"].Start.String()
		am.FixTime = time.Now().Local().String()
		alert.AlertMessage(am, nil)
		at.AT["cpu"].Broken = false
	}
}

func (at *AlterTimer) CheckMem() {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		golog.Error(err)
		return
	}

	if float64(memInfo.Used)/float64(memInfo.Total)*100 >= at.Probe.Disk {
		am := &alert.Message{
			Title: fmt.Sprintf("内存使用率超过 %.2f%%", at.Probe.Mem),
		}
		am.Use = memInfo.Used / 1024 / 1024 / 1024
		am.Total = memInfo.Total / 1024 / 1024 / 1024
		am.Top = TopMem(1)[0].ToString()
		if !at.AT["mem"].Broken {
			// 第一次发送

			alert.AlertMessage(am, nil)
			at.AT["mem"].Broken = true
			at.AT["mem"].Start = time.Now()
			at.AT["mem"].AlertTime = time.Now()
			return
		} else {
			if time.Since(at.AT["mem"].AlertTime) >= at.Probe.ContinuityInterval {
				am.BrokenTime = at.AT["mem"].Start.String()
				at.AT["mem"].AlertTime = time.Now()
				alert.AlertMessage(am, nil)
				return
			}
		}
		return
	}
	if at.AT["mem"].Broken {
		am := &alert.Message{
			Title: "内存恢复正常",
		}
		am.BrokenTime = at.AT["mem"].Start.String()
		am.FixTime = time.Now().Local().String()
		alert.AlertMessage(am, nil)
		at.AT["mem"].Broken = false
	}
}

type inout struct {
	out uint64
	in  uint64
}

func NNetwork() {
	ni := make(map[string]*inout)
	no1, err := net.IOCounters(false)
	if err != nil {
		golog.Error(err)
	}
	for _, n := range no1 {
		fmt.Println(n.Name)
		if _, ok := ni[n.Name]; !ok {
			ni[n.Name] = &inout{}
		}
		ni[n.Name].in = n.BytesRecv
		ni[n.Name].out = n.BytesSent
	}
	time.Sleep(1 * time.Second)
	no2, err := net.IOCounters(false)
	if err != nil {
		golog.Error(err)
	}
	for _, n := range no2 {
		fmt.Println(n.BytesRecv - ni[n.Name].in)
		fmt.Println(n.BytesSent - ni[n.Name].out)
	}
}