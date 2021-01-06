package command

import (
	"fmt"
	"scs/client/cliconfig"
	"scs/client/node"
	"sync"

	"github.com/spf13/cobra"
)

var updateAll bool

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update script",
	Long:  `update script`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && !updateAll {
			fmt.Println("Specify at least one parameter, or -- all")
			return
		}
		if node.UseNodes != "" {
			if nodeInfo, ok := cliconfig.Cfg.Nodes[node.UseNodes]; ok {
				nodeInfo.Update(args...)
				return
			}
		}
		if node.GroupName != "" {
			wg := &sync.WaitGroup{}
			for _, v := range cliconfig.Cfg.Group[node.GroupName] {
				if nodeInfo, ok := cliconfig.Cfg.Nodes[v]; ok {
					wg.Add(1)
					nodeInfo.Wg = wg
					nodeInfo.Update(args...)
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
			nodeInfo.Update(args...)
		}
		wg.Wait()

	},
}

func init() {
	UpdateCmd.Flags().BoolVar(&updateAll, "all", false, "update all")
	rootCmd.AddCommand(UpdateCmd)

}
