package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client/cliconfig"
	"github.com/hyahm/scs/client/node"

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
			if nodeInfo, ok := cliconfig.Cfg.Nodes[node.UseNodes]; ok {
				nodeInfo.Stop(args...)
				return
			}
		}
		if node.GroupName != "" {
			wg := &sync.WaitGroup{}
			for _, v := range cliconfig.Cfg.Group[node.GroupName] {
				if nodeInfo, ok := cliconfig.Cfg.Nodes[v]; ok {
					wg.Add(1)
					nodeInfo.Wg = wg
					nodeInfo.Stop(args...)
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
			nodeInfo.Stop(args...)
		}
		wg.Wait()

	},
}

func init() {
	StopCmd.Flags().BoolVarP(&stopAll, "all", "a", false, "stop all")
	rootCmd.AddCommand(StopCmd)

}
