package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/cliconfig"
	"github.com/hyahm/scs/server/node"

	"github.com/spf13/cobra"
)

var removeAll bool

var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove script",
	Long:  `command: scsctl remove ([flags]) || ([pname] [name])`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && !removeAll {
			fmt.Println("Specify at least one parameter, or -- all")
			return
		}
		if removeAll {
			args = nil
		}
		if node.UseNodes != "" {
			if nodeInfo, ok := cliconfig.Cfg.GetNode(node.UseNodes); ok {
				nodeInfo.Remove(args...)

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
				nodeInfo.Remove(args...)
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}

		for _, nodeInfo := range cliconfig.Cfg.GetNodes() {
			wg.Add(1)
			nodeInfo.Wg = wg
			nodeInfo.Remove(args...)
		}
		wg.Wait()

	},
}

func init() {
	RemoveCmd.Flags().BoolVarP(&removeAll, "all", "a", false, "remove all")
	rootCmd.AddCommand(RemoveCmd)

}
