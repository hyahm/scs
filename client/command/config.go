package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client"

	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "scs server config",
	Args:  cobra.MinimumNArgs(1),
	Long:  `All software has versions. This is Hugo's`,

	Run: func(cmd *cobra.Command, args []string) {

	},
}

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Print all script status",
	Long:  `All software has versions. This is Hugo's`,

	Run: func(cmd *cobra.Command, args []string) {
		client.CCfg.PrintNodes()

	},
}

var ReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "reload scs server config",
	Long:  `reload scs server config`,
	Args:  cobra.MaximumNArgs(0),
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
				node.Reload()
				wg.Done()
			}(node)

		}
		wg.Wait()
	},
}

func init() {
	ConfigCmd.AddCommand(ShowCmd)
	ConfigCmd.AddCommand(ReloadCmd)
	rootCmd.AddCommand(ConfigCmd)
}
