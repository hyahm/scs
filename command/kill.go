package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/script"

	"github.com/spf13/cobra"
)

var KillCmd = &cobra.Command{
	Use:   "kill",
	Short: "kill script",
	Long:  `command: scsctl kill (<pname> [name])`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if script.UseNodes != "" {
			if nodeInfo, ok := script.CCfg.GetNode(script.UseNodes); ok {
				nodeInfo.Kill(args...)

			} else {
				fmt.Println("not found this node")
			}
			return
		}
		if script.GroupName != "" {
			wg := &sync.WaitGroup{}
			nodes := script.CCfg.GetNodesInGroup(script.GroupName)
			for _, nodeInfo := range nodes {
				wg.Add(1)
				nodeInfo.Wg = wg
				nodeInfo.Kill(args...)
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}

		for _, nodeInfo := range script.CCfg.GetNodes() {
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
