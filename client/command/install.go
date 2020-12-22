package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"scs/client/cliconfig"
	"scs/client/node"
	"scs/internal"
	"sync"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// install <-f yaml> || <-u url> || <name>
var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install package",
	Long:  `install package`,
	Run: func(cmd *cobra.Command, args []string) {
		condition := 0
		sc := &internal.Script{}
		if len(args) > 1 {
			condition++
		}
		if url != "" {
			condition++
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = json.NewDecoder(resp.Body).Decode(sc)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if file != "" {
			condition++
			f, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = yaml.Unmarshal(f, sc)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if condition == 0 {
			fmt.Println("at lease one params : scsctl install <-f yaml> || <-u url> || <name>")
			return
		}
		if condition > 1 {
			fmt.Println("only one params : scsctl install <-f yaml> || <-u url> || <name>")
			return
		}
		if node.UseNodes != "" {
			if nodeInfo, ok := cliconfig.Cfg.Nodes[node.UseNodes]; ok {
				nodeInfo.Install(sc, env)
				return
			}
		}
		if node.GroupName != "" {
			wg := &sync.WaitGroup{}
			for _, v := range cliconfig.Cfg.Group[node.GroupName] {
				if nodeInfo, ok := cliconfig.Cfg.Nodes[v]; ok {
					wg.Add(1)
					nodeInfo.Wg = wg
					nodeInfo.Install(sc, env)
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
			nodeInfo.Install(sc, env)
		}
		wg.Wait()

	},
}
var env map[string]string
var file string
var url string

func init() {
	InstallCmd.Flags().StringToStringVarP(&env, "env", "e", nil, "set env")
	InstallCmd.Flags().StringVarP(&file, "file", "f", "", "install from file")
	InstallCmd.Flags().StringVarP(&url, "url", "u", "", "install from url")
	rootCmd.AddCommand(InstallCmd)
}
