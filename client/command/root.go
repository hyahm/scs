package command

import (
	"fmt"
	"os"

	"github.com/hyahm/scs/client/node"
	"github.com/hyahm/scs/global"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "scsctl",
	Version: global.VERSION,
	Short:   "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&node.UseNodes, "node", "n", "", "set nodes ,have priority over group")
	rootCmd.PersistentFlags().StringVarP(&node.GroupName, "group", "g", "", "set group, if node group net set, all script will be use")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
