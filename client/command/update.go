package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client"
	"github.com/spf13/cobra"
)

var updateAll bool

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update script",
	Long:  `command: scsctl update ([flags]) || ([pname] [name])`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if updateAll {
			fmt.Println("waiting update")
			return
		}
		if updateAll {
			args = nil
		}
		if UseNodes != "" {
			if nodeInfo, ok := client.CCfg.GetNode(UseNodes); ok {
				nodeInfo.Update(args...)

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
					nodeInfo.Update(args...)
					wg.Done()
				}()

			}
			wg.Wait()
			return
		}

		for _, nodeInfo := range client.CCfg.GetNodes() {
			wg.Add(1)
			go func() {
				nodeInfo.Update(args...)
				wg.Done()
			}()

		}
		wg.Wait()
	},
}

func init() {
	UpdateCmd.Flags().BoolVarP(&updateAll, "all", "a", false, "update all")
	rootCmd.AddCommand(UpdateCmd)

}
