package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client"
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
		if UseNodes != "" {
			if nodeInfo, ok := client.CCfg.GetNode(UseNodes); ok {
				nodeInfo.Restart(args...)

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
					nodeInfo.Restart(args...)
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
				nodeInfo.Restart(args...)
				wg.Done()
			}()

		}
		wg.Wait()

	},
}

func init() {
	RestartCmd.Flags().BoolVar(&restartAll, "all", false, "restart all")
	rootCmd.AddCommand(RestartCmd)
}
