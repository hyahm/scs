package command

import (
	"fmt"
	"sync"
	"time"

	"github.com/hyahm/scs/client/cliconfig"
	"github.com/hyahm/scs/client/node"

	"github.com/spf13/cobra"
)

func logConfig(cmd *cobra.Command, args []string) {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start).Seconds())
	}()
	if node.UseNodes != "" {
		if nodeInfo, ok := cliconfig.Cfg.Nodes[node.UseNodes]; ok {
			nodeInfo.Log(args[0])
			return
		}
	}
	if node.GroupName != "" {
		wg := &sync.WaitGroup{}
		for _, v := range cliconfig.Cfg.Group[node.GroupName] {
			if nodeInfo, ok := cliconfig.Cfg.Nodes[v]; ok {
				wg.Add(1)
				nodeInfo.Wg = wg
				nodeInfo.Log(args[0])
			}
		}
		wg.Wait()
		return
	}
	wg := &sync.WaitGroup{}
	for name, nodeInfo := range cliconfig.Cfg.Nodes {
		nodeInfo.Name = name
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
