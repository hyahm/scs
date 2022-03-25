package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client"
	"github.com/spf13/cobra"
)

var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "show service",
	Long:  `command: scsctl env [flags] <name>`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		wg := &sync.WaitGroup{}
		nodes := getNodes()
		if len(nodes) == 0 {
			fmt.Println("not found any nodes")
			return
		}
		for _, node := range nodes {
			wg.Add(1)
			go func(node *client.Node) {
				node.Info(args[0])
				wg.Done()
			}(node)

		}
		wg.Wait()

	},
}

func init() {
	rootCmd.AddCommand(InfoCmd)
}
