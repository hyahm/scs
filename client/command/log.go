package command

import (
	"fmt"
	"sync"
	"time"

	"github.com/hyahm/scs"
	"github.com/spf13/cobra"
)

func logConfig(cmd *cobra.Command, args []string) {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start).Seconds())
	}()
	wg := &sync.WaitGroup{}
	nodes := getNodes()
	if len(nodes) == 0 {
		fmt.Println("not found any nodes")
		return
	}
	for _, node := range nodes {
		wg.Add(1)
		go func(node *scs.Node) {
			node.Log(args[0])
			wg.Done()
		}(node)

	}
	wg.Wait()

}

var LogCmd = &cobra.Command{
	Use:   "log",
	Short: "script log",
	Long:  `command: scsctl log [flags] <name>`,
	Args:  cobra.MinimumNArgs(1),
	Run:   logConfig,
}

var tail bool

func init() {
	LogCmd.Flags().BoolVarP(&tail, "tail", "f", false, "tailf ")
	rootCmd.AddCommand(LogCmd)
}
