package node

import (
	"fmt"
	"scs/script"
	"strconv"
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
	// 排序并计算最大距离
	for _, v := range st {
		spaceLen[&maxspace1] = len(v.PName)
		spaceLen[&maxspace2] = len(v.Name)
		spaceLen[&maxspace3] = len(v.Status)
		spaceLen[&maxspace4] = len(strconv.Itoa(v.Ppid))
		spaceLen[&maxspace5] = len(v.Start)
		spaceLen[&maxspace6] = len(v.Version)
		spaceLen[&maxspace7] = len("CanNotStop")
		spaceLen[&maxspace7] = len("Failed")
		for s, l := range spaceLen {
			if l+ds > s.Int() {
				*s = space(l + ds)
			}
		}

	}

	for i := 0; i < len(st); i++ {
		min := i
		for j := i + 1; j < len(st); j++ {
			if st[j].Name < st[min].Name {
				// if st[j].PName[0] < st[min].PName[0] {
				// fmt.Println(min)
				min = j
			}
		}
		if min != i {
			st[i], st[min] = st[min], st[i]
		}
		// fmt.Println(st[i])
	}
	fmt.Printf("<node: %s, url: %s>\n", name, url)
	fmt.Println("--------------------------------------------------")
	fmt.Println("PName" + (maxspace1 - space(len("PName"))).String() +
		"Name" + (maxspace2 - space(len("Name"))).String() +
		"Status" + (maxspace3 - space(len("Status"))).String() +
		"Ppid" + (maxspace4 - space(len("Ppid"))).String() +
		"UpTime" + (maxspace5 - space(len("UpTime"))).String() +
		"Verion" + (maxspace6 - space(len("Version"))).String() +
		"CanNotStop" + (maxspace7 - space(len("CanNotStop"))).String() +
		"Failed" + (maxspace8 - space(len("Failed"))).String() +
		"Command")
	for _, v := range st {
		var canNotStopSpace space
		if v.CanNotStop {
			canNotStopSpace = 4
		} else {
			canNotStopSpace = 5
		}

		fmt.Printf("%s%s%s%s%s%s%d%s%s%s%s%s%t%s%d%s%s\n",
			v.PName, maxspace1-space(len(v.PName)),
			v.Name, maxspace2-space(len(v.Name)),
			v.Status, maxspace3-space(len(v.Status)),
			v.Ppid, maxspace4-space(len(strconv.Itoa(v.Ppid))),
			v.Start, maxspace5-space(len(v.Start)),
			v.Version, maxspace6-space(len(v.Version)),
			v.CanNotStop, maxspace7-canNotStopSpace,
			v.RestartCount, maxspace8-space(len(strconv.Itoa(v.RestartCount))),
			"cd "+v.Path+" && "+v.Command,
		)
	}
	fmt.Println("--------------------------------------------------")
}
