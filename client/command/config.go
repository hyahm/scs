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
		if UseNodes != "" {
			if nodeInfo, ok := client.CCfg.GetNode(UseNodes); ok {
				nodeInfo.Reload()

			} else {
				fmt.Println("not found this node")
			}
			return
		}
		wg := &sync.WaitGroup{}
		if GroupName != "" {

			nodeinfos := client.CCfg.GetNodesInGroup(GroupName)
			for _, nodeInfo := range nodeinfos {
				wg.Add(1)
				go func() {
					nodeInfo.Reload()
					wg.Done()
				}()

			}
			wg.Wait()
			return
		}
		for _, nodeInfo := range client.CCfg.GetNodes() {
			wg.Add(1)
			go func() {
				nodeInfo.Reload()
				wg.Done()
			}()
		}
		wg.Wait()
	},
}

func init() {
	ConfigCmd.AddCommand(ShowCmd)
	ConfigCmd.AddCommand(ReloadCmd)
	rootCmd.AddCommand(ConfigCmd)
}
