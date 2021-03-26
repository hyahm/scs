package node

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hyahm/scs/script"
)

// 打印数据相关

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

func (st status) sortAndPrint(name, url string) {
	// spaceLen := make(map[*space]int)
	// 安装 name 的顺序进行排序
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
	maxColumeLen := make([]int, 8)
	for _, v := range st {

		if v.Start > 0 {
			v.Start = time.Now().Unix() - v.Start
		}
		if len(v.PName) > maxColumeLen[0] {
			maxColumeLen[0] = len(v.PName)
		}
		if len(v.Name) > maxColumeLen[1] {
			maxColumeLen[1] = len(v.Name)
		}
		if len(v.Status) > maxColumeLen[2] {
			maxColumeLen[2] = len(v.Status)
		}
		if len(strconv.Itoa(v.Pid)) > maxColumeLen[3] {
			maxColumeLen[3] = len(strconv.Itoa(v.Pid))
		}
		if len((time.Second * time.Duration(v.Start)).String()) > maxColumeLen[4] {
			maxColumeLen[4] = len((time.Second * time.Duration(v.Start)).String())
		}
		if len(v.Version) > maxColumeLen[5] {
			maxColumeLen[5] = len(v.Version)
		}

		// spaceLen[&maxspace1] = len(v.PName)
		// maxColumeLen[0] = len(v.PName)
		// spaceLen[&maxspace2] = len(v.Name)
		// spaceLen[&maxspace3] = len(v.Status)
		// spaceLen[&maxspace4] = len(strconv.Itoa(v.Pid))
		// spaceLen[&maxspace5] = len((time.Second * time.Duration(v.Start)).String())
		// spaceLen[&maxspace6] = len(v.Version)
		// spaceLen[&maxspace7] = len("CanNotStop")
		// spaceLen[&maxspace8] = len(strconv.Itoa(info.RestartCount)))

		// for s, l := range spaceLen {
		// 	if l+ds > s.Int() {
		// 		*s = space(l + ds)
		// 	}
		// }
		// fmt.Println(spaceLen[&maxspace6])
	}
	for i := 0; i < 6; i++ {
		maxColumeLen[i] += ds
	}
	maxColumeLen[6] = 15
	maxColumeLen[7] = 15
	fmt.Printf("<node: %s, url: %s>\n", name, url)
	fmt.Println("--------------------------------------------------")
	fmt.Println("PName" + (space(maxColumeLen[0]) - space(len("PName"))).String() +
		"Name" + (space(maxColumeLen[1] - len("Name"))).String() +
		"Status" + (space(maxColumeLen[2] - len("Status"))).String() +
		"Pid" + (space(maxColumeLen[3] - len("Pid"))).String() +
		"UpTime" + (space(maxColumeLen[4] - len("UpTime"))).String() +
		"Verion" + (space(maxColumeLen[5] - len("Version"))).String() +
		"CanNotStop" + (space(maxColumeLen[6] - len("CanNotStop"))).String() +
		"Failed" + (space(maxColumeLen[7] - len("Failed"))).String() +
		"Command")
	for _, info := range st {

		var canNotStopSpace int
		if info.CanNotStop {
			canNotStopSpace = 4
		} else {
			canNotStopSpace = 5
		}
		cdpath := ""
		if info.Path != "" {
			cdpath = "cd " + info.Path + " && "
		}
		fmt.Printf("%s%s%s%s%s%s%d%s%s%s%s%s%t%s%d%s%s\n",
			info.PName, space(maxColumeLen[0]-len(info.PName)),
			info.Name, space(maxColumeLen[1]-len(info.Name)),
			info.Status, space(maxColumeLen[2]-len(info.Status)),
			info.Pid, space(maxColumeLen[3]-len(strconv.Itoa(info.Pid))),
			(time.Second * time.Duration(info.Start)).String(), space(maxColumeLen[4]-len((time.Second*time.Duration(info.Start)).String())).String(),
			info.Version, space(maxColumeLen[5]-len(info.Version)),
			info.CanNotStop, space(maxColumeLen[6]-canNotStopSpace),
			info.RestartCount, space(maxColumeLen[7]-len(strconv.Itoa(info.RestartCount))),
			cdpath+info.Command,
		)
	}
	fmt.Println("--------------------------------------------------")
}
