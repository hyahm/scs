package command

import (
	"fmt"
	"sync"
	"time"

	"github.com/hyahm/scs/client"
	"github.com/spf13/cobra"
)

func logConfig(cmd *cobra.Command, args []string) {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start).Seconds())
	}()
	if UseNodes != "" {
		if nodeInfo, ok := client.CCfg.GetNode(UseNodes); ok {
			nodeInfo.Log(args[0])

		} else {
			fmt.Println("not found this node")
		}
		return
	}
	if GroupName != "" {
		wg := &sync.WaitGroup{}
		nodes := client.CCfg.GetNodesInGroup(GroupName)
		for _, nodeInfo := range nodes {
			wg.Add(1)
			go func() {
				nodeInfo.Log(args[0])
				wg.Done()
			}()
		}
		wg.Wait()
		return
	}
	wg := &sync.WaitGroup{}

	for _, nodeInfo := range client.CCfg.GetNodes() {
		wg.Add(1)
		go func() {
			nodeInfo.Log(args[0])
			wg.Done()
		}()

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
