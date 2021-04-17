package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/script"
	"github.com/spf13/cobra"
)

var stopAll bool

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop script",
	Long:  `command: scsctl stop ([flags]) || ([pname] [name])`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && !stopAll {
			fmt.Println("Specify at least one parameter, or -- all")
			return
		}
		if stopAll {
			args = nil
		}
		if script.UseNodes != "" {
			if nodeInfo, ok := script.CCfg.GetNode(script.UseNodes); ok {
				nodeInfo.Stop(args...)

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
				nodeInfo.Stop(args...)
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}

		for _, nodeInfo := range script.CCfg.GetNodes() {
			wg.Add(1)
			nodeInfo.Wg = wg
			nodeInfo.Stop(args...)
		}
		wg.Wait()

	},
}

func init() {
	StopCmd.Flags().BoolVarP(&stopAll, "all", "a", false, "stop all")
	rootCmd.AddCommand(StopCmd)

}
