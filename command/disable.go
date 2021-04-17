package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/script"

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

		if script.UseNodes != "" {
			if nodeInfo, ok := script.CCfg.GetNode(script.UseNodes); ok {
				nodeInfo.Enable(args[0])

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
		if script.UseNodes != "" {
			if nodeInfo, ok := script.CCfg.GetNode(script.UseNodes); ok {
				nodeInfo.Disable(args[0])
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
