package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client"
	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "start assign script",
	Long:  `command: scsctl start [flags] [pname] [name]`,
	Args:  cobra.MaximumNArgs(2),
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
				node.Start(args...)
				wg.Done()
			}(node)

		}
		wg.Wait()
	},
}

func init() {

	rootCmd.AddCommand(StartCmd)

}
