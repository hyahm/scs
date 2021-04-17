package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/script"

	"github.com/spf13/cobra"
)

var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "search package",
	Long:  `search package`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if script.UseNodes != "" {
			if nodeInfo, ok := script.CCfg.Nodes[script.UseNodes]; ok {
				nodeInfo.Search(args[0])

			} else {
				fmt.Println("not found this node")
			}
			return
		}
		if script.GroupName != "" {
			wg := &sync.WaitGroup{}
			for _, v := range script.CCfg.Group[script.GroupName] {
				if nodeInfo, ok := script.CCfg.Nodes[v]; ok {
					wg.Add(1)
					nodeInfo.Wg = wg
					nodeInfo.Search(args[0])
				}
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}
		for name, nodeInfo := range script.CCfg.Nodes {
			nodeInfo.Name = name
			wg.Add(1)
			nodeInfo.Wg = wg
			nodeInfo.Search(args[0])
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(SearchCmd)
}
