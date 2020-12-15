package command

import (
	"scs/client/cliconfig"
	"scs/client/node"
	"sync"

	"github.com/spf13/cobra"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install package",
	Long:  `install package`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if node.UseNodes != "" {
			if nodeInfo, ok := cliconfig.Cfg.Nodes[node.UseNodes]; ok {
				nodeInfo.Install(args[0], env)
				return
			}
		}
		if node.GroupName != "" {
			wg := &sync.WaitGroup{}
			for _, v := range cliconfig.Cfg.Group[node.GroupName] {
				if nodeInfo, ok := cliconfig.Cfg.Nodes[v]; ok {
					wg.Add(1)
					nodeInfo.Wg = wg
					nodeInfo.Install(args[0], env)
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
			nodeInfo.Install(args[0], env)
		}
		wg.Wait()

	},
}
var env map[string]string

func init() {
	InstallCmd.Flags().StringToStringVarP(&env, "env", "e", nil, "set env")
	rootCmd.AddCommand(InstallCmd)
}
