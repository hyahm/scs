package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "get info",
	Long:  `command: scsctl get servers | scripts | alarms`,
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
			go func(node *client.Node, flag string) {
				switch flag {
				case "scripts":
					node.GetScripts()
					wg.Done()
				case "servers":
					node.GetServers()
					wg.Done()
				default:
					node.GetAlerts()
					wg.Done()
				}

			}(node, args[0])

		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(GetCmd)
}
