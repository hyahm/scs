package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client"
	"github.com/spf13/cobra"
)

var stopAll bool

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop script",
	Long:  `command: scsctl stop ([flags]) || ([pname] [name])`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && !stopAll {
			fmt.Println("Specify at least one parameter, or -- all")
			return
		}
		if stopAll {
			args = nil
		}
		wg := &sync.WaitGroup{}
		nodes := getNodes()
		if len(nodes) == 0 {
			fmt.Println("not found any nodes")
			return
		}
		for _, node := range nodes {
			wg.Add(1)
			go func(node *client.Node) {
				node.Stop(args...)
				wg.Done()
			}(node)

		}
		wg.Wait()

	},
}

func init() {
	StopCmd.Flags().BoolVarP(&stopAll, "all", "a", false, "stop all")
	rootCmd.AddCommand(StopCmd)

}
