package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs"
	"github.com/hyahm/scs/client"

	"github.com/spf13/cobra"
)

// var filter []string

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print assign script status",
	Long:  `command: scsctl status [flags] [pname] [name]`,
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		if UseNodes != "" {
			if nodeInfo, ok := client.CCfg.GetNode(UseNodes); ok {
				// nodeInfo.Filter = filter
				serverStatus, err := nodeInfo.Status(args...)
				if err == nil {
					serverStatus.SortAndPrint()
				} else {
					fmt.Println(err)
				}

			} else {
				fmt.Println("not found this node")
			}
			return
		}
		wg := &sync.WaitGroup{}
		ss := make([]*scs.ScriptStatusNode, 0)
		if GroupName != "" {
			nodes := client.CCfg.GetNodesInGroup(GroupName)
			for _, nodeInfo := range nodes {
				wg.Add(1)
				go func(nodeInfo *scs.Node) {
					serverStatus, err := nodeInfo.Status(args...)
					if err == nil {
						serverStatus.SortAndPrint()
					} else {
						fmt.Println(err)
					}
					wg.Done()
				}(nodeInfo)
			}
			wg.Wait()
			for _, s := range ss {
				s.SortAndPrint()
			}
			return
		}

		for _, nodeInfo := range client.CCfg.GetNodes() {
			wg.Add(1)
			go func(nodeInfo *scs.Node) {
				serverStatus, err := nodeInfo.Status(args...)
				if err == nil {
					serverStatus.SortAndPrint()
				} else {
					fmt.Println(err)
				}
				wg.Done()
			}(nodeInfo)
		}
		wg.Wait()
		for _, s := range ss {
			s.SortAndPrint()
		}
	},
}

func init() {
	// rootCmd.Flags().StringArrayVarP(&filter, "filter", "f", []string{}, "filter name")
	rootCmd.AddCommand(StatusCmd)
}
