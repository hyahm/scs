package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client"
	"github.com/spf13/cobra"
)

var KillCmd = &cobra.Command{
	Use:   "kill",
	Short: "kill script",
	Long:  `command: scsctl kill (<pname> [name])`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if UseNodes != "" {
			if nodeInfo, ok := client.CCfg.GetNode(UseNodes); ok {
				nodeInfo.Kill(args...)

			} else {
				fmt.Println("not found this node")
			}
			return
		}
		if GroupName != "" {
			wg := &sync.WaitGroup{}
			nodes := client.CCfg.GetNodesInGroup(GroupName)
			for _, nodeInfo := range nodes {
				wg.Add(1)
				go func() {
					nodeInfo.Kill(args...)
					wg.Done()
				}()
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}

		for _, nodeInfo := range client.CCfg.GetNodes() {
			wg.Add(1)
			go func() {
				nodeInfo.Kill(args...)
				wg.Done()
			}()

		}
		wg.Wait()

	},
}

func init() {
	rootCmd.AddCommand(KillCmd)

}
