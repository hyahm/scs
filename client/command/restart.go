package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client/cliconfig"
	"github.com/hyahm/scs/client/node"

	"github.com/spf13/cobra"
)

var restartAll bool
var RestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart assign script",
	Long:  `command: scsctl restart ([flags]) || ([pname] [name])`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && !restartAll {
			fmt.Println("Specify at least one parameter, or -- all")
			return
		}
		if restartAll {
			args = nil
		}
		if node.UseNodes != "" {
			if nodeInfo, ok := cliconfig.Cfg.GetNode(node.UseNodes); ok {
				nodeInfo.Restart(args...)
				return
			}
		}
		if node.GroupName != "" {
			wg := &sync.WaitGroup{}
			nodes := cliconfig.Cfg.GetNodesInGroup(node.GroupName)
			for _, nodeInfo := range nodes {
				wg.Add(1)
				nodeInfo.Wg = wg
				nodeInfo.Restart(args...)
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}

		for _, nodeInfo := range cliconfig.Cfg.GetNodes() {
			wg.Add(1)
			nodeInfo.Wg = wg
			nodeInfo.Restart(args...)
		}
		wg.Wait()

	},
}

func init() {
	RestartCmd.Flags().BoolVar(&restartAll, "all", false, "restart all")
	rootCmd.AddCommand(RestartCmd)
}
