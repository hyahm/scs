package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/cliconfig"
	"github.com/hyahm/scs/server/node"

	"github.com/spf13/cobra"
)

var stopAll bool

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop script",
	Long:  `command: scsctl stop ([flags]) || ([pname] [name])`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && !stopAll {
			fmt.Println("Specify at least one parameter, or -- all")
			return
		}
		if stopAll {
			args = nil
		}
		if node.UseNodes != "" {
			if nodeInfo, ok := cliconfig.Cfg.GetNode(node.UseNodes); ok {
				nodeInfo.Stop(args...)

			} else {
				fmt.Println("not found this node")
			}
			return
		}
		if node.GroupName != "" {
			wg := &sync.WaitGroup{}
			nodes := cliconfig.Cfg.GetNodesInGroup(node.GroupName)
			for _, nodeInfo := range nodes {
				wg.Add(1)
				nodeInfo.Wg = wg
				nodeInfo.Stop(args...)
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}

		for _, nodeInfo := range cliconfig.Cfg.GetNodes() {
			wg.Add(1)
			nodeInfo.Wg = wg
			nodeInfo.Stop(args...)
		}
		wg.Wait()

	},
}

func init() {
	StopCmd.Flags().BoolVarP(&stopAll, "all", "a", false, "stop all")
	rootCmd.AddCommand(StopCmd)

}
