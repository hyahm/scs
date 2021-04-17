package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs"
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
		if scs.UseNodes != "" {
			if nodeInfo, ok := scs.CCfg.GetNode(scs.UseNodes); ok {
				nodeInfo.Update(args...)

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
				nodeInfo.Update(args...)
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}

		for _, nodeInfo := range scs.CCfg.GetNodes() {
			wg.Add(1)
			nodeInfo.Wg = wg
			nodeInfo.Update(args...)
		}
		wg.Wait()
	},
}

func init() {
	UpdateCmd.Flags().BoolVarP(&updateAll, "all", "a", false, "update all")
	rootCmd.AddCommand(UpdateCmd)

}
