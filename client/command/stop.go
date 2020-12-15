package command

import (
	"fmt"
	"scs/client/cliconfig"
	"scs/client/node"
	"sync"

	"github.com/spf13/cobra"
)

var stopAll bool

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop script",
	Long:  `stop script`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && !stopAll {
			fmt.Println("Specify at least one parameter, or -- all")
			return
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
	StopCmd.Flags().BoolVar(&stopAll, "all", false, "stop all")
	rootCmd.AddCommand(StopCmd)

}
