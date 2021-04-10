package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/cliconfig"
	"github.com/hyahm/scs/server/node"

	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "start assign script",
	Long:  `command: scsctl start [flags] [pname] [name]`,
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if node.UseNodes != "" {
			if nodeInfo, ok := cliconfig.Cfg.GetNode(node.UseNodes); ok {
				nodeInfo.Start(args...)

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
				nodeInfo.Start(args...)
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}

		for _, nodeInfo := range cliconfig.Cfg.GetNodes() {
			wg.Add(1)
			nodeInfo.Wg = wg
			nodeInfo.Start(args...)
		}
		wg.Wait()
	},
}

func init() {

	rootCmd.AddCommand(StartCmd)

}
