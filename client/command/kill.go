package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client/cliconfig"
	"github.com/hyahm/scs/client/node"

	"github.com/spf13/cobra"
)

var KillCmd = &cobra.Command{
	Use:   "kill",
	Short: "kill script",
	Long:  `command: scsctl kill (<pname> [name])`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if node.UseNodes != "" {
			if nodeInfo, ok := cliconfig.Cfg.GetNode(node.UseNodes); ok {
				nodeInfo.Kill(args...)

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
				nodeInfo.Kill(args...)
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}

		for _, nodeInfo := range cliconfig.Cfg.GetNodes() {
			wg.Add(1)
			nodeInfo.Wg = wg
			nodeInfo.Kill(args...)
		}
		wg.Wait()

	},
}

func init() {
	rootCmd.AddCommand(KillCmd)

}
