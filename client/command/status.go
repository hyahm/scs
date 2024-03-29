package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client"

	"github.com/spf13/cobra"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print assign script status",
	Long:  `command: scsctl status [flags] [pname] [name]`,
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ss := make([]*client.ScriptStatusNode, 0)
		wg := &sync.WaitGroup{}
		nodes := getNodes()
		if len(nodes) == 0 {
			fmt.Println("not found any nodes")
			return
		}
		for _, node := range getNodes() {
			wg.Add(1)
			go func(node *client.Node) {
				serverStatus, err := node.Status(args...)
				if err == nil {
					ss = append(ss, serverStatus)
				} else {
					fmt.Println(err)
				}
				wg.Done()
			}(node)

		}
		wg.Wait()
		for _, s := range ss {
			s.SortAndPrint(verbose)
		}
	},
}

func init() {
	StatusCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show verbose")
	rootCmd.AddCommand(StatusCmd)
}
