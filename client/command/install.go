package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/hyahm/scs"
	"github.com/hyahm/scs/client"
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
		sc := make([]*scs.Script, 0)
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
			err = json.NewDecoder(resp.Body).Decode(&sc)
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
			err = yaml.Unmarshal(f, &sc)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if condition != 1 {
			fmt.Println("at lease one params : scsctl install <-f yaml> || <-u url> || <name>")
			return
		}

		if UseNodes != "" {
			if nodeInfo, ok := client.CCfg.GetNode(UseNodes); ok {
				nodeInfo.Install(sc, env)
				return
			}
		}
		if GroupName != "" {
			wg := &sync.WaitGroup{}
			nodes := client.CCfg.GetNodesInGroup(GroupName)
			for _, nodeInfo := range nodes {
				wg.Add(1)
				go func() {
					nodeInfo.Install(sc, env)
					wg.Done()
				}()
			}
			wg.Wait()
			return
		}
		wg := &sync.WaitGroup{}

		for _, nodeInfo := range client.CCfg.GetNodes() {
			wg.Add(1)
			go func() {
				nodeInfo.Install(sc, env)
				wg.Done()
			}()

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
