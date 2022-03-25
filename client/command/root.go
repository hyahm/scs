package command

import (
	"fmt"
	"os"
	"time"

	"github.com/hyahm/scs/client"
	"github.com/hyahm/scs/global"

	"github.com/spf13/cobra"
)

var UseNodes string
var GroupName string
var ReadTimeout time.Duration

var rootCmd = &cobra.Command{
	Use:     "scsctl",
	Version: global.VERSION,
	Short:   "scs is a server or script manager service",
	Long:    `scs service help, version: ` + global.VERSION,
	Args:    cobra.MinimumNArgs(1),
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&UseNodes, "node", "n", "", "show which nodes")
	rootCmd.PersistentFlags().StringVarP(&GroupName, "group", "g", "", "show which groupname")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getNodes() []*client.Node {
	nodes := make([]*client.Node, 0)
	if UseNodes != "" {
		if nodeInfo, ok := client.CCfg.GetNode(UseNodes); ok {
			nodes = append(nodes, nodeInfo)
		}
		return nodes

	}
	if GroupName != "" {
		return client.CCfg.GetNodesInGroup(GroupName)
	}
	return client.CCfg.GetNodes()
}
