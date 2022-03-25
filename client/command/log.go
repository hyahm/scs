package command

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func logConfig(cmd *cobra.Command, args []string) {
	nodes := getNodes()
	if len(nodes) == 0 {
		fmt.Println("not found any nodes")
		return
	}
	line := 10
	if len(args) >= 2 {
		var err error
		line, err = strconv.Atoi(args[1])
		if err != nil {
			line = 10
		}
	}
	nodes[0].Log(args[0], line)
}

var LogCmd = &cobra.Command{
	Use:   "log",
	Short: "script log",
	Long:  `command: scsctl log [flags] <name>`,
	Args:  cobra.MinimumNArgs(1),
	Run:   logConfig,
}

var tail bool

func init() {
	LogCmd.Flags().BoolVarP(&tail, "tail", "f", false, "tailf ")
	rootCmd.AddCommand(LogCmd)
}
