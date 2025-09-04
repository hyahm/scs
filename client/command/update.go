package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client"
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update server",
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
		wg := &sync.WaitGroup{}
		nodes := getNodes()
		if len(nodes) == 0 {
			fmt.Println("not found any nodes")
			return
		}
		for _, node := range nodes {
			wg.Add(1)
			go func(node *client.Node) {
				node.Update(args...)
				wg.Done()
			}(node)

		}
		wg.Wait()
	},
}

func init() {
	UpdateCmd.Flags().BoolVarP(&updateAll, "all", "a", false, "update all")
	rootCmd.AddCommand(UpdateCmd)

}
