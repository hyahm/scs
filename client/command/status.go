package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs"

	"github.com/spf13/cobra"
)

var filter []string

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print assign script status",
	Long:  `command: scsctl status [flags] [pname] [name]`,
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		if scs.UseNodes != "" {
			if nodeInfo, ok := scs.CCfg.GetNode(scs.UseNodes); ok {
				nodeInfo.Filter = filter
				if err := nodeInfo.Status(args...); err == nil {
					nodeInfo.Result.SortAndPrint()
				}

			} else {
				fmt.Println("not found this node")
			}
			return
		}
		ss := make([]*scs.ScriptStatusNode, 0)
		if scs.GroupName != "" {
			wg := &sync.WaitGroup{}
			nodes := scs.CCfg.GetNodesInGroup(scs.GroupName)
			for _, nodeInfo := range nodes {
				wg.Add(1)
				nodeInfo.Wg = wg
				nodeInfo.Filter = filter
				go func(nodeInfo *scs.Node) {
					if err := nodeInfo.Status(args...); err == nil {
						ss = append(ss, nodeInfo.Result)
					}
				}(nodeInfo)
			}
			wg.Wait()
			for _, s := range ss {
				s.SortAndPrint()
			}
			return
		}
		wg := &sync.WaitGroup{}

		for _, nodeInfo := range scs.CCfg.GetNodes() {
			wg.Add(1)
			nodeInfo.Wg = wg
			nodeInfo.Filter = filter
			go func(nodeInfo *scs.Node) {
				if err := nodeInfo.Status(args...); err == nil {
					ss = append(ss, nodeInfo.Result)
				}
			}(nodeInfo)
		}
		wg.Wait()
		for _, s := range ss {
			s.SortAndPrint()
		}
	},
}

func init() {
	rootCmd.Flags().StringArrayVarP(&filter, "filter", "f", []string{}, "filter name")
	rootCmd.AddCommand(StatusCmd)
}
