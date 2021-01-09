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
	Long:  `command: scsctl update ([flags]) || ([pname] [name])`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if updateAll {
			fmt.Println("waiting update")
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
	UpdateCmd.Flags().BoolVarP(&updateAll, "all", "a", false, "update all")
	rootCmd.AddCommand(UpdateCmd)

}
