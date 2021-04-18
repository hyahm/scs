package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client"
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
		if UseNodes != "" {
			if nodeInfo, ok := client.CCfg.GetNode(UseNodes); ok {
				nodeInfo.Stop(args...)

			} else {
				fmt.Println("not found this node")
			}
			return
		}
		wg := &sync.WaitGroup{}
		if GroupName != "" {

			nodes := client.CCfg.GetNodesInGroup(GroupName)
			for _, nodeInfo := range nodes {
				wg.Add(1)
				go func() {
					nodeInfo.Stop(args...)
					wg.Done()
				}()

			}
			wg.Wait()
			return
		}

		for _, nodeInfo := range client.CCfg.GetNodes() {
			wg.Add(1)
			go func() {
				nodeInfo.Stop(args...)
				wg.Done()
			}()
		}
		wg.Wait()

	},
}

func init() {
	StopCmd.Flags().BoolVarP(&stopAll, "all", "a", false, "stop all")
	rootCmd.AddCommand(StopCmd)

}
