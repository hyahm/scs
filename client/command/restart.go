package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client"
	"github.com/spf13/cobra"
)

var RestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart assign script",
	Long:  `command: scsctl restart ([flags]) || ([pname] [name])`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(parameter)
		if len(args) == 0 && !restartAll {
			fmt.Println("Specify at least one parameter, or --all")
			return
		}
		if restartAll {
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
				node.Restart(parameter, args...)
				wg.Done()
			}(node)

		}
		wg.Wait()

	},
}

func init() {
	RestartCmd.Flags().BoolVarP(&restartAll, "all", "a", false, "restart all")
	RestartCmd.Flags().StringVarP(&parameter, "parameter", "p", "", "restart all")
	rootCmd.AddCommand(RestartCmd)
}
