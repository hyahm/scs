package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs"
	"github.com/spf13/cobra"
)

var KillCmd = &cobra.Command{
	Use:   "kill",
	Short: "kill script",
	Long:  `command: scsctl kill (<pname> [name])`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if scs.UseNodes != "" {
			if nodeInfo, ok := scs.CCfg.GetNode(scs.UseNodes); ok {
				nodeInfo.Kill(args...)

			} else {
				fmt.Println("not found this node")
			}
			return
		}
		if scs.GroupName != "" {
			wg := &sync.WaitGroup{}
			nodes := scs.CCfg.GetNodesInGroup(scs.GroupName)
			for _, nodeInfo := range nodes {
				wg.Add(1)
				nodeInfo.Wg = wg
				nodeInfo.Kill(args...)
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}

		for _, nodeInfo := range scs.CCfg.GetNodes() {
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
