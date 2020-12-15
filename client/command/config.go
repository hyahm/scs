package command

import (
	"fmt"
	"scs/client/cliconfig"
	"scs/client/node"
	"sync"

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
		for name, v := range cliconfig.Cfg.Nodes {
			fmt.Printf("name: %s \t url: %s \t token: %s \n", name, v.Url, v.Token)
		}
	},
}

var ReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "reload scs server config",
	Long:  `reload scs server config`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if node.UseNodes != "" {
			if nodeInfo, ok := cliconfig.Cfg.Nodes[node.UseNodes]; ok {
				nodeInfo.Reload()
				return
			}
		}
		if node.GroupName != "" {
			wg := &sync.WaitGroup{}
			for _, v := range cliconfig.Cfg.Group[node.GroupName] {
				if nodeInfo, ok := cliconfig.Cfg.Nodes[v]; ok {
					wg.Add(1)
					nodeInfo.Wg = wg
					nodeInfo.Reload()
				}
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}
		for name, nodeInfo := range cliconfig.Cfg.Nodes {
			nodeInfo.Name = name
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
