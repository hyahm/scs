package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client/cliconfig"
	"github.com/hyahm/scs/client/node"

	"github.com/spf13/cobra"
)

var EnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "enable script",
	Long:  `command: scsctl enable <pname>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Specify at least one parameter, or -- all")
			return
		}

		if node.UseNodes != "" {
			if nodeInfo, ok := cliconfig.Cfg.GetNode(node.UseNodes); ok {
				nodeInfo.Enable(args[0])

			} else {
				fmt.Println("not found this node")
			}
			return
		}
		if node.GroupName != "" {
			wg := &sync.WaitGroup{}
			nodes := cliconfig.Cfg.GetNodesInGroup(node.GroupName)
			for _, nodeInfo := range nodes {
				wg.Add(1)
				nodeInfo.Wg = wg
				nodeInfo.Enable(args[0])
			}
			wg.Wait()
			return
		}

	},
}

var DisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "disable script",
	Long:  `command: scsctl disable <pname>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Specify at least one parameter, or -- all")
			return
		}
		if node.UseNodes != "" {
			if nodeInfo, ok := cliconfig.Cfg.GetNode(node.UseNodes); ok {
				nodeInfo.Disable(args[0])
			} else {
				fmt.Println("not found this node")
			}
			return
		}
		if node.GroupName != "" {
			wg := &sync.WaitGroup{}
			nodes := cliconfig.Cfg.GetNodesInGroup(node.GroupName)
			for _, nodeInfo := range nodes {
				wg.Add(1)
				nodeInfo.Wg = wg
				nodeInfo.Disable(args[0])
			}
			wg.Wait()
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(EnableCmd)
	rootCmd.AddCommand(DisableCmd)
}