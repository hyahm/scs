package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/script"

	"github.com/spf13/cobra"
)

var EnvCmd = &cobra.Command{
	Use:   "env",
	Short: "show env",
	Long:  `command: scsctl env [flags] <name>`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if script.UseNodes != "" {
			if nodeInfo, ok := script.CCfg.GetNode(script.UseNodes); ok {
				nodeInfo.Env(args[0])
			} else {
				fmt.Println("not found this node")
			}
			return
		}
		if script.GroupName != "" {
			wg := &sync.WaitGroup{}
			nodes := script.CCfg.GetNodesInGroup(script.GroupName)
			for _, nodeInfo := range nodes {
				wg.Add(1)
				nodeInfo.Wg = wg
				nodeInfo.Env(args[0])
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}

		for _, nodeInfo := range script.CCfg.GetNodes() {
			wg.Add(1)
			nodeInfo.Wg = wg
			nodeInfo.Env(args[0])
		}
		wg.Wait()

	},
}

// var ShowEnvCmd = &cobra.Command{
// 	Use:   "show",
// 	Short: "show env",
// 	Long:  `show env`,
// 	Args:  cobra.MinimumNArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		if node.UseNodes != "" {
// 			if nodeInfo, ok := cliconfig.Cfg.Nodes[node.UseNodes]; ok {
// 				nodeInfo.Env(args...)
// 				return
// 			}
// 		}
// 		if node.GroupName != "" {
// 			wg := &sync.WaitGroup{}
// 			for _, v := range cliconfig.Cfg.Group[node.GroupName] {
// 				if nodeInfo, ok := cliconfig.Cfg.Nodes[v]; ok {
// 					wg.Add(1)
// 					nodeInfo.Wg = wg
// 					nodeInfo.Env(args...)
// 				}
// 			}
// 			wg.Wait()
// 			return
// 		}
// 		wg := &sync.WaitGroup{}
// 		for name, nodeInfo := range cliconfig.Cfg.Nodes {
// 			nodeInfo.Name = name
// 			wg.Add(1)
// 			nodeInfo.Wg = wg
// 			nodeInfo.Env(args...)
// 		}
// 		wg.Wait()

// 	},
// }

func init() {
	// EnvCmd.AddCommand(ShowEnvCmd)
	rootCmd.AddCommand(EnvCmd)
}