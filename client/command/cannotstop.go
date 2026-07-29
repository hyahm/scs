package command

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
)

var CanNotStopCmd = &cobra.Command{
	Use:   "cannotstop",
	Short: "cannotstop script",
	Long:  `command: scsctl cannotstop <pname|name>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Specify at least one parameter")
			return
		}
		wg := &sync.WaitGroup{}
		nodes := getNodes()
		if len(nodes) == 0 {
			fmt.Println("not found any nodes")
			return
		}
		for _, node := range nodes {
			wg.Go(func() {
				node.CanNotStop(args...)
			})
			// go func(node *client.Node) {
			// 	node.Stop(args...)
			// 	wg.Done()
			// }(node)

		}
		wg.Wait()

	},
}

var CanStopCmd = &cobra.Command{
	Use:   "canstop",
	Short: "canstop script",
	Long:  `command: scsctl canstop <pname|name>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Specify at least one parameter")
			return
		}
		wg := &sync.WaitGroup{}
		nodes := getNodes()
		if len(nodes) == 0 {
			fmt.Println("not found any nodes")
			return
		}
		for _, node := range nodes {
			wg.Go(func() {
				node.CanStop(args...)
			})
		}
		wg.Wait()

	},
}

func init() {
	// CanNotStopCmd.Flags().BoolVarP(&stopAll, "all", "a", false, "stop all")
	rootCmd.AddCommand(CanNotStopCmd, CanStopCmd)
}
