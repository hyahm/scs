package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/cliconfig"
	"github.com/hyahm/scs/server/node"

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
		if updateAll {
			args = nil
		}
		if node.UseNodes != "" {
			if nodeInfo, ok := cliconfig.Cfg.GetNode(node.UseNodes); ok {
				nodeInfo.Update(args...)

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
				nodeInfo.Update(args...)
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}

		for _, nodeInfo := range cliconfig.Cfg.GetNodes() {
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
