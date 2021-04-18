package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs"
	"github.com/spf13/cobra"
)

var EnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "enable script",
	Long:  `command: scsctl enable <pname>`,
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
			go func(node *scs.Node) {
				node.Enable(args[0])
				wg.Done()
			}(node)

		}
		wg.Wait()
	},
}

var DisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "disable script",
	Long:  `command: scsctl disable <pname>`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		wg := &sync.WaitGroup{}
		nodes := getNodes()
		if len(nodes) == 0 {
			fmt.Println("not found any nodes")
			return
		}
		for _, node := range getNodes() {
			wg.Add(1)
			go func() {
				node.Disable(args[0])
				wg.Done()
			}()

		}
		wg.Wait()

	},
}

func init() {
	rootCmd.AddCommand(EnableCmd)
	rootCmd.AddCommand(DisableCmd)
}
