package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/hyahm/scs/client"
	"github.com/hyahm/scs/pkg/config/scripts"
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
		sc := make([]*scripts.Script, 0)
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
			path, err := filepath.Abs(file)
			if err != nil {
				log.Fatal(err)
			}
			f, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = yaml.Unmarshal(f, &sc)
			if err != nil {
				fmt.Println(err)
				return
			}
			// 如果执行目录是空，那么将yaml文件目录当家目录
			for i := range sc {
				if sc[i].Dir == "" {
					if sc[i].Env == nil {
						sc[i].Env = make(map[string]string)
					}
					sc[i].Dir = filepath.Dir(path)
					sc[i].Env["PROJECT_HOME"] = sc[i].Dir
				}
			}

		}
		if condition != 1 {
			fmt.Println("at lease one params : scsctl install <-f yaml> || <-u url> || <name>")
			return
		}

		wg := &sync.WaitGroup{}
		nodes := getNodes()
		if len(nodes) == 0 {
			fmt.Println("not found any nodes")
			return
		}
		for _, node := range nodes {
			wg.Add(1)
			go func(node *client.Node) {
				node.Install(sc, nil)
				wg.Done()
			}(node)

		}
		wg.Wait()

	},
}

func init() {
	InstallCmd.Flags().StringToStringVarP(&env, "env", "e", nil, "set env")
	InstallCmd.Flags().StringVarP(&file, "file", "f", "", "install from file")
	InstallCmd.Flags().StringVarP(&url, "url", "u", "", "install from url")
	rootCmd.AddCommand(InstallCmd)
}
