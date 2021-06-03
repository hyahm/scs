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
		wg := &sync.WaitGroup{}
		nodes := getNodes()
		if len(nodes) == 0 {
			fmt.Println("not found any nodes")
			return
		}
		for _, node := range nodes {
			wg.Add(1)
			go func(node *client.Node) {
				node.Remove(args...)
				wg.Done()
			}(node)

		}
		wg.Wait()

	},
}

func init() {
	RemoveCmd.Flags().BoolVarP(&removeAll, "all", "a", false, "remove all")
	rootCmd.AddCommand(RemoveCmd)

}
