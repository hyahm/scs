package node

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hyahm/scs/script"
)

// 打印数据相关

var (
	maxspace1 space = 10 // 最长的第一列与第二列的最短间距  PName
	maxspace2 space = 10 // Name
	maxspace3 space = 10 // Status
	maxspace4 space = 10 //Ppid
	maxspace5 space = 10 //UpTime
	maxspace6 space = 10 // Verion
	maxspace7 space = 15 // CanNotStop
	maxspace8 space = 10 // Failed
)

type status []*script.ServiceStatus

const ds = 4 // 2个列的最近间距
type space int

func (s space) String() string {
	sp := ""
	for i := 0; i < int(s); i++ {
		sp += " "
	}
	return sp
}

func (s space) Int() int {
	return int(s)
}

// fmt.Printf("%s\t%s\t%s\t%s\t%s\n", "PName", "Name", "Status", "Ppid", "CanNotStop")
func (st status) sortAndPrint(name, url string) {
	spaceLen := make(map[*space]int)
	for i := 0; i < len(st); i++ {
		min := i
		for j := i + 1; j < len(st); j++ {
			if st[j].Name < st[min].Name {
				min = j
			}
		}
		if min != i {
			st[i], st[min] = st[min], st[i]
		}
		// fmt.Println(st[i])
	}
	// 排序并计算最大距离
	for _, v := range st {

		if v.Start > 0 {
			v.Start = time.Now().Unix() - v.Start
		}
		spaceLen[&maxspace1] = len(v.PName)
		spaceLen[&maxspace2] = len(v.Name)
		spaceLen[&maxspace3] = len(v.Status)
		spaceLen[&maxspace4] = len(strconv.Itoa(v.Pid))
		spaceLen[&maxspace5] = len((time.Second * time.Duration(v.Start)).String())
		spaceLen[&maxspace6] = len(v.Version)
		spaceLen[&maxspace7] = len("CanNotStop")
		spaceLen[&maxspace7] = len("Failed")
		for s, l := range spaceLen {
			if l+ds > s.Int() {
				*s = space(l + ds)
			}
		}

	}

	fmt.Printf("<node: %s, url: %s>\n", name, url)
	fmt.Println("--------------------------------------------------")
	fmt.Println("PName" + (maxspace1 - space(len("PName"))).String() +
		"Name" + (maxspace2 - space(len("Name"))).String() +
		"Status" + (maxspace3 - space(len("Status"))).String() +
		"Pid" + (maxspace4 - space(len("Pid"))).String() +
		"UpTime" + (maxspace5 - space(len("UpTime"))).String() +
		"Verion" + (maxspace6 - space(len("Version"))).String() +
		"CanNotStop" + (maxspace7 - space(len("CanNotStop"))).String() +
		"Failed" + (maxspace8 - space(len("Failed"))).String() +
		"Command")

	for _, info := range st {

		var canNotStopSpace space
		if info.CanNotStop {
			canNotStopSpace = 4
		} else {
			canNotStopSpace = 5
		}

		fmt.Printf("%s%s%s%s%s%s%d%s%s%s%s%s%t%s%d%s%s\n",
			info.PName, maxspace1-space(len(info.PName)),
			info.Name, maxspace2-space(len(info.Name)),
			info.Status, maxspace3-space(len(info.Status)),
			info.Pid, maxspace4-space(len(strconv.Itoa(info.Pid))),
			(time.Second * time.Duration(info.Start)).String(), maxspace5-space(len((time.Second*time.Duration(info.Start)).String())),
			info.Version, maxspace6-space(len(info.Version)),
			info.CanNotStop, maxspace7-canNotStopSpace,
			info.RestartCount, maxspace8-space(len(strconv.Itoa(info.RestartCount))),
			"cd "+info.Path+" && "+info.Command,
		)
	}
	fmt.Println("--------------------------------------------------")
}
