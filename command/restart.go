package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs"

	"github.com/spf13/cobra"
)

var restartAll bool
var RestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart assign script",
	Long:  `command: scsctl restart ([flags]) || ([pname] [name])`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && !restartAll {
			fmt.Println("Specify at least one parameter, or -- all")
			return
		}
		if restartAll {
			args = nil
		}
		if scs.UseNodes != "" {
			if nodeInfo, ok := scs.CCfg.GetNode(scs.UseNodes); ok {
				nodeInfo.Restart(args...)

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
				nodeInfo.Restart(args...)
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}

		for _, nodeInfo := range scs.CCfg.GetNodes() {
			wg.Add(1)
			nodeInfo.Wg = wg
			nodeInfo.Restart(args...)
		}
		wg.Wait()

	},
}

func init() {
	RestartCmd.Flags().BoolVar(&restartAll, "all", false, "restart all")
	rootCmd.AddCommand(RestartCmd)
}
