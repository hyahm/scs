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
	maxColumeLen[7] = 12
	maxColumeLen[8] = 9
	maxColumeLen[9] = 8
	maxColumeLen[10] = 10
	maxColumeLen[11] = 10
	maxlength := make([]any, len(maxColumeLen))
	for k, v := range maxColumeLen {
		maxlength[k] = v
	}
	headerName := []any{"PName", "Name", "Status", "Pid", "UpTime", "IsCron", "Version", "CanNotStop", "Disable", "Failed", "CPU", "MEM(kb)", "Command"}
	if verbose {
		fmt.Printf("<node: %s, url: %s, server version: %s>\n", st.Name, st.Url, st.Version)
		fmt.Println("--------------------------------------------------")
		header := fmt.Sprintf("%%-%ds%%-%ds%%-%ds%%-%ds%%-%ds%%-%ds%%-%ds%%-%ds%%-%ds%%-%ds%%-%ds%%-%ds%%s\n", maxlength...)
		fmt.Printf(header, headerName...)
		row := fmt.Sprintf("%%-%ds%%-%ds%%-%ds%%-%dd%%-%ds%%-%dt%%-%ds%%-%dt%%-%dt%%-%dd%%-%ds%%-%ds%%s\n", maxlength...)
		for _, info := range st.Nodes {
			command := info.Command
			if info.Path != "" {
				if info.OS == "windows" {
					command = fmt.Sprintf("try{ cd %s ; %s } catch{$error[0];break}", info.Path, info.Command)

				} else {
					command = fmt.Sprintf("cd %s && %s", info.Path, info.Command)
				}

			}
			fmt.Printf(row, info.PName, info.Name, info.Status, info.Pid,
				(time.Second * time.Duration(info.Start)).String(), info.IsCron,
				info.Version, info.CanNotStop, info.Disable, info.RestartCount,
				fmt.Sprintf("%.2f", info.Cpu), fmt.Sprintf("%d", info.Mem/1024), command,
			)

		}
		fmt.Println("--------------------------------------------------")
	} else {
		fmt.Printf("<node: %s, url: %s, server version: %s>\n", st.Name, st.Url, st.Version)
		fmt.Println("--------------------------------------------------")

		header := fmt.Sprintf("%%-%ds%%-%ds%%-%ds%%-%ds%%-%ds%%-%ds%%-%ds%%-%ds%%-%ds\n", maxlength[0:9]...)
		fmt.Printf(header, headerName[0:9]...)
		row := fmt.Sprintf("%%-%ds%%-%ds%%-%ds%%-%dd%%-%ds%%-%dt%%-%ds%%-%dt%%-%dt\n", maxlength[0:9]...)
		for _, info := range st.Nodes {
			fmt.Printf(row, info.PName, info.Name, info.Status, info.Pid,
				(time.Second * time.Duration(info.Start)).String(), info.IsCron, info.Version, info.CanNotStop, info.Disable)
		}
		fmt.Println("--------------------------------------------------")
	}

}
