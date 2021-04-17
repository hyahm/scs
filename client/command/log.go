package command

import (
	"fmt"
	"sync"
	"time"

	"github.com/hyahm/scs/script"
	"github.com/spf13/cobra"
)

func logConfig(cmd *cobra.Command, args []string) {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start).Seconds())
	}()
	if script.UseNodes != "" {
		if nodeInfo, ok := script.CCfg.GetNode(script.UseNodes); ok {
			nodeInfo.Log(args[0])

		} else {
			fmt.Println("not found this node")
		}
		return
	}
	if script.GroupName != "" {
		wg := &sync.WaitGroup{}
		nodes := script.CCfg.GetNodesInGroup(script.GroupName)
		for _, nodeInfo := range nodes {
			wg.Add(1)
			nodeInfo.Wg = wg
			nodeInfo.Log(args[0])
		}
		wg.Wait()
		return
	}
	wg := &sync.WaitGroup{}

	for _, nodeInfo := range script.CCfg.GetNodes() {
		wg.Add(1)
		nodeInfo.Wg = wg
		nodeInfo.Log(args[0])
	}
	wg.Wait()

}

var LogCmd = &cobra.Command{
	Use:   "log",
	Short: "script log",
	Long:  `command: scsctl log [flags] <name>`,
	Args:  cobra.MinimumNArgs(1),
	Run:   logConfig,
}

var tail bool

func init() {
	LogCmd.Flags().BoolVarP(&tail, "tail", "f", false, "tailf ")
	rootCmd.AddCommand(LogCmd)
}
