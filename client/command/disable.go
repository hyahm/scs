package command

import (
	"fmt"
	"sync"

	"github.com/hyahm/scs/client"

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

		if UseNodes != "" {
			if nodeInfo, ok := client.CCfg.GetNode(UseNodes); ok {
				nodeInfo.Enable(args[0])

			} else {
				fmt.Println("not found this node")
			}
			return
		}
		wg := &sync.WaitGroup{}
		if GroupName != "" {
			nodes := client.CCfg.GetNodesInGroup(GroupName)
			for _, nodeInfo := range nodes {
				wg.Add(1)
				go func() {
					nodeInfo.Enable(args[0])
					wg.Done()
				}()
			}
			wg.Wait()
			return
		}
		for _, nodeInfo := range client.CCfg.GetNodes() {
			wg.Add(1)
			go func() {
				nodeInfo.Enable(args[0])
				wg.Done()
			}()

		}
		wg.Wait()
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
		if UseNodes != "" {
			if nodeInfo, ok := client.CCfg.GetNode(UseNodes); ok {
				nodeInfo.Disable(args[0])
			} else {
				fmt.Println("not found this node")
			}
			return
		}
		wg := &sync.WaitGroup{}
		if GroupName != "" {

			nodes := client.CCfg.GetNodesInGroup(GroupName)
			for _, nodeInfo := range nodes {
				wg.Add(1)
				go func() {
					nodeInfo.Disable(args[0])
					wg.Done()
				}()

			}
			wg.Wait()
			return
		}
		for _, nodeInfo := range client.CCfg.GetNodes() {
			wg.Add(1)
			go func() {
				nodeInfo.Disable(args[0])
				wg.Done()
			}()

		}
		wg.Wait()

	},
}

func init() {
	rootCmd.AddCommand(EnableCmd)
	rootCmd.AddCommand(DisableCmd)
}
