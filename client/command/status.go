package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs"

	"github.com/spf13/cobra"
)

// var filter []string

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print assign script status",
	Long:  `command: scsctl status [flags] [pname] [name]`,
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		ss := make([]*scs.ScriptStatusNode, 0)

		wg := &sync.WaitGroup{}
		nodes := getNodes()
		if len(nodes) == 0 {
			fmt.Println("not found any nodes")
			return
		}

		for _, node := range getNodes() {
			wg.Add(1)
			go func(node *scs.Node) {
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
			s.SortAndPrint()
		}
	},
}

func init() {
	// rootCmd.Flags().StringArrayVarP(&filter, "filter", "f", []string{}, "filter name")
	rootCmd.AddCommand(StatusCmd)
}
