package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/script"

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
		script.CCfg.PrintNodes()

	},
}

var ReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "reload scs server config",
	Long:  `reload scs server config`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if script.UseNodes != "" {
			if nodeInfo, ok := script.CCfg.GetNode(script.UseNodes); ok {
				nodeInfo.Reload()
			} else {
				fmt.Println("not found this node")
			}
			return
		}
		if script.GroupName != "" {
			wg := &sync.WaitGroup{}
			nodeinfos := script.CCfg.GetNodesInGroup(script.GroupName)
			for _, nodeInfo := range nodeinfos {
				wg.Add(1)
				nodeInfo.Wg = wg
				nodeInfo.Reload()
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}
		for _, nodeInfo := range script.CCfg.GetNodes() {
			wg.Add(1)
			nodeInfo.Wg = wg
			nodeInfo.Reload()
		}
		wg.Wait()
	},
}

func init() {
	ConfigCmd.AddCommand(ShowCmd)
	ConfigCmd.AddCommand(ReloadCmd)
	rootCmd.AddCommand(ConfigCmd)
}
