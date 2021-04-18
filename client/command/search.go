package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client"
	"github.com/spf13/cobra"
)

var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "search package",
	Long:  `search package`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if UseNodes != "" {
			if nodeInfo, ok := client.CCfg.Nodes[UseNodes]; ok {
				nodeInfo.Search(args[0])

			} else {
				fmt.Println("not found this node")
			}
			return
		}
		if GroupName != "" {
			wg := &sync.WaitGroup{}
			for _, v := range client.CCfg.Group[GroupName] {
				if nodeInfo, ok := client.CCfg.Nodes[v]; ok {
					wg.Add(1)
					go func() {
						nodeInfo.Search(args[0])
						wg.Done()
					}()

				}
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}
		for name, nodeInfo := range client.CCfg.Nodes {
			nodeInfo.Name = name
			wg.Add(1)
			go func() {
				nodeInfo.Search(args[0])
				wg.Done()
			}()
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(SearchCmd)
}
