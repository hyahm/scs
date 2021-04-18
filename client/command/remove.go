package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client"
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
		if UseNodes != "" {
			if nodeInfo, ok := client.CCfg.GetNode(UseNodes); ok {
				nodeInfo.Remove(args...)

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
					nodeInfo.Remove(args...)
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
				nodeInfo.Remove(args...)
				wg.Done()
			}()

		}
		wg.Wait()

	},
}

func init() {
	RemoveCmd.Flags().BoolVarP(&removeAll, "all", "a", false, "remove all")
	rootCmd.AddCommand(RemoveCmd)

}
