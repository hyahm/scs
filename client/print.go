package client

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hyahm/scs/pkg"
)

type ScriptStatusNode struct {
	Nodes   []pkg.ServiceStatus
	Name    string
	Version string
	Url     string
	Filter  []string
}

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

func boolSpace(b bool) int {
	if b {
		return 4
	}
	return 5
}

func (st *ScriptStatusNode) SortAndPrint(verbose bool) {

	if len(st.Filter) > 0 && len(st.Nodes) == 0 {
		return
	}

	// 按照 name 的顺序进行排序
	for i := 0; i < len(st.Nodes); i++ {
		min := i
		for j := i + 1; j < len(st.Nodes); j++ {
			if st.Nodes[j].Name < st.Nodes[min].Name {
				min = j
			}
		}
		if min != i {
			st.Nodes[i], st.Nodes[min] = st.Nodes[min], st.Nodes[i]
		}

	}
	// 计算当前列最大距离

	maxColumeLen := []int{5, 4, 6, 3, 6, 6, 7, 10, 6, 7, 3, 7}
	for i, v := range st.Nodes {

		if st.Nodes[i].Start > 0 {
			st.Nodes[i].Start = time.Now().Unix() - st.Nodes[i].Start
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
	}
	for i := 0; i < 6; i++ {
		if maxColumeLen[i] == 0 {
			maxColumeLen[i] = 10
			continue
		}
		maxColumeLen[i] += ds
	}
	if maxColumeLen[4] < 10 {
		maxColumeLen[4] = 10
	}
	maxColumeLen[6] = 12
	maxColumeLen[7] = 13
	maxColumeLen[8] = 10
	maxColumeLen[9] = 10
	maxColumeLen[10] = 10
	maxColumeLen[11] = 12

	if verbose {
		fmt.Printf("<node: %s, url: %s, server version: %s>\n", st.Name, st.Url, st.Version)
		fmt.Println("--------------------------------------------------")
		fmt.Println("PName" + (space(maxColumeLen[0]) - space(len("PName"))).String() +
			"Name" + (space(maxColumeLen[1] - len("Name"))).String() +
			"Status" + (space(maxColumeLen[2] - len("Status"))).String() +
			"Pid" + (space(maxColumeLen[3] - len("Pid"))).String() +
			"UpTime" + (space(maxColumeLen[4] - len("UpTime"))).String() +
			"IsCron" + (space(maxColumeLen[5] - len("MEM(kb)"))).String() +
			"Version" + (space(maxColumeLen[6] - len("Version"))).String() +
			"CanNotStop" + (space(maxColumeLen[7] - len("CanNotStop"))).String() +
			"Disable" + (space(maxColumeLen[8] - len("Disable"))).String() +
			"Failed" + (space(maxColumeLen[9] - len("Failed"))).String() +
			"CPU" + (space(maxColumeLen[10] - len("CPU"))).String() +
			"MEM(kb)" + (space(maxColumeLen[11] - len("MEM(kb)"))).String() +
			"Command")
		for _, info := range st.Nodes {

			// var canNotStopSpace int
			// if info.CanNotStop {
			// 	canNotStopSpace = 4
			// } else {
			// 	canNotStopSpace = 5
			// }

			cpu := fmt.Sprintf("%.2f", info.Cpu)
			mem := fmt.Sprintf("%d", info.Mem/1024)
			command := info.Command
			if info.Path != "" {
				if info.OS == "windows" {
					command = fmt.Sprintf("try{ cd %s ; %s } catch{$error[0];break}", info.Path, info.Command)

				} else {
					command = fmt.Sprintf("cd %s && %s", info.Path, info.Command)
				}

			}
			fmt.Printf("%s%s%s%s%s%s%d%s%s%s%t%s%s%s%t%s%t%s%d%s%s%s%s%s%s\n",
				info.PName, space(maxColumeLen[0]-len(info.PName)),
				info.Name, space(maxColumeLen[1]-len(info.Name)),
				info.Status, space(maxColumeLen[2]-len(info.Status)),
				info.Pid, space(maxColumeLen[3]-len(strconv.Itoa(info.Pid))),
				(time.Second * time.Duration(info.Start)).String(), space(maxColumeLen[4]-len((time.Second*time.Duration(info.Start)).String())).String(),
				info.IsCron, space(maxColumeLen[5]-boolSpace(info.IsCron)),
				info.Version, space(maxColumeLen[6]-len(info.Version)),
				info.CanNotStop, space(maxColumeLen[7]-boolSpace(info.CanNotStop)),
				info.Disable, space(maxColumeLen[8]-boolSpace(info.Disable)),
				info.RestartCount, space(maxColumeLen[9]-len(strconv.Itoa(info.RestartCount))),
				cpu, space(maxColumeLen[10]-len(cpu)),
				mem, space(maxColumeLen[11]-len(mem)),
				command,
			)
		}
		fmt.Println("--------------------------------------------------")
	} else {
		fmt.Printf("<node: %s, url: %s, server version: %s>\n", st.Name, st.Url, st.Version)
		fmt.Println("--------------------------------------------------")
		fmt.Println("PName" + (space(maxColumeLen[0]) - space(len("PName"))).String() +
			"Name" + (space(maxColumeLen[1] - len("Name"))).String() +
			"Status" + (space(maxColumeLen[2] - len("Status"))).String() +
			"Pid" + (space(maxColumeLen[3] - len("Pid"))).String() +
			"UpTime" + (space(maxColumeLen[4] - len("UpTime"))).String() +
			"IsCron" + (space(maxColumeLen[5] - len("MEM(kb)"))).String() +
			"Version" + (space(maxColumeLen[6] - len("Version"))).String() +
			"CanNotStop" + (space(maxColumeLen[7] - len("CanNotStop"))).String() +
			"Disable" + (space(maxColumeLen[8] - len("Disable"))).String())
		for _, info := range st.Nodes {
			fmt.Printf("%s%s%s%s%s%s%d%s%s%s%t%s%s%s%t%s%t%s\n",
				info.PName, space(maxColumeLen[0]-len(info.PName)),
				info.Name, space(maxColumeLen[1]-len(info.Name)),
				info.Status, space(maxColumeLen[2]-len(info.Status)),
				info.Pid, space(maxColumeLen[3]-len(strconv.Itoa(info.Pid))),
				(time.Second * time.Duration(info.Start)).String(), space(maxColumeLen[4]-len((time.Second*time.Duration(info.Start)).String())).String(),
				info.IsCron, space(maxColumeLen[5]-boolSpace(info.IsCron)),
				info.Version, space(maxColumeLen[6]-len(info.Version)),
				info.CanNotStop, space(maxColumeLen[7]-boolSpace(info.CanNotStop)),
				info.Disable, space(maxColumeLen[8]-boolSpace(info.Disable)),
			)
		}
		fmt.Println("--------------------------------------------------")
	}

}
